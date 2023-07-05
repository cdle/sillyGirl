package core

import cron "github.com/robfig/cron/v3"

var CRON *cron.Cron

func init() {
	CRON = cron.New(cron.WithSeconds())
	CRON.Start()
}
