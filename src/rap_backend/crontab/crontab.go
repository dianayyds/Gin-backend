package crontab

import (
	"rap_backend/service"

	"github.com/robfig/cron"
)

func CronInit() {
	service.RefreshCountrysLocalCache()
	service.RefreshLabelInfosLocalCache()
	service.RefreshRolePurviewCache()
	c := cron.New()
	c.AddFunc("0 */5 * * * *", service.RefreshCountrysLocalCache)
	c.AddFunc("0 */5 * * * *", service.RefreshLabelInfosLocalCache)
	c.AddFunc("0 */5 * * * *", service.RefreshRolePurviewCache)

	c.Start()
}
