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
	"encoding/json"
	"fmt"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/opsgenie/models"
	"github.com/apache/incubator-devlake/plugins/opsgenie/models/raw"
)

var _ plugin.SubTaskEntryPoint = ExtractAlerts

var ExtractAlertsMeta = plugin.SubTaskMeta{
	Name:             "extractAlerts",
	EntryPoint:       ExtractAlerts,
	EnabledByDefault: true,
	Description:      "Extract Opsgenie alerts",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractAlerts(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*OpsgenieTaskData)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_ALERTS_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			alertRaw := &raw.Alert{}

			err := errors.Convert(json.Unmarshal(row.Data, alertRaw))
			if err != nil {
				return nil, err
			}

			fmt.Println("")
			fmt.Println(string(row.Data))
			fmt.Println("")
			results := make([]interface{}, 0, 1)
			alert := models.Alert{
				ConnectionId:   data.Options.ConnectionId,
				Id:             *alertRaw.Id,
				Message:        *alertRaw.Message,
				ServiceId:      data.Options.ServiceId,
				ServiceName:    data.Options.ServiceName,
				IncidentId:     resolve(alertRaw.Details.IncidentId), // may not be associated with an Incident
				Owner:          resolve(alertRaw.Owner),
				Description:    resolve(alertRaw.Description),
				AckTime:        resolve(alertRaw.Report.AckTime),
				AcknowledgedBy: resolve(alertRaw.Report.AcknowledgedBy),
				CloseTime:      resolve(alertRaw.Report.CloseTime),
				ClosedBy:       resolve(alertRaw.Report.ClosedBy),
				Status:         models.AlertStatus(*alertRaw.Status),
				Priority:       models.AlertPriority(*alertRaw.Priority),
				CreatedAt:      *alertRaw.CreatedAt,
				UpdatedAt:      *alertRaw.UpdatedAt,
			}
			results = append(results, &alert)
			return results, nil
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}
