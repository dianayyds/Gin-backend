package dao

import (
	"errors"
	"rap_backend/db"
	"time"

	"github.com/cihub/seelog"
	"gorm.io/gorm"
)

type User struct {
	UserId    uint32 `gorm:"primaryKey;autoIncrement" json:"user_id" from:"user_id"`
	LoginName string `json:"login_name" form:"login_name"`
	Password  string `json:"password" form:"password"`
	UserName  string `json:"user_name" form:"user_name"`
	Email     string `json:"email" form:"email"`
	CountryId int    `json:"country_id" form:"country_id"`
	Roles     string `json:"roles" form:"roles"`
	DepartId  int    `json:"depart_id" form:"depart_id"`
	Status    int    `json:"status" form:"status"`
}

func (u *User) AddUserInfo() error {
	result := db.GMysalDB.Create(u)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *User) CreateOrSelectUserInfo() (uint32, error) {
	result := db.GMysalDB.Where(User{LoginName: u.LoginName}).FirstOrCreate(&u)
	return u.UserId, result.Error
}

func (u *User) UpdateUserInfo(columns ...interface{}) (int64, error) {
	if len(columns) == 0 {
		return 0, nil
	}
	upd := *u
	result := db.GMysalDB.Model(&u).Select(columns[0], columns[1:]...).Updates(upd)
	return result.RowsAffected, result.Error
}

func SelUserInfoByIDs(uids []uint32) (*[]User, map[uint32]User, error) {
	var users = make([]User, 0, len(uids))
	var userMap = make(map[uint32]User, 0)
	if len(uids) == 0 {
		return &users, userMap, nil
	}
	result := db.GMysalDB.Model(&User{}).Where("user_id IN ?", uids).Find(&users)
	for _, user := range users {
		userMap[user.UserId] = user
	}
	return &users, userMap, result.Error
}

func SelUserInfoByName(loginName string) (*User, error) {
	var user = User{}
	result := db.GMysalDB.Model(&User{}).Where("login_name = ?", loginName).Take(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

func SelUserInfoByUserName(userName string) (*User, error) {
	var user = User{}
	result := db.GMysalDB.Model(&User{}).Where("user_name = ?", userName).Take(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}

func SelUserList(name string, roleIDs []int32, status int, offset, size int, orderBy string) (*[]User, int64, error) {
	seelog.Infof("rapUserList_2 t:%d", time.Now().Unix())

	var users = make([]User, 0)
	var total int64
	queryset := db.GMysalDB.Model(&User{})
	if name != "" {
		lk := "%" + name + "%"
		queryset.Where("login_name like ? or user_name like ?", lk, lk)
	}
	if status != 0 {
		queryset.Where("status = ?", status)
	}
	if len(roleIDs) > 0 {
		q1 := db.GMysalDB.Model(&User{})
		for k, roleID := range roleIDs {
			if k == 0 {
				q1.Where("find_in_set(?, roles)", roleID)
			} else {
				q1.Or("find_in_set(?, roles)", roleID)
			}
		}
		queryset.Where(q1)
	}
	queryset.Count(&total)
	result := queryset.Offset(offset).Limit(size).Order(orderBy).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	seelog.Infof("rapUserList_3 t:%d", time.Now().Unix())

	return &users, total, nil
}
