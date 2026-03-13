package dao

import "rap_backend/db"

type Role struct {
	RoleId   int32  `gorm:"AUTO_INCREMENT" json:"role_id" from:"role_id"`
	RoleName string `json:"role_name" from:"role_name"`
	RoleData int    `json:"role_data" from:"role_data"`
	RoleFunc string `json:"role_func" from:"role_func"`
	Status   int    `json:"status" from:"status"`
}

func SelRoleInfoByIDs(ids []int32) (*[]Role, map[int32]Role, error) {
	var roles = make([]Role, 0, len(ids))
	var roleMap = make(map[int32]Role, 0)
	if len(ids) == 0 {
		return &roles, roleMap, nil
	}
	result := db.GMysalDB.Model(&Role{}).Where("role_id IN ?", ids).Find(&roles)
	for _, user := range roles {
		roleMap[user.RoleId] = user
	}
	return &roles, roleMap, result.Error
}

func GetRoleList(offset, size int) (*[]Role, int64, error) {
	var roles = make([]Role, 0)
	var total int64
	query := db.GMysalDB.Model(&Role{}).Where("status=1")
	query.Count(&total)
	result := query.Offset(offset).Limit(size).Find(&roles)
	if result.Error != nil {
		return &roles, total, result.Error
	}
	return &roles, total, nil
}
