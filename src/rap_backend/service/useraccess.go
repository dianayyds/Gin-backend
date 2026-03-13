package service

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"rap_backend/config"
	"rap_backend/dao"
	"rap_backend/internal/jwtauth"
	"rap_backend/utils"
	"time"

	"github.com/cihub/seelog"
	"github.com/dgrijalva/jwt-go"
)

type UserLoginDTO struct {
	LoginName string `json:"login_name"`
	Password  string `json:"password"`
}

type UserInfoDTO struct {
	UserID uint32 `json:"user_id"`
}

type UserInfoRetDTO struct {
	Info UserInfoDetail `json:"info"`
}

type UserListDTO struct {
	LoginName string `json:"login_name"`
	UserType  string `json:"user_type"`
	PageNum   int    `json:"page_num"`
	PageSize  int    `json:"page_size"`
	Status    int    `json:"status"`
}

type UserListRetDTO struct {
	List  []UserInfoDetail `json:"list"`
	Total int64            `json:"total"`
}

type UserInfoDetail struct {
	UserID      uint32      `json:"user_id"`
	LoginName   string      `json:"login_name"`
	UserName    string      `json:"user_name"`
	Email       string      `json:"email"`
	CountryId   int         `json:"country_id"`
	PageTab     []int32     `json:"page_tab"`
	Preview     []int32     `json:"preview"`
	Roles       []int       `json:"roles"`
	RoleInfo    []RoleInfo  `json:"role_info"`
	DepartId    int         `json:"depart_id"`
	Status      int         `json:"status"`
	CountryInfo CountryInfo `json:"country_info"`
	DataAll     bool        `json:"-"`
}

type UserCreateDTO struct {
	UserName  string  `json:"user_name"`
	LoginName string  `json:"login_name"`
	Password  string  `json:"password"`
	CountryId int     `json:"country_id"`
	Roles     []int32 `json:"roles"`
}

type ResponseRetDTO struct {
	Ret   int   `json:"ret"`
	Count int64 `json:"count"`
}

