package crontab

import (
	"time"

	cron "gitlab.com/ServerUtility/cron.v3"
	"gitlab.com/ServerUtility/crontabinfo"
	"gitlab.com/ServerUtility/foundation"
	"gitlab.com/WeberverByGo/serversetting"
)

var c *cron.Cron

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
func NewCronJob(spec string, job *crontabinfo.ParamsJob) {
	c.AddJob(spec, job)
}

// SpecToTime  conver spec string to time
func SpecToTime(spec string) time.Time {
	target, _ := cron.ParseStandard(serversetting.MaintainStartTime)
	return target.Next(foundation.ServerNow())
}

// NewLogCrontab new LogCrontab struct
func NewLogCrontab(Params func() string, FUN func(string)) *crontabinfo.LogCrontab {
	return &crontabinfo.LogCrontab{
		Params: func() string { return foundation.ServerNow().AddDate(0, 0, 1).Format("20060102") },
		FUN:    FUN,
	}
}
