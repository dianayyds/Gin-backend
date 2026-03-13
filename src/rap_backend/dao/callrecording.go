package dao

import (
	"errors"
	"rap_backend/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CallidRecording struct {
	RecId   int32  `gorm:"primaryKey;AUTO_INCREMENT" json:"rec_id" from:"rec_id"`
	CallId  string `json:"call_id" form:"call_id"`
	Address string `json:"address" form:"address"`
}

//批量插入数据
func CreateCallRecordingInBatches(tc []CallidRecording) error {
	result := db.GMysalDB.Clauses(clause.Insert{Modifier: "IGNORE"}).Model(&CallidRecording{}).CreateInBatches(tc, 100)
	return result.Error
}

//根据 callid 获取录音文件信息
func GetCallRecordingInfoByCallID(callIDs []string) (*[]CallidRecording, error) {
	var cr = make([]CallidRecording, 0, len(callIDs))
	result := db.GMysalDB.Model(&CallidRecording{}).Where("call_id IN ?", callIDs).Find(&cr)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &cr, nil
}
