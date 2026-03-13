package dao

import (
	"errors"
	"fmt"
	"rap_backend/db"
	"rap_backend/rpc/bigdata"
	"rap_backend/utils"
	"time"

	"github.com/cihub/seelog"
	"gorm.io/gorm"
)

type TaskInfo struct {
	Id             int32     `gorm:"primaryKey;AUTO_INCREMENT" json:"id" from:"id"`
	TaskId         string    `json:"task_id" from:"task_id"`
	TaskName       string    `json:"task_name" form:"task_name"`
	TaskCreator    string    `json:"task_creator" form:"task_creator"`
	TaskCreatetime time.Time `json:"task_createtime" form:"task_createtime"`
	StartTime      time.Time `json:"start_time" form:"start_time"`
	FinishTime     time.Time `json:"finish_time" form:"finish_time"`
	CallNum        int       `json:"call_num" form:"call_num"`
	Status         string    `json:"status" form:"status"`
	CreatorId      int       `json:"creator_id" form:"creator_id"`
	MultiAllocator string    `json:"multi_allocator" form:"multi_allocator"`
	MultiAnnotator string    `json:"multi_annotator" form:"multi_annotator"`
	MultiAuditor   string    `json:"multi_auditor" form:"multi_auditor"`
	MultiAnalysts  string    `json:"multi_analysts" form:"multi_analysts"`
	Labels         string    `json:"labels" form:"labels"`
}

type TaskListFilter struct {
	PageTab     string
	StatusList  []string
	UserIDs     []uint32
	SearchType  string
	SearchValue string
	PageNum     int
	PageSize    int
	TIDs        []int32
}

type TaskAndAllocatInfo struct {
	*TaskInfo
	Annotator []DoTaskUserNum `json:"annotator" form:"annotator"`
	Auditor   []DoTaskUserNum `json:"auditor" form:"auditor"`
	Analysts  []DoTaskUserNum `json:"analysts" form:"analysts"`
}

