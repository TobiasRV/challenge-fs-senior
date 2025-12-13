package handlers

import (
	"database/sql"
	"errors"
	"time"

	"github.com/TobiasRV/challenge-fs-senior/internals/interfaces"
	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	"github.com/TobiasRV/challenge-fs-senior/internals/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task (Manager only)
// @Tags Tasks
// @Accept json
// @Produce json
// @Param request body interfaces.CreateTaksPayload true "Task creation data"
// @Success 201 {object} interfaces.GetTasksResponse
// @Failure 400 {object} utils.ErrorResponse "Validation error"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Manager only"
// @Failure 404 {object} utils.ErrorResponse "Project not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /tasks [post]
func (h *Handler) CreateTask(c *fiber.Ctx) error {
	userRole := c.Locals("userRole")
	if userRole != "Manager" {
		return c.Status(fiber.StatusForbidden).JSON(utils.ErrorString("unauthorized"))
	}

	payload := interfaces.CreateTaksPayload{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	err := h.validator.Validate(payload)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewValidatorError(err))
	}

	projectUUID, err := uuid.Parse(payload.ProjectId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	var userUUID uuid.UUID

	if payload.UserId != "" {
		userUUID, err = uuid.Parse(payload.UserId)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusNotFound).JSON(utils.ErrorString("project not found"))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	task, err := h.taskRepository.CreateTask(c.Context(), database.CreateTasksParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		ProjectID: projectUUID,
		Title:     payload.Title,
		Description: sql.NullString{
			String: payload.Description,
			Valid:  payload.Description != "",
		},
		UserID: uuid.NullUUID{
			UUID:  userUUID,
			Valid: userUUID != uuid.Nil,
		},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

// GetTasks godoc
// @Summary Get all tasks
// @Description Get paginated list of tasks. Members can only see their own tasks.
// @Tags Tasks
// @Accept json
// @Produce json
// @Param projectId query string false "Filter by project ID (required for Admin/Manager)"
// @Param title query string false "Filter by title"
// @Param cursor query string false "Pagination cursor"
// @Param limit query int true "Number of items per page"
// @Success 200 {object} interfaces.TasksListResponse
// @Failure 400 {object} utils.ErrorResponse "Bad request - projectId required for Admin/Manager"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /tasks [get]
func (h *Handler) GetTasks(c *fiber.Ctx) error {

	queryParams := interfaces.GetTasksParams{}

	if err := c.QueryParser(&queryParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewError(err))
	}

	userRole := c.Locals("userRole")

	if userRole != "Member" && queryParams.ProjectId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorString("projectId required"))
	}

	var err error
	var projectUUID uuid.UUID

	if queryParams.ProjectId != "" {
		projectUUID, err = uuid.Parse(queryParams.ProjectId)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
		}
	}

	var userUUID uuid.UUID

	if userRole == "Member" {
		userId := c.Locals("userId")
		userUUID, err = uuid.Parse(userId.(string))

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
		}
	}

	cursor := queryParams.Cursor
	isFirstPage := true
	var cursorCreatedAt time.Time
	var cursorId uuid.UUID
	pointsNext := false
	if cursor != "" {
		decodedCursor, err := utils.DecodeCursor(cursor)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
		}
		cursorCreatedAt, err = time.Parse(time.RFC3339Nano, decodedCursor["created_at"].(string))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
		}
		cursorId, err = uuid.Parse(decodedCursor["id"].(string))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
		}

		pointsNext = decodedCursor["points_next"] == true

		isFirstPage = false
	}

	tasks, err := h.taskRepository.GetTasks(c.Context(), interfaces.GetTasksFilters{
		Limit:           queryParams.Limit,
		IsFirstPage:     isFirstPage,
		PointsNext:      pointsNext,
		CursorCreatedAt: cursorCreatedAt,
		CursorId:        cursorId,
		Title:           queryParams.Title,
		ProjectId:       projectUUID,
		UserId:          userUUID,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	hasPagination := len(tasks) > int(queryParams.Limit)

	if hasPagination {
		if cursor == "" || pointsNext {
			tasks = tasks[:int(queryParams.Limit)]
		} else {
			tasks = tasks[len(tasks)-int(queryParams.Limit):]
		}
	}

	if len(tasks) == 0 {
		pager := utils.GeneratePager(utils.Cursor{}, utils.Cursor{})
		return c.Status(200).JSON(fiber.Map{
			"data":       tasks,
			"pagination": pager,
		})
	}

	var nextCursor utils.Cursor
	var prevCursor utils.Cursor

	if cursor == "" {
		if hasPagination {
			nextCursor = utils.CreateCursor(tasks[len(tasks)-1].ID, tasks[len(tasks)-1].CreatedAt, true)
		}
	} else {
		if pointsNext {

			if hasPagination {
				nextCursor = utils.CreateCursor(tasks[len(tasks)-1].ID, tasks[len(tasks)-1].CreatedAt, true)
			}

			prevCursor = utils.CreateCursor(tasks[0].ID, tasks[0].CreatedAt, false)
		} else {
			nextCursor = utils.CreateCursor(tasks[len(tasks)-1].ID, tasks[len(tasks)-1].CreatedAt, true)

			if hasPagination {
				prevCursor = utils.CreateCursor(tasks[0].ID, tasks[0].CreatedAt, false)
			}
		}
	}

	pager := utils.GeneratePager(nextCursor, prevCursor)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":       tasks,
		"pagination": pager,
	})
}

