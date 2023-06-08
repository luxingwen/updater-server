package service

import (
	"encoding/json"
	"updater-server/model"
	"updater-server/pkg/app"

	"gorm.io/gorm"
)

type TaskExecutionRecordService struct{}

func (ts *TaskExecutionRecordService) CreateRecord(ctx *app.Context, record *model.TaskExecutionRecord) error {
	result := ctx.DB.Create(record)
	return result.Error
}

func (ts *TaskExecutionRecordService) UpdateRecord(ctx *app.Context, updatedRecord *model.TaskExecutionRecord) error {
	result := ctx.DB.Model(&model.TaskExecutionRecord{}).Where("record_id = ?", updatedRecord.RecordID).Updates(updatedRecord)
	return result.Error
}

func (ts *TaskExecutionRecordService) UpdaterRecordContent(ctx *app.Context, recordID string, content interface{}) error {
	b, err := json.Marshal(content)
	if err != nil {
		return err
	}

	result := ctx.DB.Model(&model.TaskExecutionRecord{}).Where("record_id = ?", recordID).Update("content", string(b))
	return result.Error
}

func (ts *TaskExecutionRecordService) DeleteRecord(ctx *app.Context, recordID string) error {
	var record model.TaskExecutionRecord
	result := ctx.DB.Where("record_id = ?", recordID).Delete(&record)
	return result.Error
}

func (ts *TaskExecutionRecordService) GetAllTaskExecRecords(ctx *app.Context, query *model.ReqTaskRecordeQuery) (*model.PagedResponse, error) {
	sess := ctx.DB.Session(&gorm.Session{})

	if query.TaskId != "" {
		sess = sess.Where("task_id = ?", query.TaskId)
	}

	if len(query.RecordIds) > 0 {
		sess = sess.Where("record_id in (?)", query.RecordIds)
	}

	var total int64
	result := sess.Model(&model.TaskExecutionRecord{}).Count(&total)
	if result.Error != nil {
		return nil, result.Error
	}

	var records []model.TaskExecutionRecord
	result = sess.Offset(query.GetOffset()).Limit(query.PageSize).Find(&records)
	if result.Error != nil {
		return nil, result.Error
	}

	return &model.PagedResponse{
		Total: total,
		Data:  records,
	}, nil
}