func (t *TaskInfo) AddTaskInfo() error {
	result := db.GMysalDB.Create(t)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetTaskListCount() (int64, error) {
	var cnt int64
	result := db.GMysalDB.Model(&TaskInfo{}).Count(&cnt)
	if result.Error != nil {
		return 0, result.Error
	}
	return cnt, nil
}

func GetTaskInfoList(searchType string, searchValue string, pageNum int, pageSize int) (*[]TaskInfo, error) {
	var taskList []TaskInfo
	var result *gorm.DB
	if searchType == utils.TASK_SEARCH_TYPE_ID {
		result = db.GMysalDB.Model(&TaskInfo{}).Where(fmt.Sprintf("task_id like '%%%s%%'", searchValue)).Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("task_createtime desc").Find(&taskList)
	} else if searchType == utils.TASK_SEARCH_TYPE_CREATOR {
		result = db.GMysalDB.Model(&TaskInfo{}).Where(fmt.Sprintf("task_creator like '%%%s%%'", searchValue)).Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("task_createtime desc").Find(&taskList)
	} else if searchType == utils.TASK_SEARCH_TYPE_NAME {
		result = db.GMysalDB.Model(&TaskInfo{}).Where(fmt.Sprintf("task_name like '%%%s%%'", searchValue)).Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("task_createtime desc").Find(&taskList)
	} else if searchType == utils.TASK_SEARCH_TYPE_STATUS {
		result = db.GMysalDB.Model(&TaskInfo{}).Where("status = ?", searchValue).Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("task_createtime desc").Find(&taskList)
	} else {
		result = db.GMysalDB.Model(&TaskInfo{}).Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("task_createtime desc").Find(&taskList)
	}
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected != 0 {
		return &taskList, nil
	} else {
		return nil, nil
	}
}

//func GetTaskInfoList(pageNum int, pageSize int) (*[]TaskInfo, error){
//	var taskList []TaskInfo
//	result := db.GMysalDB.Model(&TaskInfo{}).Limit(pageSize).Offset((pageNum-1) * pageSize).Order("task_createtime desc").Find(&taskList)
//	if result.Error != nil {
//		return nil, result.Error
//	} else if result.RowsAffected != 0 {
//		return &taskList, nil
//	} else {
//		return nil, nil
//	}
//}

func GetLabelValueFromGaussInfo(callid string) (*[]map[string]interface{}, error) {
	var results []map[string]interface{}
	ret := db.GGaussMysalDB.Table("gauss_callinfo").Where("CallID = ?", callid).Find(&results)
	if ret.Error != nil {
		seelog.Errorf("get %s default label value failed: %s", callid, ret.Error.Error())
		return nil, ret.Error
	}
	return &results, nil
}

//func GetLabelValueFromGaussInfo(callid string) (*[]string, error){
//	//colstr := ""
//	//for _, name := range(labelName){
//	//	colstr = colstr + name
//	//}
//	//colstr = colstr[0:len(colstr)-1]
//
//	sqlstr := "select * from gauss_callinfo where CallID = '" + callid + "'"
//
//	seelog.Infof("sqlstr: %s", sqlstr)
//	rows, err := db.GGaussMysalDB.Raw(sqlstr).Rows()
//	if err != nil {
//		seelog.Errorf("get %s default label value failed: %s", callid, err.Error())
//		return nil, err
//	}
//	defer rows.Close()
//	var valueList []string
//	for rows.Next() {
//		collist, _ := rows.Columns()
//		seelog.Infof("col list %s", collist)
//		valueList = make([]string, len(collist))
//		rows.Scan(&valueList)
//		seelog.Infof("value list %s", valueList)
//	}
//	return &valueList, nil
//}

func CheckTaskCallId2(callIdList []string) ([]string, []string, error) {
	var existedCallid []string
	var notExistedCallid []string

	callMap, err := bigdata.GetSvc().GetCallInfoMap(bigdata.GetCallInfoReq{
		CallIds: callIdList,
	})
	if err != nil {
		return nil, nil, err
	}

	for _, callId := range callIdList {
		if _, ok := callMap[callId]; ok {
			existedCallid = append(existedCallid, callId)
		} else {
			notExistedCallid = append(notExistedCallid, callId)
		}
	}

	return existedCallid, notExistedCallid, nil
}

func CheckTaskCallid(callidlist []string) ([]string, []string, error) {
	var existedCallid []string
	var notExistedCallid []string
	var strCallidList = ""
	for _, callid := range callidlist {
		strCallidList = strCallidList + "'" + callid + "',"
	}
	strCallidList = strCallidList[0 : len(strCallidList)-1]

	sqlstr := "select CallId from gauss_callinfo where CallID in (" + strCallidList + ")"
	seelog.Infof("sqlstr: %s", sqlstr)
	rows, err := db.GGaussMysalDB.Raw(sqlstr).Rows()
	if err != nil {
		seelog.Errorf("check gauss callid failed: %s", err.Error())
		return existedCallid, notExistedCallid, err
	}
	defer rows.Close()
	for rows.Next() {
		var extCallid string
		rows.Scan(&extCallid)
		existedCallid = append(existedCallid, extCallid)
	}
	existMap := make(map[string]bool, 0)
	for _, checkCallid := range callidlist {
		findflg := false
		for _, v := range existedCallid {
			if checkCallid == v {
				findflg = true
				break
			}
		}
		if !findflg {
			notExistedCallid = append(notExistedCallid, checkCallid)
			continue
		}
		if _, ok := existMap[checkCallid]; ok {
			notExistedCallid = append(notExistedCallid, checkCallid)
			continue
		}
		existMap[checkCallid] = true
	}
	seelog.Infof("ext:%s", existedCallid)
	seelog.Infof("notext:%s", notExistedCallid)
	return existedCallid, notExistedCallid, nil
}

// func UpdateTaskStatus(taskid string, status string) error {
// 	sqlstr := "update task_infos set status = '" + status + "' where task_id = '" + taskid + "'"
// 	result := db.GMysalDB.Exec(sqlstr)
// 	if result.Error != nil {
// 		seelog.Errorf("update task status error:%s", result.Error.Error())
// 		return result.Error
// 	}

// 	disList, err2 := GetDistributedSubTaskByTaskid(taskid)
// 	if err2 != nil {
// 		seelog.Errorf("update sub task status error %s", err2.Error())
// 		return err2
// 	}
// 	if disList == nil {
// 		seelog.Errorf("this task has no subtask task id :%s", taskid)
// 		return errors.New("this task has no subtask task id:" + taskid)
// 	}
// 	for _, sub := range *disList {
// 		err3 := UpdateSubTaskStatus(sub.SubtaskId, status)
// 		if err3 != nil {
// 			seelog.Errorf("update sub task %s status %s failed:%s", sub.SubtaskId, status, err3.Error())
// 		}
// 	}
// 	return nil
// }

// func DeleteTaskInTaskInfos(taskId string) error {
// 	result := db.GMysalDB.Exec("delete from task_infos where task_id = ?", taskId)
// 	if result.Error != nil {
// 		seelog.Errorf("delete %s task in task infos failed(del):%s", taskId, result.Error)
// 		return result.Error
// 	}
// 	return nil
// }

// 检测任务名是否存在
func CheckTaskNameExist(taskName string) (*TaskInfo, error) {
	var task = TaskInfo{}
	result := db.GMysalDB.Model(&TaskInfo{}).Where("task_name = ?", taskName).Take(&task)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &task, nil
}

// 根据taskid 更新task状态
func UpdateTaskStatusByTaskID(taskid string, status string) error {
	sqlstr := "update task_infos set status = '" + status + "' where task_id = '" + taskid + "'"
	result := db.GMysalDB.Exec(sqlstr)
	if result.Error != nil {
		seelog.Errorf("update task status error:%s", result.Error.Error())
		return result.Error
	}
	return nil
}

// 根据taskid,old status 更新task状态
func UpdateTaskStatusByTaskStatus(taskid, newStatus, oldStatus string) error {
	result := db.GMysalDB.Model(&TaskInfo{}).Where("task_id = ? and status = ?", taskid, oldStatus).Update("status", newStatus)
	if result.Error != nil {
		seelog.Errorf("update task status error:%s", result.Error.Error())
		return result.Error
	}
	return nil
}

// 根据自增id更新taskinfo.
func (u *TaskInfo) UpdateTaskInfoByID(columns ...interface{}) (int64, error) {
	if len(columns) == 0 {
		return 0, nil
	}
	upd := *u
	result := db.GMysalDB.Model(&u).Select(columns[0], columns[1:]...).Updates(upd)
	return result.RowsAffected, result.Error
}

// 根据条件筛选任务列表.
func GetTaskListByFilter(req TaskListFilter) (*[]TaskInfo, int64, error) {
	var taskList = make([]TaskInfo, 0)
	var total int64
	var query = db.GMysalDB.Model(&TaskInfo{})
	if len(req.UserIDs) > 0 {
		switch req.PageTab {
		case utils.TASK_PAGE_TAB_MANAGE:
			query.Where("creator_id IN ?", req.UserIDs)
		case utils.TASK_PAGE_TAB_ALLOCAT:
			alloc := fmt.Sprintf("%d", req.UserIDs[0])
			query.Where("FIND_IN_SET (?, multi_allocator)", alloc)
		case utils.TASK_PAGE_TAB_ANNOTAT:
			alloc := fmt.Sprintf("%d", req.UserIDs[0])
			query.Where("FIND_IN_SET (?, multi_annotator)", alloc)
		case utils.TASK_PAGE_TAB_AUDIT:
			alloc := fmt.Sprintf("%d", req.UserIDs[0])
			query.Where("FIND_IN_SET (?, multi_auditor)", alloc)
		case utils.TASK_PAGE_TAB_ANALYS:
			alloc := fmt.Sprintf("%d", req.UserIDs[0])
			query.Where("FIND_IN_SET (?, multi_analysts)", alloc)
		}
	}
	if len(req.StatusList) > 0 {
		query.Where("status IN ?", req.StatusList)
	} else {
		query.Where("status != ?", utils.TASK_STATUS_DELETED)
	}
	if len(req.TIDs) > 0 {
		query.Not("id IN ?", req.TIDs)
	}
	switch req.SearchType {
	case utils.TASK_SEARCH_TYPE_ID:
		query.Where(fmt.Sprintf("task_id like '%%%s%%'", req.SearchValue))
	case utils.TASK_SEARCH_TYPE_CREATOR:
		query.Where(fmt.Sprintf("task_creator like '%%%s%%'", req.SearchValue))
	case utils.TASK_SEARCH_TYPE_NAME:
		query.Where(fmt.Sprintf("task_name like '%%%s%%'", req.SearchValue))
	case utils.TASK_SEARCH_TYPE_STATUS:
		query.Where("status = ?", req.SearchValue)
	case utils.TASK_SEARCH_TYPE_OPERATOR:
		query.Where(
			"(status in (?) and FIND_IN_SET (?, multi_allocator)) or (status in (?, ?) and FIND_IN_SET (?, multi_annotator)) or (status in (?,?) and FIND_IN_SET (?, multi_auditor)) or (status in (?,?, ?) and FIND_IN_SET (?, multi_analysts))",
			utils.TASK_STATUS_CREATED,
			req.SearchValue,
			utils.TASK_STATUS_ALLOCATED,
			utils.TASK_STATUS_ANNOTATING,
			req.SearchValue,
			utils.TASK_STATUS_ANNOTATED,
			utils.TASK_STATUS_AUDITING,
			req.SearchValue,
			utils.TASK_STATUS_AUDITED,
			utils.TASK_STATUS_ANALYSTING,
			utils.TASK_STATUS_COMPLETED,
			req.SearchValue,
		)
	}

	query.Count(&total)
	result := query.Limit(req.PageSize).Offset((req.PageNum - 1) * req.PageSize).Order("task_createtime desc").Find(&taskList)
	if result.Error != nil {
		return &taskList, total, result.Error
	}
	return &taskList, total, nil
}

func GetTaskInfoByID(taskID string) (*TaskInfo, error) {
	var task = TaskInfo{}
	result := db.GMysalDB.Model(&TaskInfo{}).Where("task_id = ?", taskID).Take(&task)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &task, nil
}

type DoTaskUserNum struct {
	UserID uint32 `json:"user_id"`
	Num    int    `json:"num"`
}
