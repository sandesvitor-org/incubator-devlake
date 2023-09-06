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
	"github.com/apache/incubator-devlake/plugins/pokemon/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractPokemonTypeMovesMeta)
}

type PokemonTypeMovesRawData struct {
	Name  string `json:"name"`
	Moves struct {
		Name string `json:"name"`
	} `json:"moves"`
}

var ExtractPokemonTypeMovesMeta = plugin.SubTaskMeta{
	Name:             "extractPokemonTypeMoves",
	EntryPoint:       ExtractPokemonTypeMoves,
	EnabledByDefault: true,
	Description:      "Extract Pokemon type moves",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ExtractPokemonTypeMoves(taskCtx plugin.SubTaskContext) errors.Error {
	logger := taskCtx.GetLogger()
	logger.Info("extract pokemon type moves")
	data := taskCtx.GetData().(*PokemonTaskData)

	args := api.RawDataSubTaskArgs{
		Ctx:     taskCtx,
		Options: data.Options,
		Table:   RAW_POKEMON_TYPE_MOVES_TABLE,
	}

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: args,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			var userRes models.PokemonTypeMove
			err := errors.Convert(json.Unmarshal(row.Data, &userRes))

			fmt.Println(userRes)

			if err != nil {
				return nil, err
			}

			results := make([]interface{}, 0)
			PokemonMove := &models.PokemonTypeMove{
				ConnectionId:  data.Options.ConnectionId,
				Name:          userRes.Name,
				PokemonTypeId: data.Options.PokemonTypeId,
			}
			results = append(results, PokemonMove)

			fmt.Println(results)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	err = extractor.Execute()
	if err != nil {
		return err
	}

	return nil
}
