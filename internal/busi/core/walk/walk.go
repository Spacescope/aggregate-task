package walk

import (
	"aggregate-task/internal/busi/core/tasks"
	"aggregate-task/pkg/models/busi"
	"aggregate-task/pkg/utils"
	"context"
	"errors"
	"fmt"
	"sort"

	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type walkJob struct {
	jobIsRunning bool
	startTime    time.Time
	endTime      time.Time
}

type Walker struct {
	minHeight      uint64
	maxHeight      uint64
	task           string
	dependentTasks []string

	*walkJob
}

var insJob *walkJob
var once sync.Once

func NewWalker(minHeight, maxHeight uint64, task string, dependentTasks []string) *Walker {
	w := &Walker{
		minHeight:      minHeight,
		maxHeight:      maxHeight,
		task:           task,
		dependentTasks: dependentTasks,
	}

	once.Do(func() {
		insJob = &walkJob{}
	})

	w.walkJob = insJob

	return w
}

func (w *Walker) WalkChain(ctx context.Context, force bool) error {
	if w.jobIsRunning {
		str := fmt.Sprintf("The previous walk's job has begun at the time: %v, pls wait for it finishes or ctrl^c it.", w.startTime)
		log.Infof(str)
		return errors.New(str)
	} else {
		{
			w.jobIsRunning = true
			w.startTime = time.Now()

			log.Infof("Walk runs at time: %v, from: %v - to: %v", w.startTime, w.minHeight, w.maxHeight)
		}

		defer func() {
			w.jobIsRunning = false
			w.endTime = time.Now()

			log.Infof("Walk has finished the jobs: %v", w.endTime)
		}()

		exdb := utils.EngineGroup[utils.DBExtract]

		dependentTasksHeightSets := make([][]uint64, len(w.dependentTasks))
		for idx, dt := range w.dependentTasks {
			tsStates := make([]*busi.TipsetsState, 0)
			if err := exdb.Where("tipset between ? and ? and state = 1 and topic_name = ?", w.minHeight, w.maxHeight, dt).Asc("tipset").Find(&tsStates); err != nil {
				log.Errorf("execute sql error: %v", err)
				return err
			}

			for _, tsState := range tsStates {
				dependentTasksHeightSets[idx] = append(dependentTasksHeightSets[idx], tsState.Tipset)
			}
		}

		syncHeight := utils.Intersect(dependentTasksHeightSets...)
		sort.Slice(syncHeight, func(i, j int) bool { return syncHeight[i] < syncHeight[j] })

		for _, height := range syncHeight {
			log.Infof("Walk: replay height: %v", int64(height))
			RunTask(ctx, w.task, int64(height))
		}

		return nil
	}
}

func RunTask(ctx context.Context, task string, height int64) error {
	return tasks.GetTask(task).Run(ctx, height)
}
