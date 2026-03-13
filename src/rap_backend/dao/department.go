package dao

type Department struct {
	DepartId   uint32 `gorm:"primaryKey;autoIncrement" json:"depart_id" from:"depart_id"`
	DepartName string `json:"depart_name" form:"depart_name"`
	UpperLevel string `json:"upper_level" form:"upper_level"`
}
