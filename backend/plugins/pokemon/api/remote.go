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

package api

import (
	"context"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/pokemon/models"
	"github.com/apache/incubator-devlake/plugins/pokemon/models/raw"
)

type RemoteScopesChild struct {
	Type     string      `json:"type"`
	ParentId *string     `json:"parentId"`
	Id       string      `json:"id"`
	Name     string      `json:"name"`
	Data     interface{} `json:"data"`
}

type RemoteScopesOutput struct {
	Children []RemoteScopesChild `json:"children"`
}

type SearchRemoteScopesOutput struct {
	Children []RemoteScopesChild `json:"children"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
}

type PokemonTypeResponse struct {
	Count    int               `json:"count"`
	Next     string            `json:"next"`
	Previous string            `json:"previous"`
	Results  []raw.PokemonType `json:"results"`
}

const RemoteScopesPerPage int = 100
const TypeScope string = "scope"

func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.PokemonConnection{}
	err := connectionHelper.First(connection, input.Params)

	if err != nil {
		return nil, err
	}

	// create api client
	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, connection)
	if err != nil {
		return nil, err
	}

	query := initialQuery()

	var res *http.Response
	outputBody := &RemoteScopesOutput{}
	res, err = apiClient.Get("/type", query, nil)

	if err != nil {
		return nil, err
	}

	response := &PokemonTypeResponse{}
	err = api.UnmarshalResponse(res, response)
	if err != nil {
		return nil, err
	}

	// append service to output
	for _, pokemonType := range response.Results {
		child := RemoteScopesChild{
			Type: TypeScope,
			Id:   pokemonType.Name,
			Name: pokemonType.Name,
			Data: models.PokemonType{
				Url:  pokemonType.Url,
				Name: pokemonType.Name,
				Id:   pokemonType.Name,
			},
		}
		outputBody.Children = append(outputBody.Children, child)
	}

	return &plugin.ApiResourceOutput{Body: outputBody, Status: http.StatusOK}, nil
}

func initialQuery() url.Values {
	query := url.Values{}
	return query
}
