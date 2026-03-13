package dao

import (
	"errors"
	"rap_backend/db"
	"rap_backend/utils"

	"gorm.io/gorm"
)

type TaskCall struct {
	TcId            int32  `gorm:"primaryKey;AUTO_INCREMENT" json:"tc_id" from:"tc_id"`
	TaskId          string `json:"task_id" from:"task_id"`
	CallId          string `json:"call_id" form:"call_id"`
	MultiAnnotator  string `json:"multi_annotator" form:"multi_annotator"`
	MultiAuditor    string `json:"multi_auditor" form:"multi_auditor"`
	MultiAnalysts   string `json:"multi_analysts" form:"multi_analysts"`
	ActualAnnotator string `json:"actual_annotator" form:"actual_annotator"`
	ActualAuditor   string `json:"actual_auditor" form:"actual_auditor"`
	ActualAnalysts  string `json:"actual_analysts" form:"actual_analysts"`
	SerialNumber    int    `json:"serial_number" from:"serial_number"`
	Status          string `json:"status" form:"status"`
	RejectReason    string `json:"reject_reason" form:"reject_reason"`
}

// 批量插入数据
func CreateTaskCallInBatches(tc []TaskCall) error {
	result := db.GMysalDB.Model(&TaskCall{}).CreateInBatches(tc, 100)
	return result.Error
}

// 根据taskID获取全部callids对应的主键ids，方便更新数据
func GetTcIDsByTaskID(taskID string) ([]int32, error) {
	var ids = make([]int32, 0)
	result := db.GMysalDB.Model(&TaskCall{}).Select("tc_id").Where("task_id = ?", taskID).Order("tc_id asc").Find(&ids)
	return ids, result.Error
}

// 根据主键id更新task calls表数据
func (u *TaskCall) UpdateTaskCallByTcID(columns ...interface{}) (int64, error) {
	if len(columns) == 0 {
		return 0, nil
	}
	upd := *u
	result := db.GMysalDB.Model(&u).Select(columns[0], columns[1:]...).Updates(upd)
	return result.RowsAffected, result.Error
}

// 根据主键ids更新task calls表数据
func UpdateTaskCallByTcIDs(tcIDs []int32, upd map[string]interface{}) (int64, error) {
	result := db.GMysalDB.Model(&TaskCall{}).Where("tc_id IN ?", tcIDs).Updates(upd)
	return result.RowsAffected, result.Error
}

// 获取taskid下callids ，分页+类型+状态
func GetTaskCallsByTaskId(taskID string, pageTab string, userID uint32, status []string, pageNum, pageSize int) (*[]TaskCall, int64, error) {
	var calls = make([]TaskCall, 0)
	var total int64
	query := db.GMysalDB.Model(&TaskCall{}).Where("task_id = ?", taskID)
	switch pageTab {
	case utils.TASK_PAGE_TAB_ANNOTAT:
		query.Where("multi_annotator = ?", userID)
	case utils.TASK_PAGE_TAB_AUDIT:
		query.Where("multi_auditor = ?", userID)
	case utils.TASK_PAGE_TAB_ANALYS:
		query.Where("multi_analysts = ?", userID)
	}
	if len(status) != 0 {
		query.Where("status in ?", status)
	}
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	if pageSize != 0 {
		query.Offset((pageNum - 1) * pageSize).Limit(pageSize)
	}
	result := query.Order("tc_id asc").Find(&calls)
	return &calls, total, result.Error
}

// 根据用户最后一次保持数据的编号ID，查询小于等于此编号，且属于此用户的总数
func GetTaskCallNumByLastID(taskID, pageTab string, userID uint32, lastID int32) (int64, error) {
	query := db.GMysalDB.Model(&TaskCall{}).Where("task_id = ? AND serial_number <= ?", taskID, lastID)
	switch pageTab {
	case utils.TASK_PAGE_TAB_ANNOTAT:
		query.Where("multi_annotator = ?", userID)
	case utils.TASK_PAGE_TAB_AUDIT:
		query.Where("multi_auditor = ?", userID)
	case utils.TASK_PAGE_TAB_ANALYS:
		query.Where("multi_analysts = ?", userID)
	}
	var total int64
	result := query.Count(&total)
	return total, result.Error
}

