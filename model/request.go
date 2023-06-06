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
