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

package models

import (
	"time"

	"github.com/apache/incubator-devlake/core/models/common"
)

const (
	AlertStatusAcknowledged AlertStatus = "resolved"
	AlertStatusTriggered    AlertStatus = "open"
	AlertStatusResolved     AlertStatus = "closed"
)

type (
	AlertPriority string
	AlertStatus   string

	Alert struct {
		common.NoPKModel
		ConnectionId   uint64 `gorm:"primaryKey"`
		Id             string `gorm:"primaryKey"`
		Seen           bool
		Snoozed        bool
		ServiceId      string
		ServiceName    string
		Description    string
		Message        string
		Owner          string
		AcknowledgedBy string
		AckTime        int
		CloseTime      int
		ClosedBy       string
		Priority       AlertPriority
		Status         AlertStatus
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}
)

func (Alert) TableName() string {
	return "_tool_opsgenie_alerts"
}
