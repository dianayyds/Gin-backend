package dao

import (
	"rap_backend/db"
	"time"

	"github.com/cihub/seelog"
)

type TaskDistribute struct {
	TaskId          string    `json:"task_id" from:"task_id"`
	SubtaskId       string    `json:"subtask_id" from:"subtask_id"`
	TaskDistributor string    `json:"task_distributor" form:"task_distributor"`
	DistributeTime  time.Time `json:"distribute_time" form:"distribute_time"`
}

func (t *TaskDistribute) CreateTaskDistribute() error {
	result := db.GMysalDB.Create(t)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetDistributedSubTaskByTaskid(taskId string) (*[]TaskDistribute, error) {
	var distributeList []TaskDistribute
	result := db.GMysalDB.Model(&TaskDistribute{}).Where("task_id = ?", taskId).Find(&distributeList)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected != 0 {
		return &distributeList, nil
	} else {
		return nil, nil
	}
}

func GetTaskDistributedBySubTask(subtaskId string) (*TaskDistribute, error) {
	var distributeList TaskDistribute
	result := db.GMysalDB.Model(&TaskDistribute{}).Where("subtask_id = ?", subtaskId).Find(&distributeList)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected != 0 {
		return &distributeList, nil
	} else {
		return nil, nil
	}
}

func DeleteTaskInDistribute(taskId string) error {
	result := db.GMysalDB.Exec("delete from task_distributes where task_id = ?", taskId)
	if result.Error != nil {
		seelog.Errorf("delete %s task in distribute failed(del):%s", taskId, result.Error)
		return result.Error
	}
	return nil
}
