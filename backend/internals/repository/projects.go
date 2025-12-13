package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/TobiasRV/challenge-fs-senior/internals/interfaces"
	"github.com/TobiasRV/challenge-fs-senior/internals/models"
	"github.com/TobiasRV/challenge-fs-senior/internals/sqlc/database"
	"github.com/google/uuid"
)

type ProjectRepository struct {
	queries *database.Queries
	db      *sql.DB
}

func NewProjectRepository(queries *database.Queries, db *sql.DB) *ProjectRepository {
	return &ProjectRepository{
		queries: queries,
		db:      db,
	}
}

func (pr *ProjectRepository) CreateProject(c context.Context, data database.CreateProjectParams) (models.Project, error) {
	newProject, err := pr.queries.CreateProject(c, database.CreateProjectParams{
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		Name:      data.Name,
		TeamID:    data.TeamID,
		ManagerID: data.ManagerID,
	})

	if err != nil {
		return models.Project{}, err
	}

	return models.DatabaseProjectToProject(newProject), nil
}

func (pr *ProjectRepository) GetProjects(c context.Context, filters interfaces.GetProjectsFilters) ([]interfaces.GetProjectsResponse, error) {

	var sql sq.SelectBuilder

	if filters.WithStats {
		sql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select(
			"p.id",
			"p.created_at",
			"p.updated_at",
			"p.name",
			"p.team_id",
			"p.manager_id",
			"p.status",

			"COUNT(t.id) FILTER (WHERE t.status = 'ToDo') AS \"ToDoTasks\"",
			"COUNT(t.id) FILTER (WHERE t.status = 'InProgress') AS \"InProgressTasks\"",
			"COUNT(t.id) FILTER (WHERE t.status = 'Done') AS \"DoneTasks\"",
		).From("projects p").LeftJoin("tasks t ON t.project_id = p.id")
	} else {
		sql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("*").From("projects p")
	}

	if filters.TeamId != uuid.Nil {
		sql = sql.Where(sq.Eq{"p.team_id": filters.TeamId})
	}

	if filters.ManagerId != uuid.Nil {
		sql = sql.Where(sq.Eq{"p.manager_id": filters.ManagerId})
	}

	orderAsc := true

	// Handle cursor pagination
	if !filters.IsFirstPage {
		if filters.PointsNext {
			sql = sql.Where(sq.Or{
				sq.Gt{"p.created_at": filters.CursorCreatedAt},
				sq.And{
					sq.Eq{"p.created_at": filters.CursorCreatedAt},
					sq.Gt{"p.id": filters.CursorId},
				},
			})
			orderAsc = true
		} else {
			sql = sql.Where(sq.Or{
				sq.Lt{"p.created_at": filters.CursorCreatedAt},
				sq.And{
					sq.Eq{"p.created_at": filters.CursorCreatedAt},
					sq.Lt{"p.id": filters.CursorId},
				},
			})
			orderAsc = false
		}
	}

	if filters.Name != "" {
		nameLower := strings.ToLower(filters.Name)
		sql = sql.Where(sq.Like{"LOWER(p.name)": fmt.Sprintf("%%%v%%", nameLower)})
	}

	if orderAsc {
		sql = sql.OrderBy("p.created_at ASC, p.id ASC").Limit(filters.Limit + 1)
	} else {
		sql = sql.OrderBy("p.created_at DESC, p.id DESC").Limit(filters.Limit + 1)
	}

	queryString, arg, err := sql.ToSql()

	if err != nil {
		return []interfaces.GetProjectsResponse{}, err
	}

	rows, err := pr.db.Query(queryString, arg...)

	if err != nil {
		return []interfaces.GetProjectsResponse{}, err
	}

	defer rows.Close()

	var projects []interfaces.GetProjectsResponse

	if filters.WithStats {
		for rows.Next() {
			var project interfaces.GetProjectsResponse

			if err := rows.Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt, &project.Name, &project.TeamID, &project.ManagerID, &project.Status, &project.ToDoTasks, &project.InProgressTasks, &project.DoneTasks); err != nil {
				return projects, err
			}

			projects = append(projects, project)
		}
	} else {
		for rows.Next() {
			var project interfaces.GetProjectsResponse

			if err := rows.Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt, &project.Name, &project.TeamID, &project.ManagerID, &project.Status); err != nil {
				return projects, err
			}

			projects = append(projects, project)
		}
	}

	if err = rows.Err(); err != nil {
		return []interfaces.GetProjectsResponse{}, err
	}

	if !orderAsc && len(projects) > 0 {
		for i, j := 0, len(projects)-1; i < j; i, j = i+1, j-1 {
			projects[i], projects[j] = projects[j], projects[i]
		}
	}

	return projects, nil

}

func (pr *ProjectRepository) GetProjectById(c context.Context, id uuid.UUID) (models.Project, error) {

	project, err := pr.queries.GetProjectById(c, id)

	if err != nil {
		return models.Project{}, err
	}

	return models.DatabaseProjectToProject(project), nil

}

func (pr *ProjectRepository) UpdateProject(c context.Context, data interfaces.UpdateProjectData) (models.Project, error) {
	user, err := pr.queries.UpdateProject(c, database.UpdateProjectParams{
		Name:      data.Name,
		Status:    database.Projectstatus(data.Status),
		UpdatedAt: data.UpdatedAt,
		ID:        data.ID,
	})

	if err != nil {
		return models.Project{}, err
	}

	return models.DatabaseProjectToProject(user), nil
}

func (pr *ProjectRepository) DeleteProject(c context.Context, id uuid.UUID) error {
	err := pr.queries.DeleteProject(c, id)

	return err
}
