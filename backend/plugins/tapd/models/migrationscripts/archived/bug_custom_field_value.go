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

package archived

import (
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
)

type TapdBugCustomFieldValue struct {
	ConnectionId uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	WorkspaceId  uint64 `json:"workspace_id,string"`
	BugId        uint64 `gorm:"primaryKey;type:BIGINT  NOT NULL"`
	CustomField  string `gorm:"primaryKey;type:varchar(40) NOT NULL"` //TapdStoryCustomFields.custom_field
	Name         string `json:"name" gorm:"type:varchar(255)"`        // TapdStoryCustomFields.name
	CustomValue  string
	archived.NoPKModel
}

func (TapdBugCustomFieldValue) TableName() string {
	return "_tool_tapd_bug_custom_field_value"
}
