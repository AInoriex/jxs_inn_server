package scheduler

import (
	"eshop_server/src/cronjob/handler"
	"eshop_server/src/utils/log"
	"github.com/robfig/cron/v3"
)

var (
	Schedu = NewScheduler()
)

func InitScheduler() {
	// Schedu.AddJob("0 */10 * * * *", EveryTenMinuteTask) // 每10分钟执行一次
	// Schedu.AddJob("0 */2 * * * *", EveryTwoMinuteTask)  // 每2分钟执行一次
	// Schedu.AddJob("0 * * * * *", EveryMinuteTask)       // 每分钟执行一次
	// Schedu.AddJob("@every 10s", EveryTenSecondTask)     // 每10s执行一次
	// Schedu.AddJob("@every 5s", EveryFiveSecondTask)     // 每5s执行一次
	// Schedu.AddJob("@every 2s", EveryTwoSecondTask)      // 每2s执行一次

	Schedu.AddJob("@every 30s", handler.YltLoginCronjob)	// 每30s执行一次
	Schedu.AddJob("@every 5s", handler.UpdateOrderCronjob)	// 每5s执行一次
	Schedu.Start()
}

type Scheduler struct {
	cron *cron.Cron
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		cron.New(cron.WithSeconds()),
	}
}

func (s *Scheduler) AddJob(expression string, cmd func()) {
	s.cron.AddFunc(expression, cmd)
}

func (s *Scheduler) Start() {
	log.Debug("开始执行定时器")
	s.cron.Start()
}
