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

// 获取所有主机
func (s *ClientService) GetAllClient(ctx *app.Context) ([]model.Client, error) {
	var clients []model.Client
	err := ctx.DB.Model(&model.Client{}).Find(&clients).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all client: %w", err)
	}

	return clients, nil
}

func (s *ClientService) GetAllClientUuid(ctx *app.Context) ([]string, error) {
	var clients []model.Client
	err := ctx.DB.Model(&model.Client{}).Find(&clients).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all client: %w", err)
	}
	var uuids []string
	for _, client := range clients {
		uuids = append(uuids, client.Uuid)
	}

	return uuids, nil
}

func (s *ClientService) GetClientByHostInfo(ctx *app.Context, hostinfo model.HostInfo) (r []string, err error) {
	if hostinfo.All {
		return s.GetAllClientUuid(ctx)
	}
	return hostinfo.Clients, nil
}

func (s *ClientService) GetClientByUUID(ctx *app.Context, uuid string) (*model.Client, error) {
	client := &model.Client{}
	err := ctx.DB.Model(&model.Client{}).Where("uuid = ?", uuid).First(client).Error
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

// 更新客户端状态
func (s *ClientService) UpdateClientStatus(ctx *app.Context, clientUuid string, status string) error {
	err := ctx.DB.Model(&model.Client{}).Where("uuid = ?", clientUuid).Update("status", status).Error
	return err
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

func (s *ClientService) Register(ctx *app.Context, client *model.Client) error {

	// 检查是否已经注册过
	r, err := s.GetClientByUUID(ctx, client.Uuid)
	if err != nil {
		return fmt.Errorf("failed to register client: %w", err)
	}
	if err == nil && r == nil {
		client.Created = time.Now()
		client.Updated = time.Now()

		err := ctx.DB.Create(client).Error
		if err != nil {
			return fmt.Errorf("failed to register client: %w", err)
		}
		return nil
	}

	// 更新注册信息
	r.VMUUID = client.VMUUID
	r.Hostname = client.Hostname
	r.IP = client.IP
	r.SN = client.SN
	r.Updated = time.Now()
	r.Status = "online"
	r.ProxyID = client.ProxyID
	ctx.Logger.Infof("client already registered, updating..., client: %v", r)

	err = ctx.DB.Save(r).Error
	return err

}

func (s *ClientService) QueryClient(ctx *app.Context, query *model.ReqClientQuery) (r *model.PagedResponse, err error) {
	var clients []model.Client
	var total int64

	// Create a new session
	sess := ctx.DB.Session(&gorm.Session{})

	if query.Uuid != "" {
		sess = sess.Where("uuid = ?", query.Uuid)
	}
	if query.Vmuuid != "" {
		sess = sess.Where("vmuuid = ?", query.Vmuuid)
	}
	if query.Hostname != "" {
		sess = sess.Where("hostname = ?", query.Hostname)
	}
	if query.Ip != "" {
		sess = sess.Where("ip = ?", query.Ip)
	}
	if query.Sn != "" {
		sess = sess.Where("sn = ?", query.Sn)
	}

	result := sess.Limit(query.PageSize).Offset(query.GetOffset()).Find(&clients)
	if result.Error != nil {
		return nil, result.Error
	}

	result = sess.Model(&model.Client{}).Count(&total)
	if result.Error != nil {
		return nil, result.Error
	}

	response := &model.PagedResponse{
		Data:     clients,
		Current:  query.Current,
		PageSize: query.PageSize,
		Total:    total,
	}

	return response, nil
}
