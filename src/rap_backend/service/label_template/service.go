package label_template

import (
	"errors"
	"rap_backend/dao"
	"sync"
)

var labelTemplateServiceOnce sync.Once

var labelTemplateService LabelTemplateServiceI

type labelTemplateServier struct {
}

type LabelTemplateServiceI interface {
	GetLabelTemplateList(req GetLabelTemplateReq) (*GetLabelTemplateResp, error)
	CreateLabelTemplate(req CreateLabelTemplateReq) error
	UpdateLabelTemplateById(id int64, req UpdateLabelTemplateReq) error
	DeleteLabelTemplate(req DeleteLabelTemplateReq) error
}

func GetLabelTemplateSvc() LabelTemplateServiceI {
	labelTemplateServiceOnce.Do(func() {
		if labelTemplateService == nil {
			labelTemplateService = &labelTemplateServier{}
		}
	})
	return labelTemplateService
}

func (s *labelTemplateServier) GetLabelTemplateList(req GetLabelTemplateReq) (*GetLabelTemplateResp, error) {
	cnt, err := dao.GetLabelTemplate(nil).Count(req.LabelTemplateParams)
	if err != nil {
		return nil, err
	}
	list, err := dao.GetLabelTemplate(nil).List(req.LabelTemplateParams)
	if err != nil {
		return nil, err
	}
	return &GetLabelTemplateResp{
		TotalCount: cnt,
		List:       list,
	}, nil
}

func (s *labelTemplateServier) CreateLabelTemplate(req CreateLabelTemplateReq) error {
	_, err := dao.GetLabelTemplate(nil).GetByName(req.Name)
	if err == nil {
		return errors.New("This name already exists.")
	}
	return dao.GetLabelTemplate(nil).Save(&req.LabelTemplate)
}

func (s *labelTemplateServier) UpdateLabelTemplateById(id int64, req UpdateLabelTemplateReq) error {
	return dao.GetLabelTemplate(nil).UpdateById(id, &req.LabelTemplate)
}

func (s *labelTemplateServier) DeleteLabelTemplate(req DeleteLabelTemplateReq) error {
	return dao.GetLabelTemplate(nil).Delete(req.LabelTemplateParams)
}