// UpdateTask godoc
// @Summary Update a task
// @Description Update task information (Manager only)
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param request body interfaces.UpdateTaskPayload true "Task update data"
// @Success 200 {object} interfaces.GetTasksResponse
// @Failure 400 {object} utils.ErrorResponse "Validation error"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Manager only"
// @Failure 404 {object} utils.ErrorResponse "Task not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /tasks/{id} [put]
func (h *Handler) UpdateTask(c *fiber.Ctx) error {

	userRole := c.Locals("userRole")

	if userRole != "Manager" {
		return c.Status(fiber.StatusForbidden).JSON(utils.ErrorString("unauthorized"))
	}

	taskId := c.Params("id")

	if taskId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorString("invalid id"))
	}

	taskUUID, err := uuid.Parse(taskId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	_, err = h.taskRepository.GetTaskById(c.Context(), taskUUID)

	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusNotFound).JSON(utils.ErrorString("project not found"))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	payload := interfaces.UpdateTaskPayload{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	err = h.validator.Validate(payload)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewValidatorError(err))
	}

	updatedTask, err := h.taskRepository.UpdateTask(c.Context(), interfaces.UpdateTaskData{
		Title: payload.Title,
		Description: sql.NullString{
			String: payload.Description,
			Valid:  payload.Description != "",
		},
		UserId: uuid.NullUUID{
			UUID:  payload.UserId,
			Valid: payload.UserId != uuid.Nil,
		},
		Status:    payload.Status,
		UpdatedAt: time.Now().UTC(),
		ID:        taskUUID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	return c.Status(fiber.StatusOK).JSON(updatedTask)
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task by ID (Manager only)
// @Tags Tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} interfaces.DeleteTaskResponse
// @Failure 400 {object} utils.ErrorResponse "Invalid ID"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Manager only"
// @Failure 404 {object} utils.ErrorResponse "Task not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /tasks/{id} [delete]
func (h *Handler) DeleteTask(c *fiber.Ctx) error {
	userRole := c.Locals("userRole")

	if userRole != "Manager" {
		return c.Status(fiber.StatusForbidden).JSON(utils.ErrorString("unauthorized"))
	}

	taskId := c.Params("id")

	if taskId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorString("invalid id"))
	}

	taskUUID, err := uuid.Parse(taskId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	err = h.taskRepository.DeleteTask(c.Context(), taskUUID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"deleted": true,
	})
}
