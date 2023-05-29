package model

import (
	"time"

	"github.com/google/uuid"

	"sort"
)

// Program 是程序结构体
type Program struct {
	ID                 uint            `gorm:"primaryKey"`
	Uuid               string          `json:"uuid" gorm:"column:uuid"`
	ExecUser           string          `json:"execUser" gorm:"column:exec_user"`
	Name               string          `json:"name" gorm:"column:name"`
	Description        string          `json:"description" gorm:"column:description"`
	TeamID             string          `json:"teamID" gorm:"column:team_id"`
	WindowsInstallPath string          `json:"windowsInstallPath" gorm:"windows_install_path"` // 安装路径
	LinuxInstallPath   string          `json:"linuxInstallPath" gorm:"linux_install_path"`     // 安装路径
	Actions            []ProgramAction `json:"actions" gorm:"foreignKey:ProgramUUID"`
	CreatedAt          time.Time       `json:"createdAt"`
	UpdatedAt          time.Time       `json:"updatedAt"`
	Versions           []Version       `json:"versions" gorm:"foreignKey:ProgramUuid"`
}

func (Program) TableName() string {
	return "program"
}

type Version struct {
	ID          uint      `gorm:"primaryKey"`
	Uuid        string    `json:"uuid" gorm:"column:uuid"`
	ProgramUuid string    `json:"programUuid"`
	Version     string    `json:"version"`
	ReleaseNote string    `json:"releaseNote"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Packages    []Package `json:"packages" gorm:"foreignKey:VersionUuid"`
}

type Package struct {
	ID           uint      `gorm:"primaryKey"`
	Uuid         string    `json:"uuid" gorm:"primaryKey"`
	VersionUuid  string    `json:"versionUuid" gorm:"index"`
	Os           string    `json:"os"`
	Arch         string    `json:"arch"`
	StoragePath  string    `json:"storagePath"`
	DownloadPath string    `json:"downloadPath"`
	MD5          string    `json:"md5"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// 程序动作
type ProgramAction struct {
	ID          uint       `gorm:"primaryKey"`
	Uuid        string     `json:"uuid" gorm:"primaryKey"`
	ProgramUUID string     `json:"programUUID" gorm:"index"`
	Name        string     `json:"name"`
	ActionType  ActionType `json:"actionType"`
	Content     string     `json:"content"`
	Status      string     `json:"status"`
	Description string     `json:"description""`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func (ProgramAction) TableName() string {
	return "program_action"
}

type TemplateAction struct {
	ProgramActionUuid string `json:"programActionUuid"`
	Sequence          int    `json:"sequence"`
	Uuid              string `json:"uuid"`
	NextUuid          string `json:"nextUuid"`
}

func GenerateTemplateActionNextUuids(actions []*TemplateAction) {

	for i := 0; i < len(actions); i++ {
		actions[i].Uuid = uuid.New().String()
	}

	// 按照 Sequence 进行排序
	sort.Slice(actions, func(i, j int) bool {
		return actions[i].Sequence < actions[j].Sequence
	})

	for i := 0; i < len(actions); i++ {
		if i == len(actions)-1 {
			actions[i].NextUuid = ""
		} else {
			actions[i].NextUuid = actions[i+1].Uuid
		}
	}
}

type TemplateActionDetail struct {
	TemplateAction
	ProgramAction
}

// ProgramActionRecord 是程序操作动作记录结构体
type ProgramActionRecord struct {
	ID        uint      `gorm:"primaryKey"`
	ProgramID uint      `gorm:"index"`
	Action    string    `gorm:"type:varchar(255)"`
	Timestamp time.Time `gorm:"index"`
	Client    Client    `gorm:"foreignKey:ClientID"`
}

// 脚本执行内容
type ActionContent struct {
	OS      string `json:"os"`
	Content string `json:"content"`
}

type ActionType string

const (
	ActionTypeDownload  ActionType = "Download"
	ActionTypeInstall   ActionType = "Install"
	ActionTypeStart     ActionType = "Start"
	ActionTypeStop      ActionType = "Stop"
	ActionTypeUninstall ActionType = "Uninstall"
	ActionTypeBackup    ActionType = "Backup"
	ActionTypeStatus    ActionType = "Status"
	ActionTypeVersion   ActionType = "Version"
	ActionTypeSingle    ActionType = "Single"
	ActionTypeComposite ActionType = "Composite"
)

var InitialActions = []ProgramAction{
	{
		Uuid:       uuid.New().String(),
		Name:       "下载",
		ActionType: "Download",
		Content: `[{"os": "linux", "content": "Linux 下载内容"}, 
		           {"os": "windows", "content": "Windows 下载内容"}]`,
		Status:      "待处理",
		Description: "下载动作的描述信息",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		Uuid:       uuid.New().String(),
		Name:       "安装",
		ActionType: "Install",
		Content: `[{"os": "linux", "content": "Linux 安装内容"}, 
		          {"os": "windows", "content": "Windows 安装内容"}]`,
		Status:      "待处理",
		Description: "安装动作的描述信息",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		Uuid:       uuid.New().String(),
		Name:       "启动",
		ActionType: "Start",
		Content: `[{"os": "linux", "content": "Linux 启动内容"}, 
		          {"os": "windows", "content": "Windows 启动内容"}]`,
		Status:      "待处理",
		Description: "启动动作的描述信息",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		Uuid:       uuid.New().String(),
		Name:       "停止",
		ActionType: "Stop",
		Content: `[{"os": "linux", "content": "Linux 停止内容"}, 
		          {"os": "windows", "content": "Windows 停止内容"}]`,
		Status:      "待处理",
		Description: "停止动作的描述信息",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		Uuid:       uuid.New().String(),
		Name:       "卸载",
		ActionType: "Uninstall",
		Content: `[{"os": "linux", "content": "Linux 卸载内容"}, 
		          {"os": "windows", "content": "Windows 卸载内容"}]`,
		Status:      "待处理",
		Description: "卸载动作的描述信息",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		Uuid:       uuid.New().String(),
		Name:       "备份",
		ActionType: "Backup",
		Content: `[{"os": "linux", "content": "Linux 备份内容"}, 
		          {"os": "windows", "content": "Windows 备份内容"}]`,
		Status:      "待处理",
		Description: "备份动作的描述信息",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		Uuid:       uuid.New().String(),
		Name:       "状态",
		ActionType: "Status",
		Content: `[{"os": "linux", "content": "Linux 状态内容"}, 
		          {"os": "windows", "content": "Windows 状态内容"}]`,
		Status:      "待处理",
		Description: "状态动作的描述信息",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		Uuid:       uuid.New().String(),
		Name:       "版本",
		ActionType: "Version",
		Content: `[{"os": "linux", "content": "Linux 版本内容"}, 
		          {"os": "windows", "content": "Windows 版本内容"}]`,
		Status:      "待处理",
		Description: "版本动作的描述信息",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		Uuid:       uuid.New().String(),
		Name:       "单一",
		ActionType: "Single",
		Content: `[{"os": "linux", "content": "Linux 单一内容"}, 
		          {"os": "windows", "content": "Windows 单一内容"}]`,
		Status:      "待处理",
		Description: "单一动作的描述信息",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		Uuid:       uuid.New().String(),
		Name:       "组合",
		ActionType: "Composite",
		Content: `[{"os": "linux", "content": "Linux 组合内容"}, 
		          {"os": "windows", "content": "Windows 组合内容"}]`,
		Status:      "待处理",
		Description: "组合动作的描述信息",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
}
