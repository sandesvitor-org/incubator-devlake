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
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/pokemon/models"
)

type PokemonApiParams struct {
}

type PokemonOptions struct {
	ConnectionId    uint64   `json:"connectionId"`
	TimeAfter       string   `json:"time_after,omitempty"`
	PokemonTypeName string   `json:"pokemon_type_name,omitempty"`
	PokemonTypeId   string   `json:"pokemon_type_id,omitempty"`
	Tasks           []string `json:"tasks,omitempty"`
	*models.PokemonScopeConfig
}

type PokemonTaskData struct {
	Options   *PokemonOptions
	TimeAfter *time.Time
	ApiClient api.RateLimitedApiClient
}

func (p *PokemonOptions) GetParams() any {
	return models.PokemonParams{
		ConnectionId: p.ConnectionId,
		ScopeId:      p.PokemonTypeId,
	}
}

func DecodeAndValidateTaskOptions(options map[string]interface{}) (*PokemonOptions, errors.Error) {
	op, err := DecodeTaskOptions(options)
	if err != nil {
		return nil, err
	}
	err = ValidateTaskOptions(op)
	if err != nil {
		return nil, err
	}
	return op, nil
}

func DecodeTaskOptions(options map[string]interface{}) (*PokemonOptions, errors.Error) {
	var op PokemonOptions
	err := api.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func EncodeTaskOptions(op *PokemonOptions) (map[string]interface{}, errors.Error) {
	var result map[string]interface{}
	err := api.Decode(op, &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func ValidateTaskOptions(op *PokemonOptions) errors.Error {
	if op.PokemonTypeName == "" {
		return errors.BadInput.New("not enough info for pokemon execution")
	}
	if op.PokemonTypeId == "" {
		return errors.BadInput.New("not enough info for pokemon execution")
	}
	// find the needed GitHub now
	if op.ConnectionId == 0 {
		return errors.BadInput.New("connectionId is invalid")
	}
	return nil
}
