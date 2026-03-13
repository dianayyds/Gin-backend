package service

import (
	"rap_backend/utils"
	"time"

	"github.com/cihub/seelog"
)

var (
	CountryCache       map[uint32]CountryInfo
	LabelInfoCache     map[string]LabelInfoItem
	LabelNameInfoCache map[string]LabelInfoItem
	cachedateformat    = "2006-01-02 15:04:05"
	allocatPurview     = []int32{9, 10, 11, 12, 13}
	annotatPurview     = []int32{15}
	auditPurview       = []int32{17}
	analystPurview     = []int32{19}
	AllocatRoleIDs     []int32
	AnnotatRoleIDs     []int32
	AuditRoleIDs       []int32
	AnalystRoleIDs     []int32
)

//国家信息-本地缓存
func RefreshCountrysLocalCache() {
	seelog.Infof("CountrysLocalCache refresh start date:%s", time.Now().Format(cachedateformat))

	list, err := CountryList()
	if err != nil {
		seelog.Errorf("CountrysLocalCache refresh err:%s", err.Error())
		return
	}
	country := make(map[uint32]CountryInfo, 0)
	for _, li := range list {
		country[li.CountryId] = li
	}
	CountryCache = country
	seelog.Infof("CountrysLocalCache refresh success date:%s", time.Now().Format(cachedateformat))
}

//字段列表-本地缓存
func RefreshLabelInfosLocalCache() {
	seelog.Infof("RefreshLabelInfosLocalCache refresh start date:%s", time.Now().Format(cachedateformat))
	list, _, err := GetLabelList("", 1, 1000, 0, "is_editable,label_name")
	if err != nil {
		seelog.Errorf("RefreshLabelInfosLocalCache refresh err:%s", err.Error())
		return
	}
	labels := make(map[string]LabelInfoItem, 0)
	names := make(map[string]LabelInfoItem, 0)
	for _, li := range list {
		labels[li.LabelId] = li
		names[li.LabelName] = li
	}
	LabelInfoCache = labels
	LabelNameInfoCache = names
	seelog.Infof("RefreshLabelInfosLocalCache refresh success date:%s", time.Now().Format(cachedateformat))
}

//拥有分配、标注、审核、分析权限的角色
func RefreshRolePurviewCache() {
	seelog.Infof("RefreshRolePurviewCache refresh start date:%s", time.Now().Format(cachedateformat))
	list, _, err := RoleList(0, 1000)
	if err != nil {
		seelog.Errorf("RefreshRolePurviewCache refresh err:%s", err.Error())
		return
	}
	allocatMap := make(map[int32]bool, 0)
	annotatMap := make(map[int32]bool, 0)
	auditMap := make(map[int32]bool, 0)
	analystMap := make(map[int32]bool, 0)

	AllocatRoleIDs_cache := make([]int32, 0)
	AnnotatRoleIDs_cache := make([]int32, 0)
	AuditRoleIDs_cache := make([]int32, 0)
	AnalystRoleIDs_cache := make([]int32, 0)

	for _, li := range list {
		pids := utils.StringToInt32Slice(li.RoleFunc)
		for _, pid := range pids {
			for _, rid := range allocatPurview {
				if rid == pid {
					if _, ok := allocatMap[li.RoleId]; !ok {
						AllocatRoleIDs_cache = append(AllocatRoleIDs_cache, li.RoleId)
						allocatMap[li.RoleId] = true
						break
					}
				}
			}
			for _, rid := range annotatPurview {
				if rid == pid {
					if _, ok := annotatMap[li.RoleId]; !ok {
						AnnotatRoleIDs_cache = append(AnnotatRoleIDs_cache, li.RoleId)
						annotatMap[li.RoleId] = true
						break
					}
				}
			}
			for _, rid := range auditPurview {
				if rid == pid {
					if _, ok := auditMap[li.RoleId]; !ok {
						AuditRoleIDs_cache = append(AuditRoleIDs_cache, li.RoleId)
						auditMap[li.RoleId] = true
						break
					}
				}
			}
			for _, rid := range analystPurview {
				if rid == pid {
					if _, ok := analystMap[li.RoleId]; !ok {
						AnalystRoleIDs_cache = append(AnalystRoleIDs_cache, li.RoleId)
						analystMap[li.RoleId] = true
						break
					}
				}
			}
		}
	}
	AllocatRoleIDs = AllocatRoleIDs_cache
	AnnotatRoleIDs = AnnotatRoleIDs_cache
	AuditRoleIDs = AuditRoleIDs_cache
	AnalystRoleIDs = AnalystRoleIDs_cache
	seelog.Infof("RefreshRolePurviewCache refresh success date:%s, all:%v, ann:%v, audit:%v, anal:%v", time.Now().Format(cachedateformat), AllocatRoleIDs, AnnotatRoleIDs, AuditRoleIDs, AnalystRoleIDs)
}
