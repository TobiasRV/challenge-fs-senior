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

type UserRepository struct {
	queries *database.Queries
	db      *sql.DB
}

func NewUserRepository(queries *database.Queries, db *sql.DB) *UserRepository {
	return &UserRepository{
		queries: queries,
		db:      db,
	}
}

func (ur *UserRepository) CreateUser(c context.Context, userData models.User) (models.User, error) {

	newUser, err := ur.queries.CreateUser(c, database.CreateUserParams{
		Username:  userData.Username,
		Password:  userData.Password,
		Email:     userData.Email,
		Role:      database.Userroles(userData.Role),
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
		TeamID: uuid.NullUUID{
			UUID:  userData.TeamId,
			Valid: userData.TeamId != uuid.Nil,
		},
	})

	if err != nil {
		return models.User{}, err
	}

	return models.DatabaseUserToUser(newUser), nil
}

func (ur *UserRepository) GetUserByEmail(c context.Context, email string) (models.User, error) {

	user, err := ur.queries.GetUserByEmail(c, email)

	if err != nil {
		return models.User{}, err
	}

	return models.DatabaseUserToUser(user), nil
}

func (ur *UserRepository) GetUserById(c context.Context, id uuid.UUID) (models.User, error) {

	user, err := ur.queries.GetUserById(c, id)

	if err != nil {
		return models.User{}, err
	}

	return models.DatabaseUserToUser(user), nil

}

func (ur *UserRepository) GetUsers(c context.Context, filters interfaces.GetUserFilters) ([]models.User, error) {
	sql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("*").From("users")

	if filters.TeamId != uuid.Nil {
		sql = sql.Where(sq.Eq{"team_id": filters.TeamId})
	}

	orderAsc := true

	// Handle cursor pagination
	if !filters.IsFirstPage {
		if filters.PointsNext {
			sql = sql.Where(sq.Or{
				sq.Gt{"created_at": filters.CursorCreatedAt},
				sq.And{
					sq.Eq{"created_at": filters.CursorCreatedAt},
					sq.Gt{"id": filters.CursorId},
				},
			})
			orderAsc = true
		} else {
			sql = sql.Where(sq.Or{
				sq.Lt{"created_at": filters.CursorCreatedAt},
				sq.And{
					sq.Eq{"created_at": filters.CursorCreatedAt},
					sq.Lt{"id": filters.CursorId},
				},
			})
			orderAsc = false
		}
	}

	if filters.Email != "" {
		emailLower := strings.ToLower(filters.Email)
		sql = sql.Where(sq.Like{"LOWER(email)": fmt.Sprintf("%%%v%%", emailLower)})
	}

	if filters.Role != "" {
		sql = sql.Where(sq.Eq{
			"role": filters.Role,
		})
	}

	if orderAsc {
		sql = sql.OrderBy("created_at ASC, id ASC").Limit(filters.Limit + 1)
	} else {
		sql = sql.OrderBy("created_at DESC, id DESC").Limit(filters.Limit + 1)
	}

	queryString, arg, err := sql.ToSql()

	if err != nil {
		return []models.User{}, err
	}

	rows, err := ur.db.Query(queryString, arg...)

	if err != nil {
		return []models.User{}, err
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		if err := rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Username, &user.Password, &user.Email, &user.Role, &user.TeamId); err != nil {
			return users, err
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return []models.User{}, err
	}

	if !orderAsc && len(users) > 0 {
		for i, j := 0, len(users)-1; i < j; i, j = i+1, j-1 {
			users[i], users[j] = users[j], users[i]
		}
	}

	return users, nil

}

func (ur *UserRepository) UpdateUser(c context.Context, data interfaces.UpdateUserData) (models.User, error) {
	user, err := ur.queries.UpdateUser(c, database.UpdateUserParams{
		Username:  data.Username,
		Email:     data.Email,
		UpdatedAt: data.UpdatedAt,
		ID:        data.ID,
	})

	if err != nil {
		return models.User{}, err
	}

	return models.DatabaseUserToUser(user), nil
}

func (ur *UserRepository) DeleteUser(c context.Context, id uuid.UUID) error {
	err := ur.queries.DeleteUser(c, id)

	return err
}
