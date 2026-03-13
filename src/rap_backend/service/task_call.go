package service

import (
	"github.com/cihub/seelog"
	"rap_backend/dao"
)

// GetAnnotatedTaskCalls 批量更新callids状态
func UpdateCallIdsStatus(callids []dao.TaskCall, status string) error {
	var ids []int32
	for _, callid := range callids {
		ids = append(ids, callid.TcId)
	}

	upd := map[string]interface{}{
		"status": status,
	}
	_, err := dao.UpdateTaskCallByTcIDs(ids, upd)
	if err != nil {
		seelog.Errorf("UpdateTaskCallByTcIDs Operator data:%+v, err:%s", ids, err.Error())
		return err
	}
	return nil
}

// GetUserTaskCallsByStatus pageTab不传，则为task维度
func GetUserTaskCallsByStatus(taskId string, userId uint32, pageTab string, status []string) ([]dao.TaskCall, error) {
	callids, err := dao.GetUserTaskCallsByStatus(taskId, userId, pageTab, status)
	if err != nil {
		return nil, err
	}
	return *callids, nil
}

func UpdateTaskCallByTcIDs(ids []int32, upd map[string]interface{}) (int64, error) {
	return dao.UpdateTaskCallByTcIDs(ids, upd)
}
