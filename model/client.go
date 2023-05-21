package model

import "time"

type Client struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	VMUUID   string    `gorm:"column:vmuuid" json:"vmuuid"`
	SN       string    `gorm:"column:sn" json:"sn"`
	Hostname string    `gorm:"column:hostname" json:"hostname"`
	IP       string    `gorm:"column:ip" json:"ip"`
	ProxyID  string    `gorm:"column:proxy_id" json:"proxyID"`
	Status   string    `gorm:"column:status" json:"status"`
	OS       string    `gorm:"column:os" json:"os"`
	Arch     string    `gorm:"column:arch" json:"arch"`
	Created  time.Time `gorm:"column:created" json:"created"`
	Updated  time.Time `gorm:"column:updated" json:"updated"`
}
