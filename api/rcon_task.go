package api

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zaigie/palworld-server-tool/internal/database"
	"github.com/zaigie/palworld-server-tool/internal/task"
	"github.com/zaigie/palworld-server-tool/service"
)

type RconTaskResponse struct {
	database.RconTask
	RconRemark string     `json:"rcon_remark"`
	NextRunAt  *time.Time `json:"next_run_at,omitempty"`
}

func buildRconTaskResponse(rconTask database.RconTask) RconTaskResponse {
	command, _ := service.GetRconCommand(database.GetDB(), rconTask.RconUUID)
	return RconTaskResponse{
		RconTask:   rconTask,
		RconRemark: command.Remark,
		NextRunAt:  task.NextRconTaskRun(rconTask.UUID),
	}
}

func validateRconTask(rconTask database.RconTask) error {
	if strings.TrimSpace(rconTask.Name) == "" {
		return errors.New("task name is required")
	}
	if strings.TrimSpace(rconTask.RconUUID) == "" {
		return errors.New("rcon_uuid is required")
	}
	if _, err := service.GetRconCommand(database.GetDB(), rconTask.RconUUID); err != nil {
		return errors.New("Rcon command not found")
	}
	return task.ValidateCronExpression(rconTask.Cron)
}

// listRconTasks godoc
//
//	@Summary		List scheduled RCON tasks
//	@Description	List persisted RCON tasks with their command remark and next run time
//	@Tags			Rcon
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{array}		RconTaskResponse
//	@Failure		400	{object}	ErrorResponse
//	@Router			/api/rcon/tasks [get]
func listRconTasks(c *gin.Context) {
	tasks, err := service.ListRconTasks(database.GetDB())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	responses := make([]RconTaskResponse, 0, len(tasks))
	for _, rconTask := range tasks {
		responses = append(responses, buildRconTaskResponse(rconTask))
	}
	c.JSON(http.StatusOK, responses)
}

// addRconTask godoc
//
//	@Summary		Add a scheduled RCON task
//	@Description	Bind a saved RCON command to a five-field cron expression
//	@Tags			Rcon
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			task	body		database.RconTask	true	"Scheduled RCON task"
//	@Success		200		{object}	RconTaskResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/api/rcon/tasks [post]
func addRconTask(c *gin.Context) {
	var rconTask database.RconTask
	if err := c.ShouldBindJSON(&rconTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validateRconTask(rconTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	rconTask.UUID = uuid.NewString()
	rconTask.CreatedAt = now
	rconTask.UpdatedAt = now
	rconTask.LastRunAt = nil
	rconTask.LastStatus = "never"
	rconTask.LastResult = ""
	rconTask.LastError = ""
	rconTask.RunCount = 0
	if err := service.AddRconTask(database.GetDB(), rconTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := task.RegisterRconTask(database.GetDB(), rconTask); err != nil {
		_ = service.DeleteRconTask(database.GetDB(), rconTask.UUID)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, buildRconTaskResponse(rconTask))
}

// putRconTask godoc
//
//	@Summary		Update a scheduled RCON task
//	@Description	Update the bound command, arguments, schedule, or enabled state
//	@Tags			Rcon
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			uuid	path		string				true	"Task UUID"
//	@Param			task	body		database.RconTask	true	"Scheduled RCON task"
//	@Success		200		{object}	RconTaskResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Router			/api/rcon/tasks/{uuid} [put]
func putRconTask(c *gin.Context) {
	existing, err := service.GetRconTask(database.GetDB(), c.Param("uuid"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rcon task not found"})
		return
	}
	var updated database.RconTask
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := validateRconTask(updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated.UUID = existing.UUID
	updated.CreatedAt = existing.CreatedAt
	updated.UpdatedAt = time.Now()
	updated.LastRunAt = existing.LastRunAt
	updated.LastStatus = existing.LastStatus
	updated.LastResult = existing.LastResult
	updated.LastError = existing.LastError
	updated.RunCount = existing.RunCount
	if err := service.PutRconTask(database.GetDB(), updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := task.RegisterRconTask(database.GetDB(), updated); err != nil {
		_ = service.PutRconTask(database.GetDB(), existing)
		_ = task.RegisterRconTask(database.GetDB(), existing)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, buildRconTaskResponse(updated))
}

// removeRconTask godoc
//
//	@Summary		Delete a scheduled RCON task
//	@Description	Unschedule and permanently delete a RCON task
//	@Tags			Rcon
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			uuid	path		string	true	"Task UUID"
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Router			/api/rcon/tasks/{uuid} [delete]
func removeRconTask(c *gin.Context) {
	taskUUID := c.Param("uuid")
	if _, err := service.GetRconTask(database.GetDB(), taskUUID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Rcon task not found"})
		return
	}
	task.UnregisterRconTask(taskUUID)
	if err := service.DeleteRconTask(database.GetDB(), taskUUID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// runRconTask godoc
//
//	@Summary		Run a scheduled RCON task now
//	@Description	Execute the task immediately and persist its result without changing its schedule
//	@Tags			Rcon
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			uuid	path		string	true	"Task UUID"
//	@Success		200		{object}	RconTaskResponse
//	@Failure		400		{object}	ErrorResponse
//	@Router			/api/rcon/tasks/{uuid}/run [post]
func runRconTask(c *gin.Context) {
	taskUUID := c.Param("uuid")
	if err := task.ExecuteRconTask(database.GetDB(), taskUUID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rconTask, err := service.GetRconTask(database.GetDB(), taskUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, buildRconTaskResponse(rconTask))
}
