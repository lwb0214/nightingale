package models

import (
	"time"

	"github.com/ccfos/nightingale/v6/pkg/ctx"
)

type NotifyRule struct {
	ID uint `json:"id" gorm:"primarykey"`

	// 基本配置
	Name         string  `json:"name"`                                  // 名称
	Description  string  `json:"description"`                           // 备注
	Enable       bool    `json:"enable"`                                // 启用状态
	UserGroupIds []int64 `json:"user_group_ids" gorm:"serializer:json"` // 告警组ID

	// 通知配置
	NotifyConfigs []NotifyConfig `json:"notify_configs" gorm:"serializer:json"`

	CreateAt int64  `json:"create_at"`
	CreateBy string `json:"create_by"`
	UpdateAt int64  `json:"update_at"`
	UpdateBy string `json:"update_by"`
}

type NotifyConfig struct {
	Channel  string      `json:"channel"`  // 通知媒介(如：阿里云短信)
	Template string      `json:"template"` // 通知模板
	Params   interface{} `json:"params"`   // 通知参数

	Severities []int         `json:"severities"`  // 适用级别(一级告警、二级告警、三级告警)
	TimeRanges []TimeRanges  `json:"time_ranges"` // 适用时段
	LabelKeys  []LabelFilter `json:"label_keys"`  // 适用标签
}

type TimeRanges struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Week  string `json:"week"`
}

type LabelFilter struct {
	Key   string `json:"key"`
	Op    string `json:"op"` // == != in not in =~ !~
	Value string `json:"value"`
}

func (r *NotifyRule) Update(ctx *ctx.Context, ref NotifyRule) error {
	// ref.FE2DB()

	ref.ID = r.ID
	ref.CreateAt = r.CreateAt
	ref.CreateBy = r.CreateBy
	ref.UpdateAt = time.Now().Unix()

	// err := ref.Verify()
	// if err != nil {
	// 	return err
	// }
	return DB(ctx).Model(r).Select("*").Updates(ref).Error
}

func NotifyRuleGet(ctx *ctx.Context, where string, args ...interface{}) (*NotifyRule, error) {
	lst, err := NotifyRulesGet(ctx, where, args...)
	if err != nil || len(lst) == 0 {
		return nil, err
	}
	return lst[0], err
}

func NotifyRulesGet(ctx *ctx.Context, where string, args ...interface{}) ([]*NotifyRule, error) {
	lst := make([]*NotifyRule, 0)
	session := DB(ctx)
	if where != "" && len(args) > 0 {
		session = session.Where(where, args...)
	}
	err := session.Find(&lst).Error
	if err != nil {
		return nil, err
	}
	return lst, nil
}
