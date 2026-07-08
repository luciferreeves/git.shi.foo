package jobs

import (
	"context"

	"git.shi.foo/models"
)

type Runner func(runContext context.Context, record *models.Job, report Reporter) error

type Reporter interface {
	Progress(phase string, percent int)
	Final() bool
}

type Progress struct {
	JobID   uint   `json:"job_id"`
	Status  string `json:"status"`
	Phase   string `json:"phase"`
	Percent int    `json:"percent"`
	Done    bool   `json:"done"`
}
