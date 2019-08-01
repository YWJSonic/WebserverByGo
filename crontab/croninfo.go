package crontab

import (
	cron "gitlab.com/WeberverByGo/robfig/cron.v3"
)

// 字段名				是否必须	允许的值			允许的特定字符
// 秒(Seconds)			是			0-59				* / , -
// 分(Minutes)			是			0-59				* / , -
// 时(Hours)			是			0-23				* / , -
// 日(Day of month)		是			1-31				* / , – ?
// 月(Month)			是			1-12 or JAN-DEC		* / , -
// 星期(Day of week)	否			0-6 or SUM-SAT		* / , – ?

var c *cron.Cron

// ParamsJob Cron job struct attach params
type ParamsJob struct {
	Params []interface{}
	FUN    func(...interface{})
}

// Run Cron Job interface
func (p *ParamsJob) Run() {
	if p.FUN != nil {
		p.FUN(p.Params...)
	}
}

// LogCrontab crontab job interface
type LogCrontab struct {
	Params func() string
	FUN    func(string)
}

// Run Cron Job interface
func (log *LogCrontab) Run() {
	if log.FUN != nil {
		tablename := log.Params()
		log.FUN(tablename)
	}
}
