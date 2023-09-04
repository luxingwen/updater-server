package model

import (
	"encoding/json"
	"mime/multipart"
)

type Pagination struct {
	PageSize int `form:"pageSize"`
	Current  int `form:"current"`
}

func (p *Pagination) GetOffset() int {
	return p.PageSize * (p.Current - 1)
}

type PagedResponse struct {
	Data     interface{} `json:"data"`
	Current  int         `json:"current"`
	PageSize int         `json:"pageSize"`
	Total    int64       `json:"total"`
}

type ReqProgrameQuery struct {
	Pagination
	ProgramName string `form:"programName"`
	TeamId      string `form:"teamId"`
	Uuid        string `json:"uuid"`
}

type ReqVersionQuery struct {
	Pagination
	ProgramUuid string `json:"programUuid"`
}

type ReqProgramActionQuery struct {
	Pagination
	ProgramUuid string `json:"programUuid"`
}

type ReqPackageQuery struct {
	VersionUuid string `json:"versionUuid"`
}

type ReqTaskQuery struct {
	Pagination
	TaskName string   `json:"taskName"`
	TeamId   string   `json:"teamId"`
	Creater  string   `json:"creater"`
	TaskIds  []string `json:"taskIds"`
}

type ReqTaskInfoParam struct {
	TaskId string `json:"taskId"`
}

type TaskRecordInfoParam struct {
	RecordId string `json:"recordId"`
}

type ReqClientQuery struct {
	Pagination
	Uuid     string `json:"uuid"`
	Vmuuid   string `json:"vmuuid"`
	Hostname string `json:"hostname"`
	Ip       string `json:"ip"`
	Sn       string `json:"sn"`
}

type ReqDeletePackageFile struct {
	FileName string `json:"filename"`
}

type UserQuery struct {
	Pagination
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	TeamId   string `json:"teamId"`
}

type ReqTaskRecordeQuery struct {
	Pagination
	TaskId    string   `json:"taskId"`
	RecordIds []string `json:"recordIds"`
}

type ResTaskRecordItem struct {
	RecordId   string `json:"recordId"`
	ClientUuid string `json:"clientUuid"`
}

type ResTaskCreate struct {
	TaskId  string              `json:"taskId"`
	Records []ResTaskRecordItem `json:"records"`
}

type EnvItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ReqScriptContent struct {
	ScriptUuid  string    `json:"scriptUuid"`  // 脚本库uuid
	Content     string    `json:"content"`     // 脚本内容
	WorkDir     string    `json:"workDir"`     // 工作目录
	Params      []string  `json:"params"`      // 参数
	Envs        []EnvItem `json:"env"`         // 环境变量
	Timeout     int       `json:"timeout"`     // 超时时间
	Interpreter string    `json:"interpreter"` // 解释器
	Stdin       string    `json:"stdin"`       // 标准输入
}

// 单个任务请求信息
type ReqTaskSingleCreate struct {
	ClientUuid   string           `json:"clientUuid"`   // 客户端uuid
	Type         string           `json:"type"`         // 任务类型 script | file
	Script       ReqScriptContent `json:"script"`       // 脚本内容
	DownloadFile DownloadRequest  `json:"downloadFile"` // 下载文件
	Creater      string           `json:"creater"`      // 创建人
}

// 获取超时时间
func (self *ReqTaskSingleCreate) GetTimeout() int {

	if self.Type == "script" {
		if self.Script.Timeout == 0 {
			return 60
		}
		return self.Script.Timeout
	}

	if self.Type == "file" {
		if self.DownloadFile.Timeout == 0 {
			return 60
		}
		return self.DownloadFile.Timeout
	}
	return 60
}

// 多个任务请求信息
type ReqTaskMultiCreate struct {
	ClientUuids  []string         `json:"clientUuids"`  // 客户端uuid list
	TaskName     string           `json:"taskName"`     // 任务名称
	Description  string           `json:"description"`  // 任务描述
	Creater      string           `json:"creater"`      // 创建人
	Type         string           `json:"type"`         // 任务类型 script | file
	Script       ReqScriptContent `json:"script"`       // 脚本内容
	DownloadFile DownloadRequest  `json:"downloadFile"` // 下载文件
}

// 获取超时时间
func (self *ReqTaskMultiCreate) GetContentTimeout() int {

	if self.Type == "script" {
		if self.Script.Timeout == 0 {
			return 60
		}
		return self.Script.Timeout
	}

	if self.Type == "file" {
		if self.DownloadFile.Timeout == 0 {
			return 60
		}
		return self.DownloadFile.Timeout
	}
	return 60
}

