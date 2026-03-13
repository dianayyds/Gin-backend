package service

import (
	"fmt"
	"rap_backend/config"
	"rap_backend/dao"
	"strings"
	"time"

	"github.com/cihub/seelog"
	"github.com/golibs/uuid"
)

type GetLabelDTO struct {
	PageTab   string `json:"page_tab" form:"page_tab"`
	LabelName string `json:"label_name" form:"label_name"`
	PageNum   int    `json:"page_num" form:"page_num"`
	PageSize  int    `json:"page_size" form:"page_size"`
}

type GetLabelRetDTO struct {
	LabelCnt      int64           `json:"label_cnt" form:"label_cnt"`
	LabelInfoList []LabelInfoItem `json:"label_info_list" form:"label_info_list"`
}

type LabelInfoItem struct {
	LabelId         string   `json:"label_id" form:"label_id"`
	LabelName       string   `json:"label_name" form:"label_name"`
	LabelDesc       string   `json:"label_desc" form:"label_desc"`
	LabelIsOptional int      `json:"label_is_optional" form:"label_is_optional"`
	LabelIsEditable int      `json:"label_is_editable" form:"label_is_editable"`
	LabelCreator    string   `json:"label_creator" form:"label_creator"`
	LabelType       int      `json:"label_type" form:"label_type"` //0文本 1单选 2多选
	LabelOptions    []string `json:"label_options" form:"label_options"`
	IsColor         int      `json:"is_color" form:"is_color"`
	Belong          int      `json:"-" `
}

type CreateLabelDTO struct {
	LabelName       string    `json:"label_name" form:"label_name"`
	LabelDesc       string    `json:"label_desc" form:"label_desc"`
	LabelCreator    string    `json:"label_creator" form:"label_creator"`
	LabelIsOptional int       `json:"label_is_optional" form:"label_is_optional"` //0必填 1选填
	LabelCreatetime time.Time `json:"label_createtime" form:"label_createtime"`
	IsEditable      int       `json:"is_editable" form:"is_editable"` //0不可编辑，1可编辑
	LabelType       int       `json:"label_type" form:"label_type"`   //0文本 1单选 2多选
	LabelOptions    []string  `json:"label_options" form:"label_options"`
	IsColor         int       `json:"is_color" form:"is_color"` //0不可标颜色 1可标颜色
}

type EditLabelDTO struct {
	LabelId         string   `json:"label_id" form:"label_id"`
	LabelName       string   `json:"label_name" form:"label_name"`
	LabelDesc       string   `json:"label_desc" form:"label_desc"`
	LabelIsOptional int      `json:"label_is_optional" form:"label_is_optional"` //0必填 1选填
	IsEditable      int      `json:"is_editable" form:"is_editable"`             //0不可编辑，1可编辑
	LabelType       int      `json:"label_type" form:"label_type"`               //0文本 1单选 2多选
	LabelOptions    []string `json:"label_options" form:"label_options"`
	IsColor         int      `json:"is_color" form:"is_color"` //0不可标颜色 1可标颜色
}
type DelLabelDTO struct {
	LabelId string `json:"label_id" form:"label_id"`
}

func generateNewLabelId() string {
	return "lab_" + strings.Replace(uuid.Rand().Hex(), "-", "", -1)
}

func CreateNewLabel(labelInfo CreateLabelDTO) (string, error) {
	var labelDao dao.LabelInfo
	labelDao.LabelId = generateNewLabelId()
	seelog.Infof("new create label id:%s", labelDao.LabelId)
	labelDao.LabelName = labelInfo.LabelName
	labelDao.LabelCreator = labelInfo.LabelCreator
	labelDao.IsOptional = labelInfo.LabelIsOptional
	labelDao.LabelDesc = labelInfo.LabelDesc
	labelDao.LabelCreatetime = time.Now()
	labelDao.IsEditable = 1
	labelDao.LabelType = labelInfo.LabelType
	labelDao.LabelOptions = strings.Join(labelInfo.LabelOptions, ",")
	labelDao.Status = config.USER_STATUS_NORMAL
	labelDao.IsColor = labelInfo.IsColor
	err := labelDao.AddLabelInfo()
	return labelDao.LabelId, err
}

