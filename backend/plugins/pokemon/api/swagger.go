package api

import (
	"github.com/apache/incubator-devlake/plugins/pagerduty/tasks"
)

type PagerDutyTaskOptions tasks.PagerDutyOptions

// @Summary pagerduty task options for pipelines
// @Description This is a dummy API to demonstrate the available task options for pagerduty pipelines
// @Tags plugins/pagerduty
// @Accept application/json
// @Param pipeline body PagerDutyTaskOptions true "json"
// @Router /pipelines/pagerduty/pipeline-task [post]
func _() {}
