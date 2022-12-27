package busi

import (
	"aggregate-task/internal/busi/core"
	"context"
	"sort"
	"time"

	log "github.com/sirupsen/logrus"
)

func AggregateTaskStart(ctx context.Context, finalityEpoch uint64, task string, dependTasks []string) {
	// self-inspection, find the missing epochs
	log.Infof("Self-Inspection: finding the missing epochs.")
	missingHeight, err := core.SelfInspection(ctx, task, dependTasks)
	if err != nil {
		return
	}
	sort.Slice(missingHeight, func(i, j int) bool { return missingHeight[i] < missingHeight[j] })

	// replay the missing epochs
	for idx, height := range missingHeight {
		// runtask
		log.Infof("Self-Inspection: replay[%v] height: %v", idx, int64(height))
		core.RunTask(ctx, task, int64(height))
	}
	log.Infof("Self-Inspection: completed successfully.")

	// cronjob, sync the incremental epochs
	for {
		select {
		case <-ctx.Done():
			log.Errorf("ctx done, receive signal: %s", ctx.Err().Error())
			return
		case <-time.After(time.Second * 30): // heartbeat
			log.Info("Ticktack: call syncIncrementalEpoch function.")
			core.SyncIncrementalEpoch(ctx, finalityEpoch, task, dependTasks)
		}
	}
}
