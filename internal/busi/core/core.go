package core

import (
	"aggregate-task/internal/busi/core/tasks"
	"aggregate-task/pkg/models/busi"
	"aggregate-task/pkg/utils"
	"context"
	"fmt"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func RunTask(ctx context.Context, task string, height int64) error {
	return tasks.GetTask(task).Run(ctx, height)
}

func SelfInspection(ctx context.Context, task string, dependTasks []string) ([]uint64, error) {
	exdb := utils.EngineGroup[utils.DBExtract]
	taskdb := utils.EngineGroup[utils.DBOBTask]

	// get the height both from indicate task and dependent tasks
	dependentTasksHeightSets := make([][]uint64, len(dependTasks))
	for idx, dt := range dependTasks {
		tsStates := make([]*busi.TipsetsState, 0)
		if err := exdb.Where("state = 1 and topic_name = ?", dt).Asc("tipset").Find(&tsStates); err != nil {
			log.Errorf("execute sql error: %v", err)
			return []uint64{}, err
		}

		for _, tsState := range tsStates {
			dependentTasksHeightSets[idx] = append(dependentTasksHeightSets[idx], tsState.Tipset)
		}
	}

	taskSlice := make([]*busi.TipsetWithVersion, 0)
	// if err := taskdb.SQL("select height, version from ?", task).Distinct("height").Asc("height").Find(&taskSlice); err != nil {
	if err := taskdb.Table(task).Cols("height", "version").Distinct("height").Asc("height").Find(&taskSlice); err != nil {
		log.Errorf("execute sql error: %v", err)
		return []uint64{}, err
	}
	var taskHeightSet []uint64
	for _, task := range taskSlice {
		taskHeightSet = append(taskHeightSet, task.Height)
	}

	missingHeight := utils.Except(utils.Intersect(dependentTasksHeightSets...), taskHeightSet)

	return missingHeight, nil
}

func SyncIncrementalEpoch(ctx context.Context, finalityEpoch uint64, task string, dependTasks []string) {
	exdb := utils.EngineGroup[utils.DBExtract]
	taskdb := utils.EngineGroup[utils.DBOBTask]

	sql := fmt.Sprintf("select max(height) as height from %v", task)
	result, err := taskdb.QueryString(sql)
	if err != nil {
		log.Errorf("execute sql error: %v", err)
		return
	}

	if len(result) > 0 {
		height, _ := strconv.ParseUint(result[0]["height"], 10, 64)
		if height >= finalityEpoch {
			height -= finalityEpoch
		}

		dependentTasksHeightSets := make([][]uint64, len(dependTasks))
		for idx, dt := range dependTasks {
			tsStates := make([]*busi.TipsetsState, 0)
			if err := exdb.Where("tipset > ? and state = 1 and topic_name = ?", height, dt).Asc("tipset").Find(&tsStates); err != nil {
				log.Errorf("execute sql error: %v", err)
				return
			}

			for _, tsState := range tsStates {
				dependentTasksHeightSets[idx] = append(dependentTasksHeightSets[idx], tsState.Tipset)
			}
		}

		// syncHeight := utils.FindIntersections(dependentTasksHeightSets)
		syncHeight := utils.Intersect(dependentTasksHeightSets...)

		for _, height := range syncHeight {
			log.Infof("Sync: replay height: %v", int64(height))
			RunTask(ctx, task, int64(height))
		}
	}
}
