package model

import "time"

type Client struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Uuid     string    `gorm:"column:uuid" json:"uuid"`         // 客户端UUID
	VMUUID   string    `gorm:"column:vmuuid" json:"vmuuid"`     // 虚拟机UUID
	SN       string    `gorm:"column:sn" json:"sn"`             // 序列号
	Hostname string    `gorm:"column:hostname" json:"hostname"` // 主机名
	IP       string    `gorm:"column:ip" json:"ip"`             // IP
	ProxyID  string    `gorm:"column:proxy_id" json:"proxyID"`  // 代理ID
	Status   string    `gorm:"column:status" json:"status"`     // 状态
	OS       string    `gorm:"column:os" json:"os"`             // 操作系统
	Arch     string    `gorm:"column:arch" json:"arch"`         // 架构
	Version  string    `gorm:"column:version" json:"version"`   // 版本
	Created  time.Time `gorm:"column:created" json:"created"`
	Updated  time.Time `gorm:"column:updated" json:"updated"`
}
