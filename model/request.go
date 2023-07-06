package model

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
	ClientUuid string           `json:"clientUuid"` // 客户端uuid
	Type       string           `json:"type"`       // 任务类型 script | file
	Script     ReqScriptContent `json:"script"`     // 脚本内容
	Creater    string           `json:"creater"`    // 创建人
}

// 多个任务请求信息
type ReqTaskMultiCreate struct {
	ClientUuids []string         `json:"clientUuids"` // 客户端uuid list
	TaskName    string           `json:"taskName"`    // 任务名称
	Description string           `json:"description"` // 任务描述
	Creater     string           `json:"creater"`     // 创建人
	Type        string           `json:"type"`        // 任务类型 script | file
	Script      ReqScriptContent `json:"script"`      // 脚本内容
}

// 批次任务请求信息
type ReqTaskBatchCreate struct {
	//BatchTask BatchTask `json:"batchTask"` // 批次任务
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
