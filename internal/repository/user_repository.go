package repository

import (
	"context"
	"database/sql"
	"geo-microservices/user/internal/domain/entity"
)

type UserRepositoryProvider interface {
	Register(ctx context.Context, u *entity.User) (id uint64, err error)
	Delete(ctx context.Context, id uint64) error
	Profile(ctx context.Context, id uint64) (*entity.User, error)
	List(ctx context.Context, offset, limit uint64) ([]*entity.User, uint64, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, u *entity.User) (id uint64, err error) {
	// Убедитесь, что используете правильное имя таблицы (например, "users")
	query := `INSERT INTO users (password, login) VALUES ($1, $2) RETURNING id;`

	// Используем QueryRow, чтобы получить результат в одну строку
	err = r.db.QueryRowContext(ctx, query, u.Password, u.Login).Scan(&id)
	if err != nil {
		return 0, err
	}

	// Возвращаем ID добавленного пользователя
	return id, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uint64) error {

	query := `DELETE FROM users WHERE id = ($1)`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
