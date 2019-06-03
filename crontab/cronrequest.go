package crontab

import (
	"github.com/robfig/cron"
	"gitlab.com/WeberverByGo/data"
)

func init() {
	c = cron.New()

	c.AddFunc(data.MaintainStartTime, func() {
		data.Maintain = true
	})

	c.AddFunc(data.MaintainFinishTime, func() {
		data.Maintain = false
	})

	c.Start()
}

// CronStart start cron
func CronStart() {
	c.Start()
}

// CronStop stop cron
func CronStop() {
	c.Stop()
}

// NewCron add new fun
func NewCron(spec string, FUN func()) {
	c.AddFunc(spec, FUN)
}

// NewCronBaseJob cron job interface
func NewCronBaseJob(spec string, Job cron.Job) {
	c.AddJob(spec, Job)
}

// NewCronJob add new params job
func NewCronJob(spec string, Job *ParamsJob) {
	c.AddJob(spec, Job)
}
