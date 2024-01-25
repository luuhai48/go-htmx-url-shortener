package workers

import (
	"log"
	"log/slog"
	"luuhai48/short/models"

	"github.com/mborders/artifex"
)

type CronJob struct {
	dispatcher   *artifex.Dispatcher
	dispatchCron *artifex.DispatchCron
}

func (c *CronJob) Stop() {
	c.dispatchCron.Stop()
	c.dispatcher.Stop()
}

func CreateCleanupSessionCronjob() *CronJob {
	cleanupSession := artifex.NewDispatcher(1, 1)
	cleanupSession.Start()

	dc, err := cleanupSession.DispatchCron(func() {
		slog.Info("Cleaning up old sessions")
		if errr := models.DeleteOldSessions(); errr != nil {
			slog.Error(errr.Error())
		}
	}, "0 0 * * * *")
	if err != nil {
		log.Fatal(err)
	}

	return &CronJob{
		dispatcher:   cleanupSession,
		dispatchCron: dc,
	}
}
