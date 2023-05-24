package model

import (
   "time"

   "github.com/google/uuid"
)

type Task struct {
	Id uint `gorm:"primaryKey"`
    TaskID       string `json:"task_id"`
    TaskName     string    `json:"task_name"`
    TaskType     string    `json:"task_type"`
    TaskStatus   string    `json:"task_status"`
    ParentTaskID string `json:"parent_task_id"`
	Content string `json:"content"`
	Description string `json:"description"`
    Creater    string    `json:"creater"`
    TeamID       string       `json:"team_id"`
	NextTaskID string `json:"next_task_id"`
	Created  time.Time `gorm:"column:created" json:"created"`
	Updated  time.Time `gorm:"column:updated" json:"updated"`
}


type TaskExecutionRecord struct {
	Id uint `gorm:"primaryKey"`
    RecordID       string         `json:"record_id"`
    TaskID         string        `json:"task_id"`
    ClientUUID     string        `json:"client_uuid"`
    TaskType       string            `json:"task_type"`
    Status         string            `json:"status"`
    StartTime      time.Time         `json:"start_time"`
    EndTime        time.Time         `json:"end_time"`
    Stdout string            `json:"stdout"`
    Stderr    string            `json:"stderr"`
    Message        string            `json:"message"`
    ScriptExitCode int               `json:"script_exit_code"`
    Code     string               `json:"code"`
    Content        string            `json:"content"`
    Timeout        time.Duration     `json:"timeout"`
    ParentRecordID string         `json:"parent_record_id"`
	NextRecordID string `json:"next_record_id"`
}




type HostInfo struct {
	All bool `json:"all"`
	Clients []string `json:"clients"`
}

type BatchTask struct {
	count int
	Style string `json:"style"`
	Number int `json:"number"`
}



type TaskBatchesInfo struct {
	TaskID string `json:"task_id"`
	Total int `json:"total"`
	NextTaskID string `json:"next_task_id"`
	Clients []string `json:"clients"`
}

func (self BatchTask)GenerateTaskBatchesInfo(clients []string)(r []TaskBatchesInfo) {
	bs:=self.GetBatchesList(len(clients))
	for i:=0;i<len(bs);i++{
		r=append(r,TaskBatchesInfo{
			TaskID:      uuid.New().String(),
			Total:       bs[i],
			NextTaskID:  "",
			Clients:     clients[:bs[i]],
		})
		clients=clients[bs[i]:]
	}

	for i:=0;i<len(r);i++{
		if i==len(r)-1{
			r[i].NextTaskID=""
		}else{
			r[i].NextTaskID=r[i+1].TaskID
		}
	}
	return
}


func (self BatchTask)GetBatchesList(total int)(r []int) {
	if self.Style=="normal" {
		return self.GetAverage(total)
	}
	for self.count<total{
		num:=self.GetNumber(total)
		if self.count+num>=total{
			num=total-self.count
		}
		r=append(r,num)
		self.count+=num
	}
	return
}



func (self BatchTask)GetNumber(total int)int{
	if self.Style=="normal" {
		if self.Number==0{
			return total
		}
		return self.Number
	}
	return total
}



func (self BatchTask)GetAverage(total int)(r []int) {
	if self.Number<=0{
		self.Number=1
	}

	if total<=self.Number{
		self.Number=total
	}
	r=make([]int,self.Number)

	for i:=0;i<self.Number;i++{
		r[i]=0
	}
	for total>0{
		for i:=0;i<self.Number && total>0;i++{
				r[i]++
				total--
		}
	}
	return
}


// 是否是批次任务
func (self BatchTask)IsBatchTask()bool{
	if self.Style=="" || self.Style=="all"{
		return false
	}
	return true
}


type ProgramActionTask struct {
	ProgramUuid string `json:"programUuid"`
	ProgramActionUuid string `json:"programActionUuid"`
}

type ReqTask struct {
	Name string `json:"name"`
	Creater string `json:"creater"`
	TeamID string `json:"team_id"`
	Description string `json:"description"`
	HostInfo HostInfo `json:"hostInfo"`
	BatchTask BatchTask `json:"batchTask"`
}


type ReqTaskProgramAction struct {
	ReqTask
	Content ProgramActionTask `json:"content"`
}