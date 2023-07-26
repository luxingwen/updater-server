package model

import "time"

// 脚本库
type ScriptLibrary struct {
	Id           int       `json:"id"`
	Uuid         string    `json:"uuid"`
	Name         string    `json:"name"`           // 脚本名称
	Content      string    `json:"content"`        // 脚本内容
	Description  string    `json:"description"`    // 脚本描述
	Platform     string    `json:"platform"`       // 脚本平台 linux windows all
	TeamId       int       `json:"team_id"`        // 团队ID
	Creater      string    `json:"creater"`        // 创建人
	Type         string    `json:"type"`           // 类型，私有 公共 共享  private public share
	ShareTeamIds string    `json:"share_team_ids"` // 共享团队ID列表
	Status       int       `json:"status"`         // 状态 0:不可用 1:正常
	Md5          string    `json:"md5"`            // 脚本内容md5
	CreatAt      time.Time `json:"creat_at"`       // 创建时间
	UpdateAt     time.Time `json:"update_at"`      // 更新时间
}
