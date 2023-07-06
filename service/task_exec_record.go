package service

import (
	"encoding/json"
	"updater-server/model"
	"updater-server/pkg/app"

	"gorm.io/gorm"
)

type TaskExecutionRecordService struct{}

func NewTaskExecutionRecordService() *TaskExecutionRecordService {
	return &TaskExecutionRecordService{}
}

func (ts *TaskExecutionRecordService) CreateRecord(ctx *app.Context, record *model.TaskExecutionRecord) error {
	result := ctx.DB.Create(record)
	return result.Error
}

func (ts *TaskExecutionRecordService) UpdateRecord(ctx *app.Context, updatedRecord *model.TaskExecutionRecord) error {
	result := ctx.DB.Model(&model.TaskExecutionRecord{}).Where("record_id = ?", updatedRecord.RecordID).Updates(updatedRecord)
	return result.Error
}

// 获取一条记录，根据recordID
func (ts *TaskExecutionRecordService) GetRecordInfo(ctx *app.Context, recordID string) (*model.TaskExecutionRecord, error) {
	var record model.TaskExecutionRecord
	result := ctx.DB.Where("record_id = ?", recordID).First(&record)
	return &record, result.Error
}

// 更新任务状态
func (ts *TaskExecutionRecordService) UpdateRecordStatus(ctx *app.Context, recordID string, status string) error {
	result := ctx.DB.Model(&model.TaskExecutionRecord{}).Where("record_id = ?", recordID).Update("status", status)
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

// 根据map更新记录
func (ts *TaskExecutionRecordService) UpdateRecordByMap(ctx *app.Context, recordID string, data map[string]interface{}) error {
	result := ctx.DB.Model(&model.TaskExecutionRecord{}).Where("record_id = ?", recordID).Updates(data)
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

	for _, item := range records {
		if item.Status == "running" {
			isDone, err := ts.CheckTaskStatus(ctx, item.RecordID)
			if err != nil {
				return nil, err
			}
			if isDone {
				item.Status = "completed"
				ts.UpdateRecordStatus(ctx, item.RecordID, "completed")
			}
		}
	}

	return &model.PagedResponse{
		Total: total,
		Data:  records,
	}, nil
}

// 检查任务状态，检查任务是否完成
func (ts *TaskExecutionRecordService) CheckTaskStatus(ctx *app.Context, taskID string) (bool, error) {
	recordInfo, err := ts.GetRecordInfo(ctx, taskID)
	if err != nil {
		return false, err
	}
	if recordInfo.Status == "completed" || recordInfo.Status == "failed" || recordInfo.Status == "success" {
		return true, nil
	}

	taskContent := &model.TaskContentInDB{}

	err = json.Unmarshal([]byte(recordInfo.Content), taskContent)
	if err != nil {
		ctx.Logger.Error("unmarshal task content error:", err)
		return false, err
	}

	if taskContent.Type == "record" {
		tcontent := make([]model.TaskContentInfo, 0)

		err = json.Unmarshal([]byte(taskContent.Content), &tcontent)
		if err != nil {
			ctx.Logger.Error("unmarshal task content error:", err)
			return false, err
		}
		for _, item := range tcontent {
			isDone, err := ts.CheckTaskStatus(ctx, item.TaskRecordId)
			if err != nil {
				return false, err
			}
			if !isDone {
				return false, nil
			}
		}

		// 跟新状态

		err = ts.UpdateRecordStatus(ctx, taskID, "completed")
		return true, nil
	}
	return false, nil
}
