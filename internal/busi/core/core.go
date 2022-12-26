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

func RunTask(ctx context.Context, task string, height uint64) error {
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

	// {
	// 	x := utils.Intersect(dependentTasksHeightSets...)
	// 	_ = x

	// 	intersect_x1 := utils.Intersect(dependentTasksHeightSets[0], dependentTasksHeightSets[1])
	// 	_ = intersect_x1

	// 	except_x1_1 := utils.Except(dependentTasksHeightSets[0], dependentTasksHeightSets[1])
	// 	_ = except_x1_1

	// 	except_x1_2 := utils.Except(dependentTasksHeightSets[1], dependentTasksHeightSets[0])
	// 	_ = except_x1_2

	// 	x2 := utils.Intersect(dependentTasksHeightSets[1], dependentTasksHeightSets[2])
	// 	_ = x2

	// 	a := utils.Except(dependentTasksHeightSets[0], x)
	// 	b := utils.Except(dependentTasksHeightSets[1], x)
	// 	c := utils.Except(dependentTasksHeightSets[2], x)

	// 	_ = a
	// 	_ = b
	// 	_ = c

	// 	d := utils.Except([]uint64{9, 1, 2, 3, 7}, []uint64{1, 2, 3})
	// 	_ = d

	// 	fmt.Print("he")
	// }

	// 1) get the intersection set of dependent tasks
	// 2) get the difference set of task and dependent tasks
	// missingHeight := utils.SliceDifferenceInt(utils.FindIntersections(dependentTasksHeightSets), taskHeightSet)
	missingHeight := utils.Except(utils.Intersect(dependentTasksHeightSets...), taskHeightSet)

	return missingHeight, nil
}

func SyncIncrementalEpoch(ctx context.Context, task string, dependTasks []string) {
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
			log.Infof("Sync: replay height: %v", height)
			RunTask(ctx, task, height)
		}
	}
}
