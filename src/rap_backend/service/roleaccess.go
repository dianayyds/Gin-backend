package service

import (
	"rap_backend/config"
	"rap_backend/dao"
	"rap_backend/utils"
)

type RoleListDTO struct {
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}

type RoleListRetDTO struct {
	List  []RoleInfo `json:"list"`
	Total int64      `json:"total"`
}

type RoleInfo struct {
	RoleId   int32  `json:"role_id" from:"role_id"`
	RoleName string `json:"role_name" from:"role_name"`
	RoleData int    `json:"role_data" from:"role_data"`
	RoleFunc string `json:"role_func" from:"role_func"`
	Status   int    `json:"status" from:"status"`
}

func RoleInfoByIDs(ids []int32) ([]RoleInfo, error) {
	infos := make([]RoleInfo, 0)
	if len(ids) == 0 {
		return infos, nil
	}
	list, _, err := dao.SelRoleInfoByIDs(ids)
	if err != nil {
		return infos, err
	}
	for _, li := range *list {
		inf := RoleInfo{
			RoleId:   li.RoleId,
			RoleName: li.RoleName,
			RoleData: li.RoleData,
			RoleFunc: li.RoleFunc,
			Status:   li.Status,
		}
		infos = append(infos, inf)
	}
	return infos, err
}

func PurviewListByRoleIDs(ids []int32) ([]int32, []int32, error) {
	purs := []int32{}
	parents := []int32{}
	if len(ids) == 0 {
		return purs, parents, nil
	}
	roles, _, err := dao.SelRoleInfoByIDs(ids)
	if err != nil {
		return purs, parents, nil
	}
	purids := []int32{}
	for _, role := range *roles {
		if role.Status != config.USER_STATUS_NORMAL || role.RoleFunc == "" {
			continue
		}
		ids := utils.StringToInt32Slice(role.RoleFunc)
		purids = append(purids, ids...)
	}
	list, err := dao.SelPurviewInfoByIDs(purids)
	if err != nil {
		return purs, parents, nil
	}
	parentMap := make(map[int32]bool, 0)
	for _, li := range *list {
		if li.UpperLevel > 0 {
			if _, ok := parentMap[li.UpperLevel]; !ok {
				parents = append(parents, li.UpperLevel)
				parentMap[li.UpperLevel] = true
			}
		}
		purs = append(purs, li.Id)
	}
	return purs, parents, nil
}

func RoleList(offset, size int) ([]RoleInfo, int64, error) {
	infos := make([]RoleInfo, 0)
	list, total, err := dao.GetRoleList(offset, size)
	if err != nil {
		return infos, total, err
	}
	for _, li := range *list {
		inf := RoleInfo{
			RoleId:   li.RoleId,
			RoleName: li.RoleName,
			RoleData: li.RoleData,
			RoleFunc: li.RoleFunc,
			Status:   li.Status,
		}
		infos = append(infos, inf)
	}
	return infos, total, nil
}
