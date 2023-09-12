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

const RAW_INCIDENTS_TABLE = "opsgenie_incidents"

var _ plugin.SubTaskEntryPoint = CollectIncidents

type (
	collectedIncidents struct {
		TotalCount int               `json:"totalCount"`
		Data       []json.RawMessage `json:"data"`
	}

	collectedIncident struct {
		Data json.RawMessage `json:"data"`
	}
	simplifiedRawIncident struct {
		Id        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
	}
)

var CollectIncidentsMeta = plugin.SubTaskMeta{
	Name:             "collectIncidents",
	EntryPoint:       CollectIncidents,
	EnabledByDefault: true,
	Description:      "Collect Opsgenie incidents",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func CollectIncidents(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*OpsgenieTaskData)
	db := taskCtx.GetDal()
	args := api.RawDataSubTaskArgs{
		Ctx:     taskCtx,
		Options: data.Options,
		Table:   RAW_INCIDENTS_TABLE,
	}
	collector, err := api.NewStatefulApiCollectorForFinalizableEntity(api.FinalizableApiCollectorArgs{
		RawDataSubTaskArgs: args,
		ApiClient:          data.Client,
		TimeAfter:          data.TimeAfter,
		CollectNewRecordsByList: api.FinalizableApiCollectorListArgs{
			PageSize: 100,
			FinalizableApiCollectorCommonArgs: api.FinalizableApiCollectorCommonArgs{
				UrlTemplate: "v1/incidents",
				Query: func(reqData *api.RequestData, createdAfter *time.Time) (url.Values, errors.Error) {
					query := url.Values{}

					query.Set("query", fmt.Sprintf("impactedServices:%s", data.Options.ServiceId))
					query.Set("sort", "createdAt")
					query.Set("order", "desc")
					query.Set("limit", fmt.Sprintf("%d", reqData.Pager.Size))
					query.Set("offset", fmt.Sprintf("%d", reqData.Pager.Skip))
					return query, nil
				},
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					rawResult := collectedIncidents{}
					err := api.UnmarshalResponse(res, &rawResult)

					return rawResult.Data, err
				},
			},
		},
		CollectUnfinishedDetails: &api.FinalizableApiCollectorDetailArgs{
			FinalizableApiCollectorCommonArgs: api.FinalizableApiCollectorCommonArgs{
				// 2. "Input" here is the type: simplifiedRawAlert which is the element type of the returned iterator from BuildInputIterator
				UrlTemplate: "v1/incidents/{{ .Input.Id }}",
				// 3. No custom query params/headers needed for this endpoint
				Query: nil,
				// 4. Parse the response for this endpoint call into a json.RawMessage
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					rawResult := collectedIncident{}
					err := api.UnmarshalResponse(res, &rawResult)

					return []json.RawMessage{rawResult.Data}, err
				},
			},
			BuildInputIterator: func() (api.Iterator, errors.Error) {
				// 1. fetch individual "active/non-final" incident from previous collections+extractions
				cursor, err := db.Cursor(
					dal.Select("id, created_at"),
					dal.From(&models.Incident{}),
					dal.Where(
						"service_id = ? AND connection_id = ? AND status != ?",
						data.Options.ServiceId, data.Options.ConnectionId, "closed",
					),
				)
				if err != nil {
					return nil, err
				}
				return api.NewDalCursorIterator(db, cursor, reflect.TypeOf(simplifiedRawIncident{}))
			},
		},
	})
	if err != nil {
		return nil
	}
	return collector.Execute()
}
