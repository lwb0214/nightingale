package models

import (
	"time"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

// MessageTemplate 消息模板结构
type MessageTemplate struct {
	ID           uint              `json:"id" gorm:"primarykey"`
	Name         string            `json:"name"`                           // 模板名称
	Ident        string            `json:"ident"`                          // 模板标识
	Content      map[string]string `json:"content" gorm:"serializer:json"` // 模板内容
	UserGroupIds []int64           `json:"user_group_ids" gorm:"serializer:json"`
	Private      int               `json:"private"` // 0-公开 1-私有
	CreateAt     int64             `json:"create_at"`
	CreateBy     string            `json:"create_by"`
	UpdateAt     int64             `json:"update_at"`
	UpdateBy     string            `json:"update_by"`
}

type HTTPConfig struct {
	Type     string `json:"type"`
	IsGlobal bool   `json:"is_global"`
	Name     string `json:"name"`
	Ident    string `json:"ident"`
	Note     string `json:"note"` // 备注

	Enabled     bool              `json:"enabled"`     // 是否启用
	URL         string            `json:"url"`         // 回调URL
	Method      string            `json:"method"`      // HTTP方法
	Headers     map[string]string `json:"headers"`     // 请求头
	Timeout     int               `json:"timeout"`     // 超时时间(毫秒)
	Concurrency int               `json:"concurrency"` // 并发度
	RetryTimes  int               `json:"retryTimes"`  // 重试次数
	RetryDelay  int               `json:"retryDelay"`  // 重试间隔(毫秒)
	SkipVerify  bool              `json:"skipVerify"`  // 跳过SSL校验
	Proxy       string            `json:"proxy"`       // 代理地址

	// 请求参数配置
	EnableParams bool              `json:"enableParams"` // 启用Params参数
	Params       map[string]string `json:"params"`       // URL参数

	// 请求体配置
	EnableBody bool   `json:"enableBody"` // 启用Body
	Body       string `json:"body"`       // 请求体内容
}

func (t *MessageTemplate) TableName() string {
	return "message_template"
}

func (t *MessageTemplate) Update(ctx *ctx.Context, ref MessageTemplate) error {
	// ref.FE2DB()

	ref.ID = t.ID
	ref.CreateAt = t.CreateAt
	ref.CreateBy = t.CreateBy
	ref.UpdateAt = time.Now().Unix()

	// err := ref.Verify()
	// if err != nil {
	// 	return err
	// }
	return DB(ctx).Model(t).Select("*").Updates(ref).Error
}

func MessageTemplateGet(ctx *ctx.Context, where string, args ...interface{}) (*MessageTemplate, error) {
	lst, err := MessageTemplatesGet(ctx, where, args...)
	if err != nil || len(lst) == 0 {
		return nil, err
	}
	return lst[0], err
}

func MessageTemplatesGet(ctx *ctx.Context, where string, args ...interface{}) ([]*MessageTemplate, error) {
	lst := make([]*MessageTemplate, 0)
	session := DB(ctx)
	if where != "" && len(args) > 0 {
		session = session.Where(where, args...)
	}

	err := session.Debug().Find(&lst).Error
	if err != nil {
		return nil, err
	}
	return lst, nil
}
