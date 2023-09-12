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
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/opsgenie/models"
)

var ConvertIncidentsMeta = plugin.SubTaskMeta{
	Name:             "convertIncidents",
	EntryPoint:       ConvertIncidents,
	EnabledByDefault: true,
	Description:      "Convert Incidents into domain layer table issues",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertIncidents(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*OpsgenieTaskData)

	cursor, err := db.Cursor(
		dal.From(models.Incident{}),
		dal.Where("connection_id = ? AND service_id = ?", data.Options.ConnectionId, data.Options.ServiceId),
	)

	if err != nil {
		return err
	}

	defer cursor.Close()

	idGen := didgen.NewDomainIdGenerator(&models.Incident{})
	serviceIdGen := didgen.NewDomainIdGenerator(&models.Service{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_INCIDENTS_TABLE,
		},
		InputRowType: reflect.TypeOf(models.Incident{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			incident := inputRow.(*models.Incident)
			status := getStatus(incident)
			leadTime, resolutionDate := getTimes(incident)
			domainIssue := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: idGen.Generate(data.Options.ConnectionId, incident.Id),
				},
				Url:             incident.Url,
				IssueKey:        incident.Id,
				Title:           incident.Message,
				Description:     incident.Description,
				Type:            ticket.INCIDENT,
				Status:          status,
				OriginalStatus:  string(incident.Status),
				ResolutionDate:  resolutionDate,
				CreatedDate:     &incident.CreatedAt,
				UpdatedDate:     &incident.UpdatedAt,
				LeadTimeMinutes: leadTime,
				Priority:        string(incident.Priority),
			}
			var result []interface{}
			result = append(result, domainIssue)
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

func getStatus(incident *models.Incident) string {
	if incident.Status == models.IncidentStatusClosed {
		return ticket.OTHER
	}
	if incident.Status == models.IncidentStatusOpen {
		return ticket.IN_PROGRESS
	}
	if incident.Status == models.IncidentStatusResolved {
		return ticket.DONE
	}
	panic("unknown incident status encountered")
}

func getTimes(incident *models.Incident) (int64, *time.Time) {
	var leadTime int64
	var resolutionDate *time.Time
	if incident.Status == models.IncidentStatusResolved {
		resolutionDate = &incident.UpdatedAt
		leadTime = int64(resolutionDate.Sub(incident.CreatedAt).Minutes())
	}
	return leadTime, resolutionDate
}
