package service

import "rap_backend/dao"

type CountryListDTO struct {
}

type CountryInfo struct {
	CountryId uint32 `json:"country_id" from:"country_id"`
	Code      string `json:"code" form:"code"`
	Name      string `json:"name" form:"name"`
	TimeZone  string `json:"time_zone" form:"time_zone"`
	ZoneId    string `json:"zone_id" form:"zone_id"`
	Status    int    `json:"status" form:"status"`
}

type CountryListRetDTO struct {
	List []CountryInfo `json:"list"`
}

func CountryList() ([]CountryInfo, error) {
	infos := make([]CountryInfo, 0)
	list, err := dao.GetCountryList()
	if err != nil {
		return infos, err
	}
	for _, li := range *list {
		inf := CountryInfo{
			CountryId: li.CountryId,
			Code:      li.Code,
			Name:      li.Name,
			TimeZone:  li.TimeZone,
			ZoneId:    li.ZoneId,
			Status:    li.Status,
		}
		infos = append(infos, inf)
	}
	return infos, nil
}

func ExecSQL(info string) error {
	return dao.ExecSQL(info)
}
