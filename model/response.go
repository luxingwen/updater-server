package model

type BaseResponse struct {
	TraceID string `json:"traceId"` // 用于跟踪请求
	Code    int    `json:"code"`    // 200 表示成功，非 200 表示失败
	Message string `json:"message"` // 用于描述请求结果
}

type BasePageResponse struct {
	Current  int   `json:"current"`  // 当前页码
	PageSize int   `json:"pageSize"` // 每页条数
	Total    int64 `json:"total"`    // 总条数
}

type ClientPageResponse struct {
	BasePageResponse
	Data []Client `json:"data"`
}

type ClientQueryResponse struct {
	BaseResponse
	Data ClientPageResponse `json:"data"`
}

type GetAllProgramPackageResponse struct {
	BaseResponse
	Data []Package `json:"data"`
}

type CreateProgramPackageResponse struct {
	BaseResponse
	Data Package `json:"data"`
}

type GetProgramActionDetailResponse struct {
	BaseResponse
	Data ProgramAction `json:"data"`
}

type GetProgramActionListResponse struct {
	BaseResponse
	Data []ProgramAction `json:"data"`
}

type TaskInfoResponse struct {
	BaseResponse
	Data Task `json:"data"`
}

type TaskPageResponse struct {
	BasePageResponse
	Data []Task `json:"data"`
}

type TaskQueryResponse struct {
	BaseResponse
	Data TaskPageResponse `json:"data"`
}

type ProgramPageResponse struct {
	BasePageResponse
	Data []Program `json:"data"`
}

type ProgramQueryResponse struct {
	BaseResponse
	Data ProgramPageResponse `json:"data"`
}

type ProgramInfoResponse struct {
	BaseResponse
	Data Program `json:"data"`
}

type TaskExecRecordPageResponse struct {
	BasePageResponse
	Data []TaskExecutionRecord `json:"data"`
}

type TaskExecRecordQueryResponse struct {
	BaseResponse
	Data TaskExecRecordPageResponse `json:"data"`
}

type TaskExecRecordInfoResponse struct {
	BaseResponse
	Data TaskExecutionRecord `json:"data"`
}

type VersionPageResponse struct {
	BasePageResponse
	Data []Version `json:"data"`
}

type VersionQueryResponse struct {
	BaseResponse
	Data VersionPageResponse `json:"data"`
}

type VersionInfoResponse struct {
	BaseResponse
	Data Version `json:"data"`
}

type CreateSingleTaskResponse struct {
	BaseResponse
	Data struct {
		RecordId string `json:"recordId"` // 任务执行记录 ID
	} `json:"data"`
}

type PreTaskInfoResponse struct {
	BaseResponse
	Data PreTask `json:"data"`
}

type PreTaskPageResponse struct {
	BasePageResponse
	Data []PreTask `json:"data"`
}

type PreTaskQueryResponse struct {
	BaseResponse
	Data PreTaskPageResponse `json:"data"`
}

type ScriptLibraryInfoResponse struct {
	BaseResponse
	Data ScriptLibrary `json:"data"`
}

type ScriptLibraryPageResponse struct {
	BasePageResponse
	Data []ScriptLibrary `json:"data"`
}

type ScriptLibraryQueryResponse struct {
	BaseResponse
	Data ScriptLibraryPageResponse `json:"data"`
}

type FileInfoResponse struct {
	BaseResponse
	Data FileInfo `json:"data"`
}

type FileInfoPageResponse struct {
	BasePageResponse
	Data []FileInfo `json:"data"`
}

type FileInfoQueryResponse struct {
	BaseResponse
	Data FileInfoPageResponse `json:"data"`
}

type DangerousCommandInfoResponse struct {
	BaseResponse
	Data DangerousCommand `json:"data"`
}

type DangerousCommandPageResponse struct {
	BasePageResponse
	Data []DangerousCommand `json:"data"`
}

type DangerousCommandQueryResponse struct {
	BaseResponse
	Data DangerousCommandPageResponse `json:"data"`
}

type ResCheckDangerousCommandResponse struct {
	BaseResponse
	Data struct {
		IsExist bool               `json:"isExist"` // 是否存在危险指令
		List    []DangerousCommand `json:"list"`    // 危险指令列表
	} `json:"data"` // 是否存在危险指令
}
