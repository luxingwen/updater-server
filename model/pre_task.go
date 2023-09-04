package model

import "time"

// 预置任务
type PreTask struct {
	Id          int       `json:"id"`
	Uuid        string    `json:"uuid"`        // 任务UUID
	Name        string    `json:"name"`        // 任务名称
	Content     string    `json:"content"`     // 任务内容
	Description string    `json:"description"` // 任务描述
	Type        string    `json:"type"`        // 任务类型
	Category    string    `json:"category"`    // 任务分类
	Creater     string    `json:"creater"`     // 创建人
	Status      int       `json:"status"`      // 状态 0:不可用 1:正常 2:删除
	TeamId      string    `json:"team_id"`     // 团队ID
	CreateAt    time.Time `json:"create_at"`   // 创建时间
	UpdateAt    time.Time `json:"update_at"`   // 更新时间
}
