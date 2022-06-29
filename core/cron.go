package core

import cron "github.com/robfig/cron/v3"

var C *cron.Cron

func init() {
	C = cron.New()
	C.Start()
}