//编辑label
func EditLabel(req EditLabelDTO) error {
	info, err := dao.GetLabelInfoByLabelName(req.LabelName)
	if err != nil {
		return err
	}
	if info != nil && info.LabelId != "" && info.LabelId != req.LabelId {
		return fmt.Errorf("the label is already existed")
	}
	u := dao.LabelInfo{
		LabelId:      req.LabelId,
		LabelName:    req.LabelName,
		LabelDesc:    req.LabelDesc,
		IsOptional:   req.LabelIsOptional,
		LabelType:    req.LabelType,
		LabelOptions: strings.Join(req.LabelOptions, ","),
		IsColor:      req.IsColor,
		Belong:       info.Belong,
	}
	columns := []string{"label_name", "label_desc", "label_type", "label_options", "is_optional", "is_color"}
	_, err = u.UpdateLabelInfo(columns)
	return err
}

func DelLabel(req DelLabelDTO) error {
	u := dao.LabelInfo{
		LabelId: req.LabelId,
		Status:  config.USER_STATUS_FORBID,
	}
	columns := []string{"status"}
	_, err := u.UpdateLabelInfo(columns)
	return err
}

//label list.
func GetLabelList(labelName string, pageNum, pageSize, status int, orderBy string) ([]LabelInfoItem, int64, error) {
	seelog.Infof("GetLabelList, pagenum:%d, pagesize:%d", pageNum, pageSize)
	labelInfoList := make([]LabelInfoItem, 0)
	offset := (pageNum - 1) * pageSize
	labelList_0, cnt, err := dao.GetAllLabelInfoList(labelName, offset, pageSize, status, orderBy)
	if err != nil {
		seelog.Errorf("get un_editable label list failed error :%s", err.Error())
		return labelInfoList, cnt, err
	}

	if len(*labelList_0) == 0 {
		return labelInfoList, cnt, nil
	}
	seelog.Infof("get label list length :%d", len(*labelList_0))
	for _, info := range *labelList_0 {
		v := LabelInfoItem{
			LabelId:         info.LabelId,
			LabelName:       info.LabelName,
			LabelDesc:       info.LabelDesc,
			LabelCreator:    info.LabelCreator,
			LabelIsOptional: info.IsOptional,
			LabelIsEditable: info.IsEditable,
			LabelType:       info.LabelType,
			LabelOptions:    []string{},
			IsColor:         info.IsColor,
			Belong:          info.Belong,
		}
		if info.LabelOptions != "" {
			v.LabelOptions = strings.Split(info.LabelOptions, ",")
		}
		labelInfoList = append(labelInfoList, v)
	}
	return labelInfoList, cnt, nil
}

//根据labelname获取label信息
func GetLabelInfoByName(labelName string) (LabelInfoItem, error) {
	label := LabelInfoItem{}
	info, err := dao.GetLabelInfoByLabelName(labelName)
	if err != nil {
		seelog.Errorf("GetLabelInfoByName failed error :%s", err.Error())
		return label, err
	}

	label = LabelInfoItem{
		LabelId:         info.LabelId,
		LabelName:       info.LabelName,
		LabelDesc:       info.LabelDesc,
		LabelCreator:    info.LabelCreator,
		LabelIsOptional: info.IsOptional,
		LabelIsEditable: info.IsEditable,
		LabelType:       info.LabelType,
		LabelOptions:    []string{},
		IsColor:         info.IsColor,
		Belong:          info.Belong,
	}
	if info.LabelOptions != "" {
		label.LabelOptions = strings.Split(info.LabelOptions, ",")
	}
	return label, nil
}
