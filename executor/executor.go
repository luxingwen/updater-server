package executor

import (
	"context"
	"encoding/json"
	"runtime/debug"
	"time"
	"updater-server/pkg/app"
	"updater-server/service"
	"updater-server/wsserver"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type ExecutorServer struct {
	WsContext                  *wsserver.Context
	TaskService                *service.TaskService
	TaskExecutionRecordService *service.TaskExecutionRecordService
	ClientService              *service.ClientService
}

const (
	TaskQueueKey = "task_queue"
)

type TaskExecItem struct {
	TaskID   string `json:"task_id"`   // 任务ID
	Category string `json:"category"`  // 任务分类 task/record
	TaskType string `json:"task_type"` // 任务类型 root/batches
	TraceId  string `json:"trace_id"`  // 跟踪id
}

func EnqueueTask(ctx *app.Context, task TaskExecItem) error {
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return ctx.Redis.Enqueue(context.Background(), TaskQueueKey, string(taskJSON))
}

func (es *ExecutorServer) Execute(ctx *app.Context) error {
	taskJSON, err := ctx.Redis.Dequeue(context.Background(), TaskQueueKey)
	if err != nil {
		//ctx.Logger.Error("dequeue task error:", err)
		return err
	}

	task := TaskExecItem{}
	err = json.Unmarshal([]byte(taskJSON), &task)
	if err != nil {
		ctx.Logger.Error("unmarshal task error:", err)
		return err
	}

	ctx.Logger = ctx.Logger.With(zap.String("traceID", task.TraceId))

	if task.Category == "task" {
		es.ExecuteTask(ctx, task)
	}

	// 分批任务
	if task.Category == "record" {
		es.ExecuteTaskRecord(ctx, task)
	}
	return nil
}

func (es *ExecutorServer) Worker(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			es.WsContext.Logger.Errorf("Recovered from panic in Worker: %v\n%s", r, debug.Stack())
		}
	}()

	for {
		select {
		case <-ctx.Done():
			// The context has been cancelled, stop the worker
			return
		default:
			// Continue to the next task
		}

		err := es.Execute(es.WsContext.AppContext())
		if err != nil {
			if err == redis.Nil {
				// No task in the queue, sleep for a while and try again
				time.Sleep(time.Second)
				continue
			}
			// Log the error and continue to the next task
			es.WsContext.Logger.Errorf("Failed to execute task: %v", err)
			continue
		}

	}
}
