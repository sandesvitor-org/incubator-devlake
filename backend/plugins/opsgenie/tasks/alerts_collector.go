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
	pagingInfo struct {
		Limit      *int `json:"limit"`
		Offset     *int `json:"offset"`
		TotalCount *int `json:"totalCount"`
	}
	collectedAlerts struct {
		pagingInfo
		Alerts []json.RawMessage `json:"alerts"`
	}

	collectedAlert struct {
		pagingInfo
		Alert json.RawMessage `json:"alert"`
	}
	simplifiedRawAlert struct {
		Id        int       `json:"id"`
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
			FinalizableApiCollectorCommonArgs: api.FinalizableApiCollectorCommonArgs{
				UrlTemplate: "v2/alerts",
				Query: func(reqData *api.RequestData, createdAfter *time.Time) (url.Values, errors.Error) {
					query := url.Values{}

					query.Set("query", fmt.Sprintf("details.key:impacted-services+details.value:%s", data.Options.ServiceId))
					query.Set("sort", "createdAt")
					query.Set("order", "desc")
					query.Set("limit", fmt.Sprintf("%d", reqData.Pager.Size))
					query.Set("offset", fmt.Sprintf("%d", reqData.Pager.Skip))
					return query, nil
				},
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					rawResult := collectedAlerts{}
					err := api.UnmarshalResponse(res, &rawResult)
					return rawResult.Alerts, err
				},
			},
		},
		CollectUnfinishedDetails: &api.FinalizableApiCollectorDetailArgs{
			FinalizableApiCollectorCommonArgs: api.FinalizableApiCollectorCommonArgs{
				// 2. "Input" here is the type: simplifiedRawAlert which is the element type of the returned iterator from BuildInputIterator
				UrlTemplate: "v1/alerts/{{ .Input.Id }}",
				// 3. No custom query params/headers needed for this endpoint
				Query: nil,
				// 4. Parse the response for this endpoint call into a json.RawMessage
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					rawResult := collectedAlert{}
					err := api.UnmarshalResponse(res, &rawResult)
					return []json.RawMessage{rawResult.Alert}, err
				},
			},
			BuildInputIterator: func() (api.Iterator, errors.Error) {
				// 1. fetch individual "active/non-final" alerts from previous collections+extractions
				cursor, err := db.Cursor(
					dal.Select("id, created_date"),
					dal.From(&models.Alert{}),
					dal.Where(
						"service_id = ? AND connection_id = ? AND status != ?",
						data.Options.ServiceId, data.Options.ConnectionId, "resolved",
					),
				)
				if err != nil {
					return nil, err
				}
				return api.NewDalCursorIterator(db, cursor, reflect.TypeOf(simplifiedRawAlert{}))
			},
		},
	})
	if err != nil {
		return nil
	}
	return collector.Execute()
}
