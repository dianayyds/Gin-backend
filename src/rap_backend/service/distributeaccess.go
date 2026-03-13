package service

import (
	"rap_backend/dao"
	"time"
)

func CreateNewDistribute(taskId string, subTaskId string, distributor string) error {
	var taskDistribute dao.TaskDistribute
	taskDistribute.TaskId = taskId
	taskDistribute.SubtaskId = subTaskId
	taskDistribute.DistributeTime = time.Now()
	taskDistribute.TaskDistributor = distributor
	err := taskDistribute.CreateTaskDistribute()
	return err
}
