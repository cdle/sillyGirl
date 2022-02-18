package core

import cron "github.com/robfig/cron/v3"

var c *cron.Cron

func init() {
	c = cron.New()
	c.Start()
}
