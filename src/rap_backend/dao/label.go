package dao

import (
	"errors"
	"rap_backend/db"
	"strings"
	"time"

	"github.com/cihub/seelog"
)

type LabelInfo struct {
	LabelId         string    `gorm:"primarykey" json:"label_id" from:"label_id"`
	LabelName       string    `json:"label_name" form:"label_name"`
	LabelCreator    string    `json:"label_creator" form:"label_creator"`
	LabelCreatetime time.Time `json:"label_createtime" form:"label_createtime"`
	IsOptional      int       `json:"is_optional" form:"is_optional"`
	IsEditable      int       `json:"is_editable" form:"is_editable"`
	LabelDesc       string    `json:"label_desc" form:"label_desc"`
	LabelType       int       `json:"label_type" form:"label_type"`
	LabelOptions    string    `json:"label_options" form:"label_options"`
	IsColor         int       `json:"is_color" form:"is_color"`
	Belong          int       `json:"belong" form:"belong"`
	Status          int       `json:"status" form:"status"`
}

func (t *LabelInfo) AddLabelInfo() error {
	seelog.Infof("AddLabelInfo:%s", t)
	var checklabelList []LabelInfo
	result := db.GMysalDB.Model(&LabelInfo{}).Where("label_name=?", strings.TrimSpace(t.LabelName)).Find(&checklabelList)

	if result.Error != nil {
		return result.Error
	} else if result.RowsAffected != 0 {
		return errors.New("the label is already existed")
	}

	result2 := db.GMysalDB.Create(t)
	if result2.Error != nil {
		seelog.Errorf("addlabelinfo err:%s", result2.Error.Error())
		return result2.Error
	}
	return nil
}

func (t *LabelInfo) UpdateLabelInfo(columns ...interface{}) (int64, error) {
	if len(columns) == 0 {
		return 0, nil
	}
	upd := *t
	result := db.GMysalDB.Model(&t).Select(columns[0], columns[1:]...).Updates(upd)
	return result.RowsAffected, result.Error
}

func GetLabelListCount() (int64, error) {
	var cnt int64
	result := db.GMysalDB.Model(&LabelInfo{}).Count(&cnt)
	if result.Error != nil {
		return 0, result.Error
	}
	return cnt, nil
}

func GetAllLabelInfoList(labelName string, offset, pageSize, status int, orderBy string) (*[]LabelInfo, int64, error) {
	var labelList = make([]LabelInfo, 0)
	var total int64
	query := db.GMysalDB.Model(&LabelInfo{})
	if labelName != "" {
		name := "%" + labelName + "%"
		query.Where("label_name like ?", name)
	}
	if status > 0 {
		query.Where("status = ?", status)
	}
	query.Count(&total)

	result := query.Limit(pageSize).Offset(offset).Order(orderBy).Find(&labelList)
	if result.Error != nil {
		return &labelList, total, result.Error
	}
	return &labelList, total, nil
}

func GetLabelInfoByLabelId(labelId string) (*LabelInfo, error) {
	var labInfo LabelInfo
	result := db.GMysalDB.Model(&LabelInfo{}).Where("label_id = ?", labelId).Find(&labInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &labInfo, nil
}

func GetLabelInfoByLabelName(name string) (*LabelInfo, error) {
	var labInfo LabelInfo
	result := db.GMysalDB.Model(&LabelInfo{}).Where("label_name = ?", name).Find(&labInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &labInfo, nil
}

func GetLabelInfoByIDs(labelIds []string) (map[string]LabelInfo, error) {
	var labInfo = make([]LabelInfo, 0)
	var labMap = make(map[string]LabelInfo, 0)
	result := db.GMysalDB.Model(&LabelInfo{}).Where("label_id IN ?", labelIds).Find(&labInfo)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, lab := range labInfo {
		labMap[lab.LabelId] = lab
	}
	return labMap, nil
}

func GetAllLabels() (*[]LabelInfo, map[string]LabelInfo, error) {
	var labelList = make([]LabelInfo, 0)
	var labelMap = make(map[string]LabelInfo, 0)
	result := db.GMysalDB.Model(&LabelInfo{}).Find(&labelList)
	if result.Error != nil {
		return &labelList, labelMap, result.Error
	}
	for _, ll := range labelList {
		labelMap[ll.LabelId] = ll
	}
	return &labelList, labelMap, nil
}
