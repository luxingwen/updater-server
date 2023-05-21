package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"updater-server/model"

	"updater-server/pkg/app"
)

type ClientService struct {
}

func NewClientService() *ClientService {
	return &ClientService{}
}

func (s *ClientService) CreateClient(ctx *app.Context, client *model.Client) error {
	client.Created = time.Now()
	client.Updated = time.Now()

	err := ctx.DB.Create(client).Error
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	return nil
}

func (s *ClientService) GetClientByID(ctx *app.Context, id uint) (*model.Client, error) {
	client := &model.Client{}
	err := ctx.DB.First(client, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 返回 nil 表示未找到记录
		}
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	return client, nil
}

func (s *ClientService) UpdateClient(ctx *app.Context, client *model.Client) error {
	client.Updated = time.Now()

	err := ctx.DB.Save(client).Error
	if err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}

	return nil
}

func (s *ClientService) DeleteClient(ctx *app.Context, client *model.Client) error {
	err := ctx.DB.Delete(client).Error
	if err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}

	return nil
}

func (s *ClientService) FindClientByCriteria(ctx *app.Context, vmuuid, sn, hostname, ip string) (*model.Client, error) {
	client := &model.Client{}
	query := ctx.DB.Model(&model.Client{})

	if vmuuid != "" {
		query = query.Where("vm_uuid = ?", vmuuid)
	}
	if sn != "" {
		query = query.Where("sn = ?", sn)
	}
	if hostname != "" {
		query = query.Where("hostname = ?", hostname)
	}
	if ip != "" {
		query = query.Where("ip = ?", ip)
	}

	err := query.First(client).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("failed to find client: %w", err)
	}

	return client, nil
}
