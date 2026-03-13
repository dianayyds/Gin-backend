package dao

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	dbutils "rap_backend/db"
	"time"
)

type LabelTemplate struct {
	ID        int64          `json:"id" gorm:"column:id" form:"id"`                                       // 主键
	Name      string         `json:"name" gorm:"column:name" form:"name"`                                 // 模版名字
	LabelIds  string         `json:"label_ids" from:"label_ids" gorm:"column:label_ids" form:"label_ids"` //标签数组
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at"`                                 // 创建时间
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at"`                                 // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *LabelTemplate) TableName() string {
	return "label_template"
}

type LabelTemplateList []*LabelTemplate

type LabelTemplateModel struct {
	ctx context.Context
	db  *gorm.DB
}

func GetLabelTemplate(db *gorm.DB) *LabelTemplateModel {
	db = dbutils.GetTxDB(context.Background(), db)
	model := &LabelTemplateModel{
		db: db,
	}
	return model
}

type LabelTemplateParams struct {
	LabelTemplate
	*Paging
	IdList  []int64 `json:"id_list" form:"id_list"`
	OrderBy string
}

func (m *LabelTemplateModel) setParams(p *LabelTemplateParams) *gorm.DB {
	db := m.db
	if p == nil {
		return db
	}
	if p.ID != 0 {
		db = db.Where("id = ?", p.ID)
	}
	if len(p.IdList) > 0 {
		db = db.Where("id in (?)", p.IdList)
	}
	if p.Name != "" {
		db = db.Where("name like ?", fmt.Sprintf("%%%s%%", p.Name))
	}
	if p.OrderBy != "" {
		db = db.Order(p.OrderBy)
	} else {
		db = db.Order("id desc")
	}
	if p.Paging != nil && p.Paging.Page != -1 && p.Paging.PageSize != -1 {
		db = db.Offset(p.Paging.Offset()).Limit(p.Paging.Limit())
	}
	return db
}

func (m *LabelTemplateModel) GetByName(name string) (*LabelTemplate, error) {
	row := new(LabelTemplate)
	err := m.db.Where("name = ?", name).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return row, nil
}

// Get 取得一条记录
func (m *LabelTemplateModel) Get(params *LabelTemplateParams) (*LabelTemplate, error) {
	row := new(LabelTemplate)
	err := m.setParams(params).Take(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (m *LabelTemplateModel) Save(tag *LabelTemplate) error {
	if err := m.db.Save(tag).Error; err != nil {
		return fmt.Errorf("save label template model err: %+v", err)
	}
	return nil
}

func (m *LabelTemplateModel) Delete(params *LabelTemplateParams) error {
	if err := m.setParams(params).Delete(&LabelTemplate{}).Error; err != nil {
		return fmt.Errorf("delete label template model err: %+v", err)
	}
	return nil
}

// Count 取得总行数
func (m *LabelTemplateModel) Count(params *LabelTemplateParams) (int64, error) {
	var total int64 = 0
	if err := m.setParams(params).Model(&LabelTemplate{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// List 取得列表
func (m *LabelTemplateModel) List(params *LabelTemplateParams) (LabelTemplateList, error) {
	rows := make(LabelTemplateList, 0)
	sql := m.setParams(params)
	if err := sql.Find(&rows).Error; err != nil {
		return nil, fmt.Errorf("list label template model err: %+v", err)
	}
	return rows, nil
}

func (m *LabelTemplateModel) UpdateColumn(params *LabelTemplateParams, field string, value interface{}) error {
	return m.setParams(params).Model(LabelTemplate{}).UpdateColumn(field, value).Error
}

func (m *LabelTemplateModel) UpdateById(id int64, do *LabelTemplate) error {
	return m.db.Model(LabelTemplate{}).Where("id = ?", id).Updates(do).Error
}
