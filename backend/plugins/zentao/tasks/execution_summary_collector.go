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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_EXECUTION_SUMMARY_TABLE = "zentao_api_execution_summary"

var _ plugin.SubTaskEntryPoint = CollectExecutionSummary

func CollectExecutionSummary(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_EXECUTION_SUMMARY_TABLE,
		},
		ApiClient:   data.ApiClient,
		UrlTemplate: fmt.Sprintf("/projects/%d/executions", data.Options.ProjectId),
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("limit", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				Executions []json.RawMessage `json:"executions"`
			}
			err := api.UnmarshalResponse(res, &data)
			if errors.Is(err, api.ErrEmptyResponse) {
				return nil, nil
			}
			if err != nil {
				return nil, errors.Default.Wrap(err, "error reading endpoint response by Zentao bug collector")
			}
			return data.Executions, nil
		},
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}

var CollectExecutionSummaryMeta = plugin.SubTaskMeta{
	Name:             "collectExecutionSummary",
	EntryPoint:       CollectExecutionSummary,
	EnabledByDefault: true,
	Description:      "Collect Execution summary data from Zentao api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}
