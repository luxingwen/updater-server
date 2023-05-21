package model

import "time"

// TaskTemplate 是任务模板结构体
type TaskTemplate struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Type        string     `json:"type"`
	TeamID      string     `json:"teamID"`
	Tasks       []TaskItem `json:"tasks" gorm:"foreignKey:TaskTemplateID"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
