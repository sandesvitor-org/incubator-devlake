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

package impl

import (
	"fmt"
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/subtaskmeta/sorter"
	"github.com/apache/incubator-devlake/plugins/pokemon/api"
	"github.com/apache/incubator-devlake/plugins/pokemon/models"
	"github.com/apache/incubator-devlake/plugins/pokemon/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/pokemon/tasks"
)

// make sure interface is implemented

var _ interface {
	plugin.PluginMeta
	plugin.PluginInit
	plugin.PluginModel
	plugin.PluginMigration
	plugin.PluginTask
	plugin.PluginApi
	plugin.DataSourcePluginBlueprintV200
	plugin.CloseablePluginTask
	plugin.PluginSource
} = (*Pokemon)(nil)

type Pokemon struct{}

func init() {
	// check subtask meta loop when init subtask meta
	if _, err := sorter.NewDependencySorter(tasks.SubTaskMetaList).Sort(); err != nil {
		panic(err)
	}
}

func (p Pokemon) Description() string {
	return "collect some Pokemon data"
}

// PkgPath information lost when compiled as plugin(.so)
func (p Pokemon) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/pokemon"
}

func (p Pokemon) Name() string {
	return "pokemon"
}

func (p Pokemon) Scope() plugin.ToolLayerScope {
	return &models.PokemonType{}
}

func (p Pokemon) ScopeConfig() dal.Tabler {
	return nil
}

func (p Pokemon) Init(basicRes context.BasicRes) errors.Error {
	api.Init(basicRes, p)

	return nil
}

func (p Pokemon) Connection() dal.Tabler {
	return &models.PokemonConnection{}
}

func (p Pokemon) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.PokemonConnection{},
		&models.PokemonType{},
		&models.PokemonTypeMove{},
	}
}

func (p Pokemon) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectPokemonTypeMovesStatefullMeta,
		tasks.ExtractPokemonTypeMovesMeta,
	}
}

func (p Pokemon) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
		p.Name(),
	)
	connection := &models.PokemonConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get Pokemon connection by the given connection ID")
	}
	var timeAfter *time.Time
	if op.TimeAfter != "" {
		convertedTime, err := errors.Convert01(time.Parse(time.RFC3339, op.TimeAfter))
		if err != nil {
			return nil, errors.BadInput.Wrap(err, fmt.Sprintf("invalid value for `timeAfter`: %s", timeAfter))
		}
		timeAfter = &convertedTime
	}
	client, err := helper.NewApiClient(taskCtx.GetContext(), connection.Endpoint, nil, 0, connection.Proxy, taskCtx)
	if err != nil {
		return nil, err
	}
	asyncClient, err := helper.CreateAsyncApiClient(taskCtx, client, nil)
	if err != nil {
		return nil, err
	}
	return &tasks.PokemonTaskData{
		Options:   op,
		TimeAfter: timeAfter,
		ApiClient: asyncClient,
	}, nil
}

func (p Pokemon) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (p Pokemon) Close(taskCtx plugin.TaskContext) errors.Error {
	_, ok := taskCtx.GetData().(*tasks.PokemonTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	return nil
}

func (p Pokemon) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"GET":    api.GetConnection,
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
		},
		"connections/:connectionId/remote-scopes": {
			"GET": api.RemoteScopes,
		},
		"connections/:connectionId/scopes": {
			"GET": api.GetScopeList,
			"PUT": api.PutScope,
		},
		"connections/:connectionId/scopes/:scopeId": {
			"GET":    api.GetScope,
			"PATCH":  api.UpdateScope,
			"DELETE": api.DeleteScope,
		},
	}
}

func (p Pokemon) MakeDataSourcePipelinePlanV200(connectionId uint64, scopes []*plugin.BlueprintScopeV200, syncPolicy plugin.BlueprintSyncPolicy,
) (plugin.PipelinePlan, []plugin.Scope, errors.Error) {
	return api.MakeDataSourcePipelinePlanV200(p.SubTaskMetas(), connectionId, scopes, &syncPolicy)
}
