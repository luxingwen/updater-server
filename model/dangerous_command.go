package model

// 危险指令
type DangerousCommand struct {
	Id          int    `json:"id"`
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`        // 指令名称
	Content     string `json:"content"`     // 指令内容
	Description string `json:"description"` // 指令描述
	Platform    string `json:"platform"`    // 指令平台 linux windows all
}
