package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Id           uint      `gorm:"primaryKey"`
	TaskID       string    `json:"taskId"`       //	任务ID
	TaskName     string    `json:"taskName"`     // 任务名称
	TaskType     string    `json:"taskType"`     // 任务类型
	Category     string    `json:"category"`     // 任务类别 root
	TaskStatus   string    `json:"taskStatus"`   // 任务状态
	ParentTaskID string    `json:"parentTaskId"` // 父任务ID
	Content      string    `json:"content"`      // 任务内容
	Description  string    `json:"description"`  // 任务描述
	Creater      string    `json:"creater"`      // 创建者
	TeamID       string    `json:"teamId"`       // 团队ID
	NextTaskID   string    `json:"nextTaskId"`   // 下一个任务ID
	Ext          string    `json:"ext"`          // 扩展字段
	TraceId      string    `json:"traceId"`      // 跟踪id
	CreatedAt    time.Time `gorm:"column:created_at" json:"created"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated"`
}

type TaskExecutionRecord struct {
	Id             uint      `gorm:"primaryKey"`
	RecordID       string    `json:"recordId"`       // 执行记录ID
	TaskID         string    `json:"taskId"`         // 任务ID
	ClientUUID     string    `json:"clientUuid"`     // 客户端UUID
	Category       string    `json:"category"`       // 任务类别
	Name           string    `json:"name"`           // 任务名称
	TaskType       string    `json:"taskType"`       // 任务类型
	Status         string    `json:"status"`         // 任务状态
	StartTime      time.Time `json:"startTime"`      // 任务开始时间
	EndTime        time.Time `json:"endTime"`        // 任务结束时间
	Stdout         string    `json:"stdout"`         // 标准输出
	Stderr         string    `json:"stderr"`         // 标准错误
	Message        string    `json:"message"`        // 消息
	ScriptExitCode int       `json:"scriptExitCode"` // 脚本退出码
	Code           string    `json:"code"`           // 任务执行码
	Content        string    `json:"content"`        // 任务内容
	Timeout        int       `json:"timeout"`        // 任务超时时间
	ParentRecordID string    `json:"parentRecordId"` // 父任务执行记录ID
	NextRecordID   string    `json:"nextRecordId"`   // 下一个任务执行记录ID
	TraceId        string    `json:"traceId"`        // 跟踪id
	CreatedAt      time.Time `gorm:"column:created_at" json:"created"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated"`
}

type HostInfo struct {
	All     bool     `json:"all"`
	Clients []string `json:"clients"`
}

type BatchTask struct {
	count  int
	Style  string `json:"style"`
	Number int    `json:"number"`
}

type TaskBatchesInfo struct {
	TaskID     string   `json:"task_id"`
	Total      int      `json:"total"`
	NextTaskID string   `json:"next_task_id"`
	Sequence   int      `json:"sequence"`
	Clients    []string `json:"clients"`
}

type TaskContentInDB struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (self BatchTask) GenerateTaskBatchesInfo(clients []string) (r []TaskBatchesInfo) {
	bs := self.GetBatchesList(len(clients))
	for i := 0; i < len(bs); i++ {
		r = append(r, TaskBatchesInfo{
			TaskID:     uuid.New().String(),
			Total:      bs[i],
			NextTaskID: "",
			Sequence:   i + 1,
			Clients:    clients[:bs[i]],
		})
		clients = clients[bs[i]:]
	}

	for i := 0; i < len(r); i++ {
		if i == len(r)-1 {
			r[i].NextTaskID = ""
		} else {
			r[i].NextTaskID = r[i+1].TaskID
		}
	}
	return
}

func (self BatchTask) GetBatchesList(total int) (r []int) {
	if self.Style == "average" {
		return self.GetAverage(total)
	}
	for self.count < total {
		num := self.GetNumber(total)
		if self.count+num >= total {
			num = total - self.count
		}
		r = append(r, num)
		self.count += num
	}
	return
}

func (self BatchTask) GetNumber(total int) int {
	if self.Style == "normal" {
		if self.Number == 0 {
			return total
		}
		return self.Number
	}
	return total
}

func (self BatchTask) GetAverage(total int) (r []int) {
	if self.Number <= 0 {
		self.Number = 1
	}

	if total <= self.Number {
		self.Number = total
	}
	r = make([]int, self.Number)

	for i := 0; i < self.Number; i++ {
		r[i] = 0
	}
	for total > 0 {
		for i := 0; i < self.Number && total > 0; i++ {
			r[i]++
			total--
		}
	}
	return
}

// 是否是批次任务
func (self BatchTask) IsBatchTask() bool {
	if self.Style == "" || self.Style == "all" {
		return false
	}
	return true
}

type ProgramActionTask struct {
	ProgramUuid        string `json:"programUuid"`        // 程序uuid
	ProgramActionUuid  string `json:"programActionUuid"`  // 程序动作uuid
	ProgramVersionUuid string `json:"programVersionUuid"` // 程序版本uuid
	Timeout            int    `json:"timeout"`            // 超时时间
}

type ReqTask struct {
	Name        string    `json:"name"`
	Creater     string    `json:"creater"`
	TeamID      string    `json:"team_id"`
	Description string    `json:"description"`
	HostInfo    HostInfo  `json:"hostInfo"`
	BatchTask   BatchTask `json:"batchTask"`
}

type ReqTaskProgramAction struct {
	ReqTask
	Content ProgramActionTask `json:"content"`
}

type TaskContent struct {
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type TaskContentInfo struct {
	TaskID       string `json:"taskId"`
	Sequence     int    `json:"sequence"`
	TaskRecordId string `json:"taskRecordId"`
}

type TaskType string

const (
	TaskTypeScript        TaskType = "script"
	TaskTypeProgram       TaskType = "program"
	TaskTypeProgramAction TaskType = "programAction"
	TaskTypeBatch         TaskType = "batch"
)

type TaskCategory string

const (
	TaskCategoryRoot TaskCategory = "root"
	TaskCategorySub  TaskCategory = "sub"
)

const (
	TaskStatusPreparing = "Preparing" // 准备中
	TaskStatusReady     = "Ready"     // 准备完成,就绪状态，等待执行
	TaskStatusRunning   = "Running"   // 执行中
	TaskStatusPaused    = "Paused"    // 暂停
	TaskStatusAbandoned = "Abandoned" // 废弃
	TaskStatusCompleted = "Completed" // 完成
	TaskStatusSuceess   = "Success"   // 成功
	TaskStatusFailed    = "Failed"    // 失败
)

// DownloadRequest 是下载请求参数
type DownloadRequest struct {
	DownLoadPath     string `json:"downloadPath"`     // 下载路径
	URL              string `json:"url"`              // 下载 URL
	DestPath         string `json:"destPath"`         // 目标路径
	AutoCreateDir    bool   `json:"autoCreateDir"`    // 是否自动创建文件夹
	OverwriteExisted bool   `json:"overwriteExisted"` // 文件存在是否覆盖文件
	Timeout          int    `json:"timeout"`          // 超时时间
}

type FileDownloadTask struct {
	Content DownloadRequest `json:"content"`
	Clients []string        `json:"clients"`
	Name    string          `json:"name"`
	Creater string          `json:"creater"`
}
