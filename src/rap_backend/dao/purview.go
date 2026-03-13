package dao

import "rap_backend/db"

type Purview struct {
	Id          int32  `gorm:"AUTO_INCREMENT" json:"id" from:"id"`
	PurviewName string `json:"purview_name" from:"purview_name"`
	ShowName    string `json:"show_name" from:"show_name"`
	ApiName     string `json:"api_name" from:"api_name"`
	UpperLevel  int32  `json:"upper_level" from:"upper_level"`
}

func SelPurviewInfoByIDs(ids []int32) (*[]Purview, error) {
	var purs = make([]Purview, 0, len(ids))
	result := db.GMysalDB.Model(&Purview{}).Where("id IN ?", ids).Find(&purs)
	return &purs, result.Error
}
