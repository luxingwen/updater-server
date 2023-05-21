package model

import "time"

// Program 是程序结构体
type Program struct {
	ID          uint            `gorm:"primaryKey"`
	Uuid        string          `json:"uuid" gorm:"primaryKey"`
	ExecUser    string          `json:"execUser"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	TeamID      string          `json:"teamID"`
	Actions     []ProgramAction `json:"actions" gorm:"foreignKey:ProgramUUID"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	Versions    []Version       `json:"versions" gorm:"foreignKey:ProgramID"`
}

type Version struct {
	Uuid        string    `json:"uuid" gorm:"primaryKey"`
	ProgramID   string    `json:"programID"`
	Version     string    `json:"version"`
	ReleaseNote string    `json:"releaseNote"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Packages    []Package `json:"packages" gorm:"foreignKey:ProgramUUID"`
}

// ProgramActionRecord 是程序操作动作记录结构体
type ProgramActionRecord struct {
	ID        uint      `gorm:"primaryKey"`
	ProgramID uint      `gorm:"index"`
	Action    string    `gorm:"type:varchar(255)"`
	Timestamp time.Time `gorm:"index"`
	Client    Client    `gorm:"foreignKey:ClientID"`
}

type Package struct {
	ID           uint      `gorm:"primaryKey"`
	VersionUuid  string    `json:"versionUuid" gorm:"index"`
	Os           string    `json:"os"`
	Arch         string    `json:"architecture"`
	StoragePath  string    `json:"storagePath"`
	DownloadPath string    `json:"downloadPath"`
	Version      string    `json:"version"`
	MD5          string    `json:"md5"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// 脚本执行内容
type ActionContent struct {
	OS      string `json:"os"`
	Content string `json:"content"`
}

// 程序动作
type ProgramAction struct {
	ID          uint      `gorm:"primaryKey"`
	ProgramUUID string    `json:"programUUID" gorm:"index"`
	ActionType  string    `json:"actionType"`
	Content     string    `json:"content"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ProgramActionTemplate struct {
	ID           uint            `gorm:"primaryKey"`
	TemplateName string          `json:"templateName"`
	Actions      []ProgramAction `json:"actions" gorm:"foreignKey:TemplateID"`
	Status       string          `json:"status"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}
