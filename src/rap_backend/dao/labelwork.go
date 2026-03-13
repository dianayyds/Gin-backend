package dao

import (
	"rap_backend/db"
	"strconv"
	"time"

	"github.com/cihub/seelog"
)

type LabelworkInfo struct {
	Id             int32     `gorm:"primaryKey;AUTO_INCREMENT" json:"id" from:"id"`
	TaskId         string    `json:"task_id" from:"task_id"`
	SubtaskId      string    `json:"subtask_id" from:"subtask_id"`
	CallId         string    `json:"call_id" from:"call_id"`
	LabelId        string    `json:"label_id" from:"label_id"`
	LabelName      string    `json:"label_name" form:"label_name"`
	LabelValue     string    `json:"label_value" form:"label_value"`
	AuditorContent string    `json:"auditor_content" form:"auditor_content"`
	LabelTime      time.Time `json:"label_time" form:"label_time"`
	LabelDesc      string    `json:"label_desc" form:"label_value"`
	IsOptional     int       `json:"is_optional" form:"is_optional"`
	IsEditable     int       `json:"is_editable" form:"is_editable"`
	Status         string    `json:"status" form:"status"`
}

func (t *LabelworkInfo) AddLabelworkInfo() error {
	result := db.GMysalDB.Create(t)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (u *LabelworkInfo) UpdateLabelworkInfo(columns ...interface{}) (int64, error) {
	if len(columns) == 0 {
		return 0, nil
	}
	upd := *u
	result := db.GMysalDB.Model(&u).Select(columns[0], columns[1:]...).Updates(upd)
	return result.RowsAffected, result.Error
}

func GetAllCallidListBySubtaskId(subtaskId string) (map[string]string, error) {
	var labelworkList []LabelworkInfo
	var callidList = make(map[string]string)
	result := db.GMysalDB.Model(&LabelworkInfo{}).Where("subtask_id = ?", subtaskId).Order("call_id").Group("call_id").Find(&labelworkList)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected != 0 {
		for _, v := range labelworkList {
			callidList[v.CallId] = v.Status
		}
		return callidList, nil
	} else {
		return nil, nil
	}
}

func GetAllLabelWorkListBySubtaskId(subtaskId string) (*[]LabelworkInfo, error) {
	var labelworkList []LabelworkInfo
	result := db.GMysalDB.Model(&LabelworkInfo{}).Where("subtask_id = ?", subtaskId).Order("call_id").Find(&labelworkList)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected != 0 {
		return &labelworkList, nil
	} else {
		return nil, nil
	}
}

func GetLabelWorkListBySubtaskId(subtaskId string, pageNum int, pageSize int) (*[]LabelworkInfo, error) {
	var labelworkList []LabelworkInfo
	result := db.GMysalDB.Model(&LabelworkInfo{}).Where("subtask_id = ?", subtaskId).Limit(pageSize).Offset((pageNum - 1) * pageSize).Group("call_id").Order("call_id").Find(&labelworkList)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected != 0 {
		return &labelworkList, nil
	} else {
		return nil, nil
	}
}

func GetLabelWorkListBySubtaskId2(subtaskIdList []string, pageNum int, pageSize int) (*[]LabelworkInfo, error) {
	var labelworkList []LabelworkInfo
	strSubtaskIdList := ""
	for _, subTaskId := range subtaskIdList {
		strSubtaskIdList = strSubtaskIdList + "'" + subTaskId + "',"
	}
	strSubtaskIdList = strSubtaskIdList[0 : len(strSubtaskIdList)-1]

	sqlstr := "select call_id from labelwork_infos where subtask_id in (" + strSubtaskIdList + ") group by call_id  limit " + strconv.Itoa(pageSize) + " offset " + strconv.Itoa((pageNum-1)*pageSize)
	seelog.Infof("sqlstr: %s", sqlstr)
	rows, err := db.GMysalDB.Raw(sqlstr).Rows()
	if err != nil {
		seelog.Errorf("get callid failed: %s", err.Error())
		return nil, err
	}
	defer rows.Close()
	var showCallidList []string
	for rows.Next() {
		var extCallid string
		rows.Scan(&extCallid)
		showCallidList = append(showCallidList, extCallid)
	}
	if len(showCallidList) == 0 {
		seelog.Infof("No callid")
		return nil, err
	}

	result := db.GMysalDB.Model(&LabelworkInfo{}).Debug().Where("call_id IN ? and subtask_id IN ?", showCallidList, subtaskIdList).Find(&labelworkList)
	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected != 0 {
		return &labelworkList, nil
	} else {
		return nil, nil
	}
}

func GetLabelWorkDetail(taskId string, callid string) (*[]LabelworkInfo, error) {
	var labelworkDetail = make([]LabelworkInfo, 0)
	result := db.GMysalDB.Model(&LabelworkInfo{}).Where("task_id = ? and call_id = ?", taskId, callid).Order("label_name").Find(&labelworkDetail)
	if result.Error != nil {
		return nil, result.Error
	}
	return &labelworkDetail, nil
}

// func UpdteLabelWorkDetail(taskId string, callid string, updatedDetail *[]LabelworkInfo) error {
// 	var err2 error
// 	for _, updatedContent := range *updatedDetail {
// 		result := db.GMysalDB.Exec("delete from labelwork_infos where subtask_id = ? and call_id = ? and label_id = ?", taskId, callid, updatedContent.LabelId)
// 		if result.Error != nil {
// 			seelog.Errorf("updated label work failed(del):%s", result.Error)
// 			continue
// 		} else {
// 			err2 = updatedContent.AddLabelworkInfo()
// 			if err2 != nil {
// 				seelog.Errorf("updated label work failed(ins):%s", updatedContent)
// 				continue
// 			}
// 		}
// 	}
// 	if err2 != nil {
// 		return err2
// 	} else {
// 		return nil
// 	}

// }

// func DeleteSubtaskLabelWork(subtaskId string) error {
// 	result := db.GMysalDB.Exec("delete from labelwork_infos where subtask_id = ?", subtaskId)
// 	if result.Error != nil {
// 		seelog.Errorf("delete %s label work failed(del):%s", subtaskId, result.Error)
// 		return result.Error
// 	}
// 	return nil
// }

func CreateTaskLabelWorkInBatches(tc []LabelworkInfo) error {
	result := db.GMysalDB.Model(&LabelworkInfo{}).CreateInBatches(tc, 500)
	return result.Error
}

func GetLabelWorkDetailByCallIDs(taskId string, callids []string) (map[string]map[string]LabelworkInfo, error) {
	var labelworkDetail = make([]LabelworkInfo, 0)
	var labelMap = make(map[string]map[string]LabelworkInfo, 0)
	result := db.GMysalDB.Model(&LabelworkInfo{}).Where("task_id = ? and call_id IN ?", taskId, callids).Order("id").Find(&labelworkDetail)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, lab := range labelworkDetail {
		if _, ok := labelMap[lab.CallId]; !ok {
			labelMap[lab.CallId] = make(map[string]LabelworkInfo)
		}
		labelMap[lab.CallId][lab.LabelName] = lab
	}
	return labelMap, nil
}

func UpdateLabelworkInfoByCallIDLabID(taskID, callID string, labids []string, upd map[string]interface{}) (int64, error) {
	result := db.GMysalDB.Model(&LabelworkInfo{}).Where("task_id = ? and call_id = ? and label_id IN ?", taskID, callID, labids).Updates(upd)
	return result.RowsAffected, result.Error
}
