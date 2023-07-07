package model

import (
	"time"

	"github.com/google/uuid"

	"sort"
)

// Program 是程序结构体
type Program struct {
	ID                 uint            `gorm:"primaryKey"`
	Uuid               string          `json:"uuid" gorm:"column:uuid"`                        // 程序唯一标识
	ExecUser           string          `json:"execUser" gorm:"column:exec_user"`               // 程序执行用户
	Name               string          `json:"name" gorm:"column:name"`                        // 程序名称
	Description        string          `json:"description" gorm:"column:description"`          // 程序描述
	TeamID             string          `json:"teamID" gorm:"column:team_id"`                   // 团队ID
	WindowsInstallPath string          `json:"windowsInstallPath" gorm:"windows_install_path"` // 安装路径
	LinuxInstallPath   string          `json:"linuxInstallPath" gorm:"linux_install_path"`     // 安装路径
	Actions            []ProgramAction `json:"actions" gorm:"foreignKey:ProgramUUID"`          // 程序动作列表
	CreatedAt          time.Time       `json:"createdAt"`                                      // 创建时间
	UpdatedAt          time.Time       `json:"updatedAt"`                                      // 更新时间
	Versions           []Version       `json:"versions" gorm:"foreignKey:ProgramUuid"`         // 程序版本列表
}

func (Program) TableName() string {
	return "program"
}

type Version struct {
	ID          uint      `gorm:"primaryKey"`                             // 程序版本ID
	Uuid        string    `json:"uuid" gorm:"column:uuid"`                // 程序版本UUID
	ProgramUuid string    `json:"programUuid"`                            // 程序UUID
	Version     string    `json:"version"`                                // 程序版本
	ReleaseNote string    `json:"releaseNote"`                            // 程序版本发布说明
	CreatedAt   time.Time `json:"createdAt"`                              // 创建时间
	UpdatedAt   time.Time `json:"updatedAt"`                              // 更新时间
	Packages    []Package `json:"packages" gorm:"foreignKey:VersionUuid"` // 程序包列表
}

type Package struct {
	ID           uint      `gorm:"primaryKey"`               // 程序包ID
	Uuid         string    `json:"uuid" gorm:"primaryKey"`   // 程序包UUID
	VersionUuid  string    `json:"versionUuid" gorm:"index"` // 程序版本UUID
	Os           string    `json:"os"`                       // 程序包操作系统
	Arch         string    `json:"arch"`                     // 程序包架构
	StoragePath  string    `json:"storagePath"`              // 程序包存储路径
	DownloadPath string    `json:"downloadPath"`             // 程序包下载路径
	MD5          string    `json:"md5"`                      // 程序包MD5
	CreatedAt    time.Time `json:"createdAt"`                // 创建时间
	UpdatedAt    time.Time `json:"updatedAt"`                // 更新时间
}

// 程序动作
type ProgramAction struct {
	ID          uint       `gorm:"primaryKey"`               // 程序动作ID
	Uuid        string     `json:"uuid" gorm:"primaryKey"`   // 程序动作UUID
	ProgramUUID string     `json:"programUUID" gorm:"index"` // 程序UUID
	Name        string     `json:"name"`                     // 程序动作名称
	ActionType  ActionType `json:"actionType"`               // 程序动作类型
	Content     string     `json:"content"`                  // 程序动作内容
	Status      string     `json:"status"`                   // 程序动作状态
	Description string     `json:"description""`             // 程序动作描述
	CreatedAt   time.Time  `json:"createdAt"`                // 创建时间
	UpdatedAt   time.Time  `json:"updatedAt"`                // 更新时间
}

func (ProgramAction) TableName() string {
	return "program_action"
}

type TemplateAction struct {
	ProgramActionUuid string `json:"programActionUuid"` // 程序动作UUID
	Sequence          int    `json:"sequence"`          // 程序动作顺序
	Uuid              string `json:"uuid"`              // 模板动作UUID
	NextUuid          string `json:"nextUuid"`          // 下一个模板动作UUID
	TaskRecordId      string `json:"taskRecordId"`      // 任务记录ID
	NextTaskRecordId  string `json:"nextTaskRecordId"`  // 下一个任务记录ID
}

func GenerateTemplateActionNextUuids(actions []*TemplateAction) {

	for i := 0; i < len(actions); i++ {
		actions[i].TaskRecordId = uuid.New().String()
	}

	// 按照 Sequence 进行排序
	sort.Slice(actions, func(i, j int) bool {
		return actions[i].Sequence < actions[j].Sequence
	})

	for i := 0; i < len(actions); i++ {
		if i == len(actions)-1 {
			actions[i].NextTaskRecordId = ""
		} else {
			actions[i].NextTaskRecordId = actions[i+1].TaskRecordId
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
		Content: `[{"os": "linux", "content": "echo hello"}, 
		          {"os": "windows", "content": "echo hello"}]`,
		Status:      "待处理",
		Description: "安装动作的描述信息",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		Uuid:       uuid.New().String(),
		Name:       "启动",
		ActionType: "Start",
		Content: `[{"os": "linux", "content": "echo hello"}, 
		          {"os": "windows", "content": "echo hello"}]`,
		Status:      "待处理",
		Description: "启动动作的描述信息",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		Uuid:       uuid.New().String(),
		Name:       "停止",
		ActionType: "Stop",
		Content: `[{"os": "linux", "content": "echo hello"}, 
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
