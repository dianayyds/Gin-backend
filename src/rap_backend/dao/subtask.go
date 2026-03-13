package dao

import (
	"rap_backend/db"
	"time"

	"github.com/cihub/seelog"
)

type SubtaskInfo struct {
	Id                int32     `gorm:"AUTO_INCREMENT" json:"id" from:"id"`
	SubtaskId         string    `json:"subtask_id" from:"subtask_id"`
	SubtaskOperator   string    `json:"subtask_operator" form:"subtask_operator"`
	StartTime         time.Time `json:"start_time" form:"start_time"`
	FinishTime        time.Time `json:"finish_time" form:"finish_time"`
	SubtaskCreatetime time.Time `json:"subtask_createtime" form:"subtask_createtime"`
	Status            string    `json:"status" from:"status"`
}

func (t *SubtaskInfo) AddSubTaskInfo() error {
	result := db.GMysalDB.Create(t)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetSubTaskInfoBySubTaskid(subtaskId string) (*SubtaskInfo, error) {
	var subTask SubtaskInfo
	result := db.GMysalDB.Model(&SubtaskInfo{}).Where("subtask_id = ?", subtaskId).Find(&subTask)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected != 0 {
		return &subTask, nil
	} else {
		return nil, nil
	}
}

func UpdateSubTaskStatus(subtaskId string, status string) error {
	sqlstr := "update subtask_infos set status = '" + status + "' where subtask_id = '" + subtaskId + "'"
	result := db.GMysalDB.Exec(sqlstr)
	if result.Error != nil {
		seelog.Errorf("update sub task status error:%s", result.Error.Error())
		return result.Error
	}
	return nil
}

func DeleteSubtaskInSubTaskInfos(subtaskId string) error {
	result := db.GMysalDB.Exec("delete from subtask_infos where subtask_id = ?", subtaskId)
	if result.Error != nil {
		seelog.Errorf("delete %s sub task in sub task infos failed(del):%s", subtaskId, result.Error)
		return result.Error
	}
	return nil
}
