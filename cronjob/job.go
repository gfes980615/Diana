package cronjob

import (
	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/service"
	"github.com/robfig/cron"
)

var job *cronJob

func init() {
	job = &cronJob{
		c: cron.New(),
	}
	injection.AutoRegister(job)
}

type cronJob struct {
	c                    *cron.Cron
	mapleBulletinService service.MapleBulletinService `injection:"mapleBulletinService"`
	remoteService        service.RemoteService        `injection:"remoteService"`
}

func (cj *cronJob) MapleBulletinPushJob() {
	cj.c.AddFunc("0 */1 * * * ?", func() {
		cj.mapleBulletinService.PushToLine()
	})
}

// use cron to call medusa api, to get bulletin content
func (cj *cronJob) CallMedusaBulletinAPI() {
	cj.c.AddFunc("0 */2 * * * ?", func() {
		cj.remoteService.Simple("http://127.0.0.1:5000/maple/realtime_bulletin")
	})
}

func (cj *cronJob) Start() {
	cj.c.Start()
}

func InitJob() {
	job.MapleBulletinPushJob()
	job.CallMedusaBulletinAPI()
	job.Start()
}
