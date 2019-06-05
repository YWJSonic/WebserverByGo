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
func NewCron(spec string, fun func()) {
	c.AddFunc(spec, fun)
}

// NewCronBaseJob cron job interface
func NewCronBaseJob(spec string, job cron.Job) {
	c.AddJob(spec, job)
}

// NewCronJob add new params job
func NewCronJob(spec string, job *ParamsJob) {
	c.AddJob(spec, job)
}
