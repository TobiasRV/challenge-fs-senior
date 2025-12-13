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

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new project (Manager only)
// @Tags Projects
// @Accept json
// @Produce json
// @Param request body interfaces.CreateProjectPayload true "Project creation data"
// @Success 201 {object} interfaces.GetProjectsResponse
// @Failure 400 {object} utils.ErrorResponse "Validation error or manager not found"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Manager only"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /projects [post]
func (h *Handler) CreateProject(c *fiber.Ctx) error {
	userRole := c.Locals("userRole")
	userId := c.Locals("userId")

	if userRole != "Manager" {
		return c.Status(fiber.StatusForbidden).JSON(utils.ErrorString("unauthorized"))
	}

	userUUID, err := uuid.Parse(userId.(string))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	payload := interfaces.CreateProjectPayload{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	err = h.validator.Validate(payload)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewValidatorError(err))
	}

	user, err := h.userRepository.GetUserById(c.Context(), userUUID)

	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorString("manager doesn't exists"))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	project, err := h.projectRepository.CreateProject(c.Context(), database.CreateProjectParams{
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      payload.Name,
		TeamID:    user.TeamId,
		ManagerID: user.ID,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	return c.Status(fiber.StatusCreated).JSON(project)

}

// GetProjects godoc
// @Summary Get all projects
// @Description Get paginated list of projects (Admin or Manager only)
// @Tags Projects
// @Accept json
// @Produce json
// @Param name query string false "Filter by name"
// @Param teamId query string false "Filter by team ID"
// @Param managerId query string false "Filter by manager ID"
// @Param withStats query bool false "Include task statistics"
// @Param cursor query string false "Pagination cursor"
// @Param limit query int true "Number of items per page"
// @Success 200 {object} interfaces.ProjectsListResponse
// @Failure 400 {object} utils.ErrorResponse "Bad request - teamId or managerId required"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Admin or Manager only"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /projects [get]
func (h *Handler) GetProjects(c *fiber.Ctx) error {

	userRole := c.Locals("userRole")

	if userRole != "Admin" && userRole != "Manager" {
		return c.Status(fiber.StatusForbidden).JSON(utils.ErrorString("unauthorized"))
	}

	queryParams := interfaces.GetProjectsParams{}

	if err := c.QueryParser(&queryParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewError(err))
	}

	if queryParams.ManagerId == uuid.Nil && queryParams.TeamId == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorString("teamId or managerId required"))
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

	projects, err := h.projectRepository.GetProjects(c.Context(), interfaces.GetProjectsFilters{
		Name:            queryParams.Name,
		TeamId:          queryParams.TeamId,
		ManagerId:       queryParams.ManagerId,
		Limit:           queryParams.Limit,
		IsFirstPage:     isFirstPage,
		PointsNext:      pointsNext,
		CursorCreatedAt: cursorCreatedAt,
		CursorId:        cursorId,
		WithStats:       queryParams.WithStats,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	hasPagination := len(projects) > int(queryParams.Limit)

	if hasPagination {
		if cursor == "" || pointsNext {
			projects = projects[:int(queryParams.Limit)]
		} else {
			projects = projects[len(projects)-int(queryParams.Limit):]
		}
	}

	if len(projects) == 0 {
		pager := utils.GeneratePager(utils.Cursor{}, utils.Cursor{})
		return c.Status(200).JSON(fiber.Map{
			"data":       projects,
			"pagination": pager,
		})
	}

	var nextCursor utils.Cursor
	var prevCursor utils.Cursor

	if cursor == "" {
		if hasPagination {
			nextCursor = utils.CreateCursor(projects[len(projects)-1].ID, projects[len(projects)-1].CreatedAt, true)
		}
	} else {
		if pointsNext {

			if hasPagination {
				nextCursor = utils.CreateCursor(projects[len(projects)-1].ID, projects[len(projects)-1].CreatedAt, true)
			}

			prevCursor = utils.CreateCursor(projects[0].ID, projects[0].CreatedAt, false)
		} else {
			nextCursor = utils.CreateCursor(projects[len(projects)-1].ID, projects[len(projects)-1].CreatedAt, true)

			if hasPagination {
				prevCursor = utils.CreateCursor(projects[0].ID, projects[0].CreatedAt, false)
			}
		}
	}

	pager := utils.GeneratePager(nextCursor, prevCursor)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":       projects,
		"pagination": pager,
	})
}

// UpdateProject godoc
// @Summary Update a project
// @Description Update project information (Manager only, must be the project manager)
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param request body interfaces.UpdateProjectPayload true "Project update data"
// @Success 200 {object} interfaces.GetProjectsResponse
// @Failure 400 {object} utils.ErrorResponse "Validation error"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Must be the project manager"
// @Failure 404 {object} utils.ErrorResponse "Project not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /projects/{id} [put]
func (h *Handler) UpdateProject(c *fiber.Ctx) error {

	userRole := c.Locals("userRole")
	userId := c.Locals("userId")

	userUUID, err := uuid.Parse(userId.(string))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	if userRole != "Manager" {
		return c.Status(fiber.StatusForbidden).JSON(utils.ErrorString("unauthorized"))
	}

	projectId := c.Params("id")

	if projectId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorString("invalid id"))
	}

	projectUUID, err := uuid.Parse(projectId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	existingProject, err := h.projectRepository.GetProjectById(c.Context(), projectUUID)

	if errors.Is(err, sql.ErrNoRows) {
		return c.Status(fiber.StatusNotFound).JSON(utils.ErrorString("project not found"))
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	if existingProject.ManagerID != userUUID {
		return c.Status(fiber.StatusForbidden).JSON(utils.ErrorString("unauthorized"))
	}

	payload := interfaces.UpdateProjectPayload{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	err = h.validator.Validate(payload)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(utils.NewValidatorError(err))
	}

	updatedProject, err := h.projectRepository.UpdateProject(c.Context(), interfaces.UpdateProjectData{
		Name:      payload.Name,
		Status:    payload.Status,
		UpdatedAt: time.Now().UTC(),
		ID:        projectUUID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	return c.Status(fiber.StatusOK).JSON(updatedProject)
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Delete a project by ID (Manager only)
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} interfaces.MessageResponse
// @Failure 400 {object} utils.ErrorResponse "Invalid ID"
// @Failure 403 {object} utils.ErrorResponse "Forbidden - Manager only"
// @Failure 404 {object} utils.ErrorResponse "Project not found"
// @Failure 500 {object} utils.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /projects/{id} [delete]
func (h *Handler) DeleteProject(c *fiber.Ctx) error {
	userRole := c.Locals("userRole")

	if userRole != "Manager" {
		return c.Status(fiber.StatusForbidden).JSON(utils.ErrorString("unauthorized"))
	}

	projectId := c.Params("id")

	if projectId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(utils.ErrorString("invalid id"))
	}

	projectUUID, err := uuid.Parse(projectId)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	err = h.projectRepository.DeleteProject(c.Context(), projectUUID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.NewError(err))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"deleted": true,
	})
}
