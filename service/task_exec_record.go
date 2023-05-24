package service

import (
	"updater-server/model"
	"updater-server/pkg/app"
)

type TaskExecutionRecordService struct {}

func (ts *TaskExecutionRecordService) CreateRecord(ctx *app.Context, record *model.TaskExecutionRecord) error {
	result := ctx.DB.Create(record)
	return result.Error
}

func (ts *TaskExecutionRecordService) UpdateRecord(ctx *app.Context, updatedRecord *model.TaskExecutionRecord) error {
	result := ctx.DB.Model(&model.TaskExecutionRecord{}).Where("record_id = ?", updatedRecord.RecordID).Updates(updatedRecord)
	return result.Error
}

func (ts *TaskExecutionRecordService) DeleteRecord(ctx *app.Context, recordID string) error {
	var record model.TaskExecutionRecord
	result := ctx.DB.Where("record_id = ?", recordID).Delete(&record)
	return result.Error
}
