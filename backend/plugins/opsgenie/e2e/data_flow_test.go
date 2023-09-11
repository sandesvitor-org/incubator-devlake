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

package e2e

import (
	"fmt"
	"testing"

	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/opsgenie/impl"
	"github.com/apache/incubator-devlake/plugins/opsgenie/models"
	"github.com/apache/incubator-devlake/plugins/opsgenie/tasks"
	"github.com/stretchr/testify/require"
)

func TestDataFlow(t *testing.T) {
	var plugin impl.Opsgenie
	dataflowTester := e2ehelper.NewDataFlowTester(t, "opsgenie", plugin)
	options := tasks.OpsgenieOptions{
		ConnectionId: 1,
		ServiceId:    "5faa5cc2-37d8-4182-b441-dea974ca0b67",
		ServiceName:  "TestService",
		Tasks:        nil,
	}
	taskData := &tasks.OpsgenieTaskData{
		Options: &options,
	}

	dataflowTester.FlushTabler(&models.Service{})

	service := models.Service{
		ConnectionId: options.ConnectionId,
		Url:          fmt.Sprintf("https://sandesvitor.app.opsgenie.com/service/%s", options.ServiceId),
		Id:           options.ServiceId,
		Name:         options.ServiceName,
	}
	// scope
	require.NoError(t, dataflowTester.Dal.CreateOrUpdate(&service))

	// import raw data table
	dataflowTester.ImportCsvIntoRawTable("./raw_tables/_raw_opsgenie_incidents.csv", "_raw_opsgenie_incidents")

	dataflowTester.FlushTabler(&models.Incident{})

	// verify service converter
	dataflowTester.Subtask(tasks.ConvertServicesMeta, taskData)
	// dataflowTester.VerifyTableWithOptions(
	// 	models.Alert{},
	// 	e2ehelper.TableOptions{
	// 		CSVRelPath:  "./snapshot_tables/_tool_opsgenie_alerts.csv",
	// 		IgnoreTypes: []any{common.Model{}},
	// 	},
	// )

	// // verify service converter
	// dataflowTester.Subtask(tasks.ExtractAlertsMeta, taskData)
	// dataflowTester.VerifyTableWithOptions(
	// 	models.Alert{},
	// 	e2ehelper.TableOptions{
	// 		CSVRelPath:  "./snapshot_tables/_tool_opsgenie_alerts.csv",
	// 		IgnoreTypes: []any{common.Model{}},
	// 	},
	// )
}
