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

package devops

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/domainlayer"
)

type CICDPipeline struct {
	domainlayer.DomainEntity
	Name         string `gorm:"type:varchar(255)"`
	Result       string `gorm:"type:varchar(100)"`
	Status       string `gorm:"type:varchar(100)"`
	Type         string `gorm:"type:varchar(100);comment: to indicate this is CI or CD"`
	DurationSec  uint64
	Environment  string `gorm:"type:varchar(255)"`
	CreatedDate  time.Time
	FinishedDate *time.Time
	CicdScopeId  string `gorm:"index;type:varchar(255)"`
}

func (CICDPipeline) TableName() string {
	return "cicd_pipelines"
}

// this is for the field `result` in table.cicd_pipelines and table.cicd_tasks
const (
	RESULT_SUCCESS = "SUCCESS"
	RESULT_FAILURE = "FAILURE"
	RESULT_ABORT   = "ABORT"
	RESULT_MANUAL  = "MANUAL"
	RESULT_SKIPPED = "SKIPPED"
)

// this is for the field `status` in table.cicd_pipelines and table.cicd_tasks
const (
	STATUS_IN_PROGRESS = "IN_PROGRESS"
	STATUS_DONE        = "DONE"
)

type ResultRule struct {
	Success []string
	Failed  []string
	Abort   []string
	Manual  []string
	Skipped []string
	Default string
}
type StatusRule[T comparable] struct {
	InProgress []T
	Done       []T
	Manual     []T
	Default    string
}

// GetResult compare the input with rule for return the enmu value of result
func GetResult(rule *ResultRule, input interface{}) string {
	for _, suc := range rule.Success {
		if suc == input {
			return RESULT_SUCCESS
		}
	}
	for _, fail := range rule.Failed {
		if fail == input {
			return RESULT_FAILURE
		}
	}
	for _, abort := range rule.Abort {
		if abort == input {
			return RESULT_ABORT
		}
	}
	for _, manual := range rule.Manual {
		if manual == input {
			return RESULT_MANUAL
		}
	}
	for _, skipped := range rule.Skipped {
		if skipped == input {
			return RESULT_SKIPPED
		}
	}
	return rule.Default
}

// GetStatus compare the input with rule for return the enmu value of status
func GetStatus[T comparable](rule *StatusRule[T], input T) string {
	for _, inp := range rule.InProgress {
		if inp == input {
			return STATUS_IN_PROGRESS
		}
	}
	for _, done := range rule.Done {
		if done == input {
			return STATUS_DONE
		}
	}
	for _, manual := range rule.Manual {
		if manual == input {
			return RESULT_MANUAL
		}
	}
	return rule.Default
}