type UserEditDTO struct {
	UserID    uint32 `json:"user_id"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
	CountryId int    `json:"country_id"`
	Roles     []int  `json:"roles"`
}

type UserStatusDTO struct {
	UserID uint32 `json:"user_id"`
	Status int    `json:"status"`
}

func CreateUser(req UserCreateDTO) (uint32, error) {
	u1, err := dao.SelUserInfoByName(req.LoginName)
	if err != nil {
		return 0, err
	}
	if u1 != nil { //loginname已存在
		return 0, fmt.Errorf("账号已存在")
	}
	var user = dao.User{
		UserName:  req.UserName,
		LoginName: req.LoginName,
		Password:  UserPasswordInit(req.Password),
		CountryId: req.CountryId,
		Roles:     utils.Int32SliceToString(req.Roles),
		Status:    config.USER_STATUS_NORMAL,
	}
	err = user.AddUserInfo()
	if err != nil {
		return 0, err
	}
	return user.UserId, err
}

func UserLogin(loginName, password string) (string, error) {
	user, err := dao.SelUserInfoByName(loginName)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user does not exist")
	}
	if user.Password != UserPasswordInit(password) {
		return "", errors.New("loginname or pwd not match")
	}
	if user.Status != config.USER_STATUS_NORMAL {
		return "", errors.New("user is forbid")
	}
	return SetUserToken(*user)
}

func UserInfo(uid uint32) (*UserInfoDetail, error) {
	_, users, err := dao.SelUserInfoByIDs([]uint32{uid})
	if err != nil {
		return nil, err
	}
	u, ok := users[uid]
	if !ok {
		return nil, nil
	}
	info := UserInfoDetail{
		UserID:      u.UserId,
		LoginName:   u.LoginName,
		UserName:    u.UserName,
		Email:       u.Email,
		CountryId:   u.CountryId,
		PageTab:     []int32{},
		Preview:     []int32{},
		Roles:       utils.StringToIntSlice(u.Roles),
		RoleInfo:    []RoleInfo{},
		DepartId:    u.DepartId,
		Status:      u.Status,
		CountryInfo: CountryInfo{},
	}
	if ci, ok := CountryCache[uint32(u.CountryId)]; ok {
		info.CountryInfo = ci
	}
	roleIDs := []int32{}
	if u.Roles != "" {
		ids := utils.StringToInt32Slice(u.Roles)
		roleIDs = append(roleIDs, ids...)
		roleInfo, _ := RoleInfoByIDs(roleIDs)

		for _, role := range roleInfo {
			if role.RoleData == 2 {
				info.DataAll = true
				break
			}
		}
		info.RoleInfo = roleInfo
	}
	info.Preview, info.PageTab, _ = PurviewListByRoleIDs(roleIDs)
	return &info, nil
}

func UserList(loginName, userType string, userStatus int, offset, size int) ([]UserInfoDetail, int64, error) {
	seelog.Infof("rapUserList_1 t:%d", time.Now().Unix())

	infos := make([]UserInfoDetail, 0)
	roleIDs := []int32{}
	orderBy := "user_id asc"
	switch userType {
	case utils.USER_TYPE_ALLOCATOR:
		roleIDs = AllocatRoleIDs
	case utils.USER_TYPE_ANNOTATOR:
		orderBy = "user_name asc"
		roleIDs = AnnotatRoleIDs
	case utils.USER_TYPE_AUDITOR:
		orderBy = "user_name asc"
		roleIDs = AuditRoleIDs
	case utils.USER_TYPE_ANALYSTS:
		orderBy = "user_name asc"
		roleIDs = AnalystRoleIDs
	}
	list, total, err := dao.SelUserList(loginName, roleIDs, userStatus, offset, size, orderBy)
	seelog.Infof("rapUserList_4 t:%d", time.Now().Unix())

	if err != nil {
		return infos, total, err
	}
	for _, u := range *list {
		inf := UserInfoDetail{
			UserID:    u.UserId,
			LoginName: u.LoginName,
			UserName:  u.UserName,
			Email:     u.Email,
			CountryId: u.CountryId,
			Preview:   []int32{},
			Roles:     utils.StringToIntSlice(u.Roles),
			RoleInfo:  []RoleInfo{},
			DepartId:  u.DepartId,
			Status:    u.Status,
		}
		if u.Roles != "" {
			ids := utils.StringToInt32Slice(u.Roles)
			roleInfo, _ := RoleInfoByIDs(ids)
			inf.RoleInfo = roleInfo
		}
		infos = append(infos, inf)
	}
	seelog.Infof("rapUserList_5 t:%d", time.Now().Unix())

	return infos, total, nil
}

func EditUser(req UserEditDTO) error {
	u := dao.User{
		UserId:    req.UserID,
		UserName:  req.UserName,
		CountryId: req.CountryId,
		Roles:     utils.IntSliceToString(req.Roles),
	}
	columns := []string{"user_name", "country_id", "roles"}
	if req.Password != "" {
		columns = append(columns, "password")
		u.Password = UserPasswordInit(req.Password)
	}
	_, err := u.UpdateUserInfo(columns)
	return err
}

func EditUserStatus(req UserStatusDTO) error {
	u := dao.User{
		UserId: req.UserID,
		Status: req.Status,
	}
	columns := []string{"status"}
	_, err := u.UpdateUserInfo(columns)
	return err
}

func SetUserToken(user dao.User) (string, error) {
	before := time.Now().Unix() - 1000
	expiresAt := time.Now().Unix() + 24*3600
	j := jwtauth.NewJWT()
	claims := jwtauth.CustomClaims{
		UserId:   int64(user.UserId),
		Username: user.LoginName,
		StandardClaims: jwt.StandardClaims{
			NotBefore: before,
			ExpiresAt: expiresAt,
			Issuer:    "rap_management",
		},
	}
	token, err := j.CreateToken(claims)
	return token, err
}

func UserPasswordInit(password string) string {
	var pwd = "Sfa8biXQfyLVaPfz"
	if password != "" {
		pwd = password
	}
	return sha1password(pwd)
}

func sha1password(str string) string {
	o := sha1.New()
	o.Write([]byte(str))
	return hex.EncodeToString(o.Sum(nil))
}

func UserShortInfoByIDs(uid []uint32) (map[uint32]UserInfoDetail, error) {
	_, users, err := dao.SelUserInfoByIDs(uid)
	userMap := make(map[uint32]UserInfoDetail, 0)
	if err != nil {
		return userMap, err
	}
	for _, u := range users {
		info := UserInfoDetail{
			UserID:    u.UserId,
			LoginName: u.LoginName,
			UserName:  u.UserName,
			Email:     u.Email,
			CountryId: u.CountryId,
			PageTab:   []int32{},
			Preview:   []int32{},
			RoleInfo:  []RoleInfo{},
			DepartId:  u.DepartId,
			Status:    u.Status,
		}
		userMap[u.UserId] = info
	}
	return userMap, nil
}

// 根据uid获取所属国家信息
func GetCountryInfoByUID(userID uint32) CountryInfo {
	var country = CountryInfo{}
	user, err := UserShortInfoByIDs([]uint32{userID})
	if err != nil {
		return country
	}
	u, ok := user[userID]
	if !ok {
		return country
	}
	if c, ok := CountryCache[uint32(u.CountryId)]; ok {
		country = c
	}
	return country
}

func GetUserInfoByLoginName(loginName string) (UserInfoDetail, error) {
	info := UserInfoDetail{}
	u, err := dao.SelUserInfoByName(loginName)
	if err != nil || u == nil {
		return info, err
	}
	info = UserInfoDetail{
		UserID:    u.UserId,
		LoginName: u.LoginName,
		UserName:  u.UserName,
		Email:     u.Email,
		CountryId: u.CountryId,
		PageTab:   []int32{},
		Preview:   []int32{},
		RoleInfo:  []RoleInfo{},
		DepartId:  u.DepartId,
		Status:    u.Status,
	}
	return info, err
}

func GetUserInfoByUserName(userName string) (UserInfoDetail, error) {
	info := UserInfoDetail{}
	u, err := dao.SelUserInfoByUserName(userName)
	if err != nil || u == nil {
		return info, err
	}
	info = UserInfoDetail{
		UserID:    u.UserId,
		LoginName: u.LoginName,
		UserName:  u.UserName,
		Email:     u.Email,
		CountryId: u.CountryId,
		PageTab:   []int32{},
		Preview:   []int32{},
		RoleInfo:  []RoleInfo{},
		DepartId:  u.DepartId,
		Status:    u.Status,
	}
	return info, err
}
