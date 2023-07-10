package model

import "time"

type FileInfo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Uuid      string    `gorm:"column:uuid" json:"uuid"`            // 文件UUID
	Filename  string    `gorm:"column:filename" json:"filename"`    // 文件名
	Extension string    `gorm:"column:extension" json:"extension"`  // 文件扩展名
	Path      string    `gorm:"column:path" json:"path"`            // 文件路径
	StorePath string    `gorm:"column:store_path" json:"storePath"` // 文件存储路径
	Size      int64     `gorm:"column:size" json:"size"`            // 文件大小
	IsDir     bool      `gorm:"column:is_dir" json:"isDir"`         // 是否是目录
	Type      string    `gorm:"column:type" json:"type"`            // 文件类型
	Status    string    `gorm:"column:status" json:"status"`        // 状态 1: 正常 2: 删除
	Md5       string    `gorm:"column:md5" json:"md5"`              // 文件MD5
	Creater   string    `gorm:"column:creater" json:"creater"`      // 创建者
	TeamId    string    `gorm:"column:team_id" json:"teamId"`       // 团队ID
	ParentId  string    `gorm:"column:parent_id" json:"parentId"`   // 父级ID
	CreateAt  time.Time `gorm:"column:create_at" json:"create_at"`  // 创建时间
	UpdateAt  time.Time `gorm:"column:update_at" json:"update_at"`  // 更新时间
}
