package crontab

import (
	"fmt"
	"time"

	"gitlab.com/WeberverByGo/data"
	cron "gitlab.com/WeberverByGo/robfig/cron.v3"
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
	target, err := cron.ParseStandard(data.MaintainStartTime)
	if err != nil {
		fmt.Println("Cron SpecToTime:", err)
	}
	return target.Next(time.Now())
}
