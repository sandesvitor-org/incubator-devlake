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
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/opsgenie/models"
)

const RAW_ALERTS_TABLE = "opsgenie_alerts"

var _ plugin.SubTaskEntryPoint = CollectAlerts

type (
	collectedAlerts struct {
		TotalCount int               `json:"totalCount"`
		Data       []json.RawMessage `json:"data"`
	}

	collectedAlert struct {
		Data      json.RawMessage `json:"data"`
		Took      float64         `json:"took"`
		RequestId string          `json:"requestId"`
	}
	simplifiedRawAlert struct {
		Id        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	}
)

var CollectAlertsMeta = plugin.SubTaskMeta{
	Name:             "collectAlerts",
	EntryPoint:       CollectAlerts,
	EnabledByDefault: true,
	Description:      "Collect Opsgenie alerts",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func CollectAlerts(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*OpsgenieTaskData)
	db := taskCtx.GetDal()
	args := api.RawDataSubTaskArgs{
		Ctx:     taskCtx,
		Options: data.Options,
		Table:   RAW_ALERTS_TABLE,
	}
	collector, err := api.NewStatefulApiCollectorForFinalizableEntity(api.FinalizableApiCollectorArgs{
		RawDataSubTaskArgs: args,
		ApiClient:          data.Client,
		TimeAfter:          data.TimeAfter,
		CollectNewRecordsByList: api.FinalizableApiCollectorListArgs{
			PageSize: 100,
			GetNextPageCustomData: func(prevReqData *api.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
				pager := prevReqData.Pager
				if pager.Skip+pager.Size >= 20_000 { // API limit. Can't exceed this or it'll error out
					return nil, api.ErrFinishCollect
				}
				return nil, nil
			},
			FinalizableApiCollectorCommonArgs: api.FinalizableApiCollectorCommonArgs{
				UrlTemplate: "v2/alerts",
				Query: func(reqData *api.RequestData, createdAfter *time.Time) (url.Values, errors.Error) {
					query := url.Values{}

					query.Set("query", fmt.Sprintf("detailsPair(impacted-services:%s)", data.Options.ServiceId))
					query.Set("sort", "createdAt")
					query.Set("order", "desc")
					query.Set("limit", fmt.Sprintf("%d", reqData.Pager.Size))
					query.Set("offset", fmt.Sprintf("%d", reqData.Pager.Skip))
					return query, nil
				},
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					rawResult := collectedAlerts{}
					err := api.UnmarshalResponse(res, &rawResult)

					return rawResult.Data, err
				},
			},
		},
		// dal.Select("pi.*, pu.*, pa.assigned_at"),
		// dal.From("_tool_pagerduty_incidents AS pi"),
		// dal.Join(`LEFT JOIN _tool_pagerduty_assignments AS pa ON pa.incident_number = pi.number`),
		// dal.Join(`LEFT JOIN _tool_pagerduty_users AS pu ON pa.user_id = pu.id`),
		// dal.Where("pi.connection_id = ? AND pi.service_id = ?", data.Options.ConnectionId, data.Options.ServiceId),
		CollectUnfinishedDetails: &api.FinalizableApiCollectorDetailArgs{
			BuildInputIterator: func() (api.Iterator, errors.Error) {
				fmt.Printf("\n### DENTRO DE BuildInputIterator\n")
				cursor, err := db.Cursor(
					dal.Select("id"),
					dal.From(&models.Alert{}),
					// dal.Where(
					// 	"service_id = ? AND connection_id = ?",
					// 	data.Options.ServiceId, data.Options.ConnectionId,
					// ),
				)

				fmt.Println(cursor.Columns())
				fmt.Println(cursor)
				fmt.Println(err)
				// cursor, err := db.Cursor(
				// 	dal.Select("id, created_at"),
				// 	dal.From(&models.Alert{}),
				// 	dal.Where(
				// 		"service_id = ? AND connection_id = ? AND status != ?",
				// 		data.Options.ServiceId, data.Options.ConnectionId, "closed",
				// 	),
				// )

				if err != nil {
					return nil, err
				}
				return api.NewDalCursorIterator(db, cursor, reflect.TypeOf(simplifiedRawAlert{}))
			},
			FinalizableApiCollectorCommonArgs: api.FinalizableApiCollectorCommonArgs{
				UrlTemplate: "v2/alerts/{{ .Input.Id }}",
				Query:       nil,
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					fmt.Printf("\nDENTRO DO FinalizableApiCollectorCommonArgs.ResponseParser\n")
					rawResult := collectedAlert{}
					err := api.UnmarshalResponse(res, &rawResult)
					return []json.RawMessage{rawResult.Data}, err
				},
			},
		},
	})
	if err != nil {
		return nil
	}
	return collector.Execute()
}
