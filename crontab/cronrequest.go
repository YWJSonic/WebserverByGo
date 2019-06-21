package crontab

import (
	"time"

	"github.com/robfig/cron"
	"gitlab.com/WeberverByGo/data"
	"gitlab.com/WeberverByGo/foundation"
)

func init() {
	c = cron.New()
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

// SpecToTime  conver spec string to time
func SpecToTime(spec string) time.Time {
	target, _ := cron.Parse(data.MaintainStartTime)
	return target.Next(foundation.ServerNow())
}
