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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

func init() {
	RegisterSubtaskMeta(&CollectPokemonTypeMovesStatefullMeta)
}

type PokemonTypeMovesResponse struct {
	Name  string            `json:"name"`
	Moves []json.RawMessage `json:"moves"`
}

const RAW_POKEMON_TYPE_MOVES_TABLE = "pokemon_type_moves"

var CollectPokemonTypeMovesStatefullMeta = plugin.SubTaskMeta{
	Name:             "collectPokemonTypeMovesStatefull",
	EntryPoint:       CollectPokemonTypeMovesStatefull,
	EnabledByDefault: true,
	Description:      "Collect Pokemon type moves (i.e. 'normal' type has 'slan' move)",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func CollectPokemonTypeMovesStatefull(taskCtx plugin.SubTaskContext) errors.Error {
	logger := taskCtx.GetLogger()
	logger.Info("collect pokemon type moves")
	data := taskCtx.GetData().(*PokemonTaskData)

	args := api.RawDataSubTaskArgs{
		Ctx:     taskCtx,
		Options: data.Options,
		Table:   RAW_POKEMON_TYPE_MOVES_TABLE,
	}
	collectorWithState, err := api.NewStatefulApiCollector(args, data.TimeAfter)
	if err != nil {
		return err
	}

	err = collectorWithState.InitCollector(api.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		UrlTemplate: "type/{{ .Params.ScopeId }}",
		Query:       nil,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			fmt.Println("Inside ResponseParser: body raw")

			rawResult := PokemonTypeMovesResponse{}

			err := api.UnmarshalResponse(res, &rawResult)
			if err != nil {
				return nil, err
			}
			return rawResult.Moves, nil
		},
	})
	if err != nil {
		logger.Error(err, "collect iteration error")
		return err
	}
	return collectorWithState.Execute()
}
