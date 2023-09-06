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

package raw

import "time"

type Alert struct {
	Seen           *bool      `json:"seen"`
	Id             *string    `json:"id"`
	TinyId         *string    `json:"tinyId"`
	Alias          *string    `json:"alias"`
	Message        *string    `json:"message"`
	Status         *string    `json:"status"`
	Acknowledged   *bool      `json:"acknowledged"`
	IsSeen         *bool      `json:"isSeen"`
	Tags           *[]any     `json:"tags"`
	Snoozed        *bool      `json:"snoozed"`
	Count          *int       `json:"count"`
	LastOccurredAt *time.Time `json:"lastOccurredAt"`
	CreatedAt      *time.Time `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
	Source         *string    `json:"source"`
	Owner          *string    `json:"owner"`
	Priority       *string    `json:"priority"`
	Teams          []struct {
		Id string `json:"id"`
	} `json:"teams"`
	Responders []struct {
		Type string `json:"type"`
		Id   string `json:"id"`
	} `json:"responders"`
	Report struct {
		AckTime        *int    `json:"ackTime"`
		AcknowledgedBy *string `json:"acknowledgedBy"`
	} `json:"report"`
	OwnerTeamId string  `json:"ownerTeamId"`
	Actions     []any   `json:"actions"`
	Entity      *string `json:"entity"`
	Description *string `json:"description"`
	Details     struct {
		IncidentAlertType *string `json:"incident-alert-type"`
		IncidentId        *string `json:"incident-id"`
		ImpactedServices  *string `json:"impacted-services"`
	} `json:"details"`
}
