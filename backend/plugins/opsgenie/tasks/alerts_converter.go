/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tasks

import (
	"fmt"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/opsgenie/models"
)

var ConvertAlertsMeta = plugin.SubTaskMeta{
	Name:             "convertAlerts",
	EntryPoint:       ConvertAlerts,
	EnabledByDefault: true,
	Description:      "Convert alerts into domain layer table issues",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

type (
	// AlertWithUser struct that represents the joined query result
	AlertWithUser struct {
		common.NoPKModel
		models.Alert
		*models.User
		AssignedAt time.Time
	}
)

func ConvertAlerts(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*OpsgenieTaskData)
	cursor, err := db.Cursor(
		dal.Select("pi.*, pu.*, pa.assigned_at"),
		dal.From("_tool_opsgenie_alertss AS pi"),
		dal.Join(`LEFT JOIN _tool_opsgenie_assignments AS pa ON pa.alert_id = pi.id`),
		dal.Join(`LEFT JOIN _tool_opsgenie_users AS pu ON pa.user_id = pu.id`),
		dal.Where("pi.connection_id = ? AND pi.service_id = ?", data.Options.ConnectionId, data.Options.ServiceId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()
	seenAlerts := map[int]*AlertWithUser{}
	idGen := didgen.NewDomainIdGenerator(&models.Alert{})
	serviceIdGen := didgen.NewDomainIdGenerator(&models.Service{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_ALERTS_TABLE,
		},
		InputRowType: reflect.TypeOf(AlertWithUser{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			combined := inputRow.(*AlertWithUser)
			alert := combined.Alert
			if seen, ok := seenAlerts[alert.Id]; ok {
				if combined.AssignedAt.Before(seen.AssignedAt) {
					// skip this one (it's an older assignee)
					return nil, nil
				}
			}
			status := getStatus(&alert)
			leadTime, resolutionDate := getTimes(&alert)
			domainIssue := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: idGen.Generate(data.Options.ConnectionId, alert.Number),
				},
				Url:             alert.Url,
				IssueKey:        fmt.Sprintf("%d", alert.Id),
				Description:     alert.Description,
				Type:            ticket.INCIDENT,
				Status:          status,
				OriginalStatus:  string(alert.Status),
				ResolutionDate:  resolutionDate,
				CreatedDate:     &alert.CreatedDate,
				UpdatedDate:     &alert.UpdatedDate,
				LeadTimeMinutes: leadTime,
				Priority:        string(alert.Urgency),
			}
			var result []interface{}
			if combined.User != nil {
				domainIssue.AssigneeId = combined.User.Id
				domainIssue.AssigneeName = combined.User.UserName
				issueAssignee := &ticket.IssueAssignee{
					IssueId:      domainIssue.Id,
					AssigneeId:   domainIssue.AssigneeId,
					AssigneeName: domainIssue.AssigneeName,
				}
				result = append(result, issueAssignee)
			}
			result = append(result, domainIssue)
			seenAlerts[alert.Id] = combined
			boardIssue := &ticket.BoardIssue{
				BoardId: serviceIdGen.Generate(data.Options.ConnectionId, data.Options.ServiceId),
				IssueId: domainIssue.Id,
			}
			result = append(result, boardIssue)
			return result, nil
		},
	})
	if err != nil {
		return err
	}
	return converter.Execute()
}

func getStatus(alert *models.Alert) string {
	if alert.Status == models.AlertStatusTriggered {
		return ticket.TODO
	}
	if alert.Status == models.AlertStatusAcknowledged {
		return ticket.IN_PROGRESS
	}
	if alert.Status == models.AlertStatusResolved {
		return ticket.DONE
	}
	panic("unknown alert status encountered")
}

func getTimes(alert *models.Alert) (int64, *time.Time) {
	var leadTime int64
	var resolutionDate *time.Time
	if alert.Status == models.AlertStatusResolved {
		resolutionDate = &alert.UpdatedAt
		leadTime = int64(resolutionDate.Sub(alert.CreatedAt).Minutes())
	}
	return leadTime, resolutionDate
}
