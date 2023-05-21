package model

import "time"

// TaskExecutionRecord 是任务执行记录结构体
type TaskExecutionRecord struct {
	TaskID      string    `json:"taskID"`
	TaskType    string    `json:"taskType"`
	Description string    `json:"description"`
	TeamID      string    `json:"teamID"`
	Status      string    `json:"status"`
	ParentID    string    `json:"parentID"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	LogOutput   string    `json:"logOutput"`
	TaskContent string    `json:"taskContent"`
	Creater     string    `json:"creater"`
	// 其他任务执行记录字段...
}