// 更改task calls 某个callid的状态，newstatus新状态，oldstatus是老状态，为空表示不判断当前状态
func UpdateTaskCallStatus(taskID, callID, newStatus, oldStatus string) (int64, error) {
	query := db.GMysalDB.Model(&TaskCall{}).Where("task_id = ? AND call_id = ?", taskID, callID)
	if oldStatus != "" {
		query.Where("status = ?", oldStatus)
	}
	result := query.Update("status", newStatus)
	return result.RowsAffected, result.Error
}

// 根据taskid,callid,更新TaskCall状态等内容
func UpdateTaskCallByWhere(taskID, callID string, upd map[string]interface{}) (int64, error) {
	result := db.GMysalDB.Model(&TaskCall{}).Where("task_id = ? and call_id = ?", taskID, callID).Updates(upd)
	return result.RowsAffected, result.Error
}

// 获取taskid+user_id的callids中各个状态的数量
func GetTaskCallStatusNum(taskID, pageTab string, userID uint32) (map[string]int, error) {
	var list = make(map[string]int, 0)
	query := db.GMysalDB.Model(&TaskCall{}).Select("status, count(1) as num").Where("task_id = ?", taskID)

	switch pageTab {
	case utils.TASK_PAGE_TAB_ANNOTAT:
		query.Where("multi_annotator = ?", userID)
	case utils.TASK_PAGE_TAB_AUDIT:
		query.Where("multi_auditor = ?", userID)
	case utils.TASK_PAGE_TAB_ANALYS:
		query.Where("multi_analysts = ?", userID)
	}
	rows, err := query.Group("status").Rows()

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

func GetTaskCallStatusNumByStatus(taskID, pageTab string, userID uint32, status []string) (int64, error) {
	query := db.GMysalDB.
		Model(&TaskCall{}).
		Select("status, count(1) as num").
		Where("task_id = ?", taskID).
		Where("status in ?", status)

	switch pageTab {
	case utils.TASK_PAGE_TAB_ANNOTAT:
		query.Where("multi_annotator = ?", userID)
	case utils.TASK_PAGE_TAB_AUDIT:
		query.Where("multi_auditor = ?", userID)
	case utils.TASK_PAGE_TAB_ANALYS:
		query.Where("multi_analysts = ?", userID)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

// 根据taskid callid 获取详情
func GetTaskCallInfoByTaskIDCallID(taskID, callID string) (*TaskCall, error) {
	var tc = TaskCall{}
	result := db.GMysalDB.Model(&TaskCall{}).Where("task_id = ? and call_id = ?", taskID, callID).Take(&tc)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &tc, nil
}

// 根据taskID获取全部callids
func GetCallIDsByTaskID(taskID string) ([]string, error) {
	var ids = make([]string, 0)
	result := db.GMysalDB.Model(&TaskCall{}).Select("call_id").Where("task_id = ?", taskID).Order("tc_id asc").Find(&ids)
	return ids, result.Error
}

// GetUserTaskCallsByStatus   pageTab不传，则为task维度
func GetUserTaskCallsByStatus(taskID string, userID uint32, pageTab string, status []string) (*[]TaskCall, error) {
	var calls = make([]TaskCall, 0)
	query := db.GMysalDB.Model(&TaskCall{}).Where("task_id = ?", taskID)
	switch pageTab {
	case utils.TASK_PAGE_TAB_ANNOTAT:
		query.Where("multi_annotator = ?", userID)
	case utils.TASK_PAGE_TAB_AUDIT:
		query.Where("multi_auditor = ?", userID)
	case utils.TASK_PAGE_TAB_ANALYS:
		query.Where("multi_analysts = ?", userID)
	}
	if len(status) != 0 {
		query.Where("status in ?", status)
	}
	result := query.Order("tc_id asc").Find(&calls)
	return &calls, result.Error
}

type TaskProgressDetail struct {
	Status string `json:"status"`
	Num    int    `json:"num"`
}

func GetTaskProgressDetails(taskID string) ([]*TaskProgressDetail, error) {
	var calls = make([]*TaskProgressDetail, 0)
	query := db.GMysalDB.Model(&TaskCall{}).
		Where("task_id = ?", taskID)
	result := query.
		Select("status, count(1) as num").
		Group("status").
		Find(&calls)
	return calls, result.Error
}