func (self *ReqTaskMultiCreate) GetContentStr() (r string, err error) {

	if self.Type == "script" {
		b, err := json.Marshal(self.Script)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	if self.Type == "file" {
		b, err := json.Marshal(self.DownloadFile)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	return
}

// 批次任务请求信息
type ReqTaskBatchCreate struct {
	BatchTask BatchTask `json:"batchTask"` // 批次任务
	ReqTaskMultiCreate
}

// Uuid请求参数
type ReqUuidParam struct {
	Uuid string `json:"uuid"`
}

// uuid和uuids请求参数
type ReqUuidsParam struct {
	Uuid  string   `json:"uuid"`
	Uuids []string `json:"uuids"`
}

// 查询预设任务请求参数
type ReqPresetTaskQuery struct {
	Pagination
	Name     string `json:"name"`     // 任务名称
	Category string `json:"category"` // 任务分类
	Type     string `json:"type"`     // 任务类型
	Creater  string `json:"creater"`  // 创建人
}

// 创建预设任务请求参数
type ReqPresetTaskCreate struct {
	Category          string               `json:"category"`          // 任务分类  mutilTask|batchesTask|programActionTask
	MutilTask         ReqTaskMultiCreate   `json:"mutilTask"`         // 多个任务信息
	BatchesTask       ReqTaskBatchCreate   `json:"batchesTask"`       // 批次任务信息
	ProgramActionTask ReqTaskProgramAction `json:"programActionTask"` // 程序行为任务信息
}

// 设置任务信息
func (param *ReqPresetTaskCreate) SetPreTask(preTask *PreTask) {
	if param.Category == "mutilTask" {
		preTask.Name = param.MutilTask.TaskName
		preTask.Type = param.MutilTask.Type
		preTask.Description = param.MutilTask.Description
		preTask.Creater = param.MutilTask.Creater

		b, _ := json.Marshal(param.MutilTask)
		preTask.Content = string(b)
	}

	if param.Category == "batchesTask" {
		preTask.Name = param.BatchesTask.Creater
		preTask.Type = param.BatchesTask.Type
		preTask.Description = param.BatchesTask.Description
		preTask.Creater = param.BatchesTask.Creater

		b, _ := json.Marshal(param.BatchesTask)
		preTask.Content = string(b)
	}

	if param.Category == "programActionTask" {
		preTask.Name = param.ProgramActionTask.Name
		preTask.Type = "program_action"
		preTask.Description = param.ProgramActionTask.Description
		preTask.Creater = param.ProgramActionTask.Creater

		b, _ := json.Marshal(param.ProgramActionTask)
		preTask.Content = string(b)
	}
	return
}

// 更新预设任务请求参数
type ReqPresetTaskUpdate struct {
	Uuid string `json:"uuid"` // 预设任务uuid
	ReqPresetTaskCreate
}

// 脚本库查询请求参数
type ReqScriptLibQuery struct {
	Pagination
	Name    string `json:"name"`    // 脚本库名称
	Creater string `json:"creater"` // 创建人
	TeamId  string `json:"teamId"`  // 团队id
	Type    string `json:"type"`    // 脚本类型
	Status  int    `json:"status"`  // 状态
}

// 文件查询请求参数
type ReqFileQuery struct {
	Pagination
	Name    string `json:"name"`    // 文件名称
	Creater string `json:"creater"` // 创建人
	TeamId  string `json:"teamId"`  // 团队id
	Status  int    `json:"status"`  // 状态
	Type    string `json:"type"`    // 文件类型
	DirUuid string `json:"dirUuid"` // 目录uuid
}

// 文件上传表单参数
type ReqFileUpload struct {
	TeamId  string                `from:"teamId"`  // 团队id
	Creater string                `from:"creater"` // 创建人
	Dir     string                `from:"dir"`     // 目录
	File    *multipart.FileHeader `from:"file"`    // 文件
}

// 创建目录
type ReqDirCreate struct {
	TeamId  string `json:"teamId"`  // 团队id
	Creater string `json:"creater"` // 创建人
	Dir     string `json:"dir"`     // 目录
}

// 危险指令查询参数
type ReqDangerousCommandQuery struct {
	Pagination
	Name string `json:"name"` // 指令名称
}

type ReqDangerousCommandCheck struct {
	Content string `json:"content"` // 指令内容
	CmdType string `json:"cmdType"` // 指令类型
}
