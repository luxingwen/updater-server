package model

import "time"

type TaskItem struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Content     string    `json:"content"`
	Status      string    `json:"status"`
	ItemUUID    string    `json:"itemUUID"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	TeamID      string    `json:"teamID"`
	Platform    string    `json:"platform"`
}
