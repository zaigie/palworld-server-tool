package task

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/robfig/cron/v3"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/internal/logger"
	"github.com/zaigie/palworld-server-tool/internal/tool"
	"github.com/zaigie/palworld-server-tool/service"
	"go.etcd.io/bbolt"
)

const rconTaskTagPrefix = "rcon-task:"

var rconExecutionMu sync.Mutex

func ValidateCronExpression(expression string) error {
	expression = strings.TrimSpace(expression)
	if expression == "" {
		return errors.New("cron expression is required")
	}
	if _, err := cron.ParseStandard(expression); err != nil {
		return fmt.Errorf("invalid cron expression: %w", err)
	}
	return nil
}

func RegisterRconTask(db *bbolt.DB, rconTask database.RconTask) error {
	UnregisterRconTask(rconTask.UUID)
	if !rconTask.Enabled {
		return nil
	}
	if err := ValidateCronExpression(rconTask.Cron); err != nil {
		return err
	}
	scheduler := getScheduler()
	_, err := scheduler.NewJob(
		gocron.CronJob(rconTask.Cron, false),
		gocron.NewTask(ExecuteRconTask, db, rconTask.UUID),
		gocron.WithName(rconTask.Name),
		gocron.WithTags(rconTaskTagPrefix+rconTask.UUID),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	return err
}

func UnregisterRconTask(taskUUID string) {
	getScheduler().RemoveByTags(rconTaskTagPrefix + taskUUID)
}

func LoadRconTasks(db *bbolt.DB) error {
	tasks, err := service.ListRconTasks(db)
	if err != nil {
		return err
	}
	for _, rconTask := range tasks {
		if err := RegisterRconTask(db, rconTask); err != nil {
			logger.Errorf("Failed to schedule RCON task %s: %v\n", rconTask.UUID, err)
		}
	}
	return nil
}

func NextRconTaskRun(taskUUID string) *time.Time {
	tag := rconTaskTagPrefix + taskUUID
	for _, job := range getScheduler().Jobs() {
		for _, jobTag := range job.Tags() {
			if jobTag != tag {
				continue
			}
			next, err := job.NextRun()
			if err != nil || next.IsZero() {
				return nil
			}
			return &next
		}
	}
	return nil
}

func ExecuteRconTask(db *bbolt.DB, taskUUID string) error {
	return executeRconTask(db, taskUUID, tool.CustomCommand)
}

func executeRconTask(db *bbolt.DB, taskUUID string, execute func(string) (string, error)) error {
	rconExecutionMu.Lock()
	defer rconExecutionMu.Unlock()

	rconTask, err := service.GetRconTask(db, taskUUID)
	if err != nil {
		return err
	}
	rconCommand, err := service.GetRconCommand(db, rconTask.RconUUID)
	if err != nil {
		ranAt := time.Now()
		_ = service.UpdateRconTaskExecution(db, taskUUID, "failed", "", err.Error(), ranAt)
		return err
	}
	execCommand := strings.TrimSpace(strings.Join([]string{rconCommand.Command, rconTask.Content}, " "))
	result, executeErr := execute(execCommand)
	ranAt := time.Now()
	if executeErr != nil {
		if updateErr := service.UpdateRconTaskExecution(db, taskUUID, "failed", result, executeErr.Error(), ranAt); updateErr != nil {
			logger.Errorf("Failed to save RCON task result %s: %v\n", taskUUID, updateErr)
		}
		return executeErr
	}
	if err := service.UpdateRconTaskExecution(db, taskUUID, "success", result, "", ranAt); err != nil {
		return err
	}
	logger.Infof("Scheduled RCON task %s executed successfully\n", taskUUID)
	return nil
}
