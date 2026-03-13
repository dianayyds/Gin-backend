package label_template

import "rap_backend/dao"

type LabelTemplate struct {
	dao.LabelTemplate
}

type GetLabelTemplateReq struct {
	*dao.LabelTemplateParams
}

type GetLabelTemplateResp struct {
	List       dao.LabelTemplateList `json:"list"`
	TotalCount int64                 `json:"total_count"`
}
type CreateLabelTemplateReq LabelTemplate
type UpdateLabelTemplateReq LabelTemplate
type DeleteLabelTemplateReq struct {
	*dao.LabelTemplateParams
}
