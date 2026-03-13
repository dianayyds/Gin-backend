package dao

import (
	"errors"
	"rap_backend/db"

	"gorm.io/gorm"
)

type TaskUserRelation struct {
	TurId    int32  `gorm:"primaryKey;AUTO_INCREMENT" json:"tur_id" from:"tur_id"`
	Tid      int32  `json:"tid" from:"tid"`
	TaskId   string `json:"task_id" from:"task_id"`
	UserId   uint32 `json:"user_id" form:"user_id"`
	TaskRole string `json:"task_role" form:"task_role"`
	LastId   int32  `json:"last_id" from:"last_id"`
	TotalNum int    `json:"total_num" from:"total_num"`
	DoneNum  int    `json:"done_num" from:"done_num"`
	Status   string `json:"status" form:"status"`
}

// 批量插入TaskUserRelation.
func CreateTaskUserRelationInBatches(tc []TaskUserRelation) error {
	result := db.GMysalDB.Model(&TaskUserRelation{}).CreateInBatches(tc, 100)
	return result.Error
}

// 根据主键id更新TaskUserRelation表数据.
func (u *TaskUserRelation) UpdateTaskUserRelationByID(columns ...interface{}) (int64, error) {
	if len(columns) == 0 {
		return 0, nil
	}
	upd := *u
	result := db.GMysalDB.Model(&u).Select(columns[0], columns[1:]...).Updates(upd)
	return result.RowsAffected, result.Error
}

// 根据taskid,uid,role,获取TaskUserRelation对应的内容
func GetTaskUserRelationInfo(taskId string, userID uint32, taskRole string) (*TaskUserRelation, error) {
	var rel = TaskUserRelation{}
	result := db.GMysalDB.Model(&TaskUserRelation{}).Where("user_id = ? AND task_role = ? AND task_id = ?", userID, taskRole, taskId).Take(&rel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &rel, nil
}

// 根据taskid,uid,role,更新TaskUserRelation状态等内容
func UpdateTaskUserRelationByWhere(taskID string, userID uint32, taskRole string, upd map[string]interface{}) (int64, error) {
	result := db.GMysalDB.Model(&TaskUserRelation{}).Where("user_id = ? and task_role=? and task_id = ?", userID, taskRole, taskID).Updates(upd)
	return result.RowsAffected, result.Error
}

//根据taskid，role，status, 获取任务下用户所属内容完成情况，用于判断任务是否进入下一阶段
// func GetTaskUserRelationByStatus(taskID, taskRole, status string) (int64, error) {
// 	var total int64
// 	result := db.GMysalDB.Model(TaskUserRelation{}).Where("task_id = ? and task_role=? and status = ?", taskID, taskRole, status).Count(&total)
// 	return total, result.Error
// }

// 获取taskid+role的userid各个状态的数量
func GetTaskUserRelationStatusNum(taskID, taskRole string) (map[string]int, error) {
	var list = make(map[string]int, 0)
	rows, err := db.GMysalDB.Model(&TaskUserRelation{}).Select("status, count(1) as num").Where("task_id = ? AND task_role = ?", taskID, taskRole).Group("status").Rows()
	if err != nil {
		return list, err
	}
	defer rows.Close()

	for rows.Next() {
		var st string
		var num int
		err = rows.Scan(&st, &num)
		list[st] = num
	}
	return list, nil
}

// 根据userid taskrole state获取用户都完成了哪些任务ids
func GetTaskIDsByUserRole(taskRole, status string, userID uint32) ([]int32, error) {
	var ids = make([]int32, 0)
	result := db.GMysalDB.Model(&TaskUserRelation{}).Select("tid").Where("user_id = ? and task_role = ? and status = ?", userID, taskRole, status).Find(&ids)
	return ids, result.Error
}

// 获取taskid+role的userid下完成的数量
func GetTaskUserRelationDoneNum(taskID, taskRole string) (*[]TaskUserRelation, error) {
	var list = make([]TaskUserRelation, 0)
	result := db.GMysalDB.Model(&TaskUserRelation{}).Where("task_id = ? AND task_role = ?", taskID, taskRole).Find(&list)
	if result.Error != nil {
		return &list, result.Error
	}
	return &list, nil
}

// 获取uid+role+taskids下的taskuserrelation信息
func GetTaskUserRelationByUIDRole(userID uint32, taskRole string, tids []int32) (*[]TaskUserRelation, error) {
	var list = make([]TaskUserRelation, 0)
	result := db.GMysalDB.Model(&TaskUserRelation{}).Where("user_id = ? AND task_role = ? and tid IN ?", userID, taskRole, tids).Find(&list)
	if result.Error != nil {
		return &list, result.Error
	}
	return &list, nil
}

func GetTaskUserRelationListByTaskId(taskID string) ([]TaskUserRelation, error) {
	var list = make([]TaskUserRelation, 0)
	err := db.GMysalDB.Model(&TaskUserRelation{}).
		Where("task_id = ? ", taskID).
		Order("tur_id asc").
		Find(&list).Error
	if err != nil {
		return list, err
	}
	return list, nil
}

func DeleteTaskUserRelationByTaskID(taskID string) error {
	return db.GMysalDB.
		Where("task_id = ?", taskID).
		Delete(&TaskUserRelation{}).
		Error
}

func DeleteTaskUserRelationByRole(taskID, role string) error {
	return db.GMysalDB.
		Where("task_id = ? AND task_role = ?", taskID, role).
		Delete(&TaskUserRelation{}).
		Error
}
