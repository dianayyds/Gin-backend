package dao

import "rap_backend/db"

type Country struct {
	CountryId uint32 `gorm:"primaryKey;autoIncrement" json:"country_id" from:"country_id"`
	Code      string `json:"code" form:"code"`
	Name      string `json:"name" form:"name"`
	TimeZone  string `json:"time_zone" form:"time_zone"`
	ZoneId    string `json:"zone_id" form:"zone_id"`
	Status    int    `json:"status" form:"status"`
}

func (c *Country) TableName() string {
	return "countrys"
}

func GetCountryList() (*[]Country, error) {
	var countrys = make([]Country, 0)
	result := db.GMysalDB.Model(&Country{}).Where("status=1").Find(&countrys)
	if result.Error != nil {
		return &countrys, result.Error
	}
	return &countrys, nil
}

func ExecSQL(s string) error {
	result := db.GMysalDB.Exec(s)
	return result.Error
}
