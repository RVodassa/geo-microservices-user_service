package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/RVodassa/geo-microservices-user/internal/domain/entity"
)

type UserRepositoryProvider interface {
	Register(ctx context.Context, u *entity.User) (id uint64, err error)
	Delete(ctx context.Context, id uint64) error
	Profile(ctx context.Context, id uint64) (*entity.User, error)
	List(ctx context.Context, offset, limit uint64) ([]*entity.User, uint64, error)
	Login(ctx context.Context, login, password string) (bool, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Register(ctx context.Context, u *entity.User) (id uint64, err error) {
	query := `INSERT INTO users (password, login) VALUES ($1, $2) RETURNING id;`

	err = r.db.QueryRowContext(ctx, query, u.Password, u.Login).Scan(&u.ID)
	if err != nil {
		return 0, err
	}

	// Возвращаем ID добавленного пользователя
	return u.ID, nil
}

func (r *UserRepository) Delete(ctx context.Context, id uint64) error {

	query := `DELETE FROM users WHERE id = ($1)`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) Profile(ctx context.Context, id uint64) (*entity.User, error) {
	var u entity.User
	query := `SELECT id, login FROM users WHERE id = $1;`
	row := r.db.QueryRowContext(ctx, query, id)
	if err := row.Scan(&u.ID, &u.Login); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) List(ctx context.Context, offset, limit uint64) ([]*entity.User, uint64, error) {
	var users []*entity.User
	var totalCount uint64

	// Получаем общее количество записей
	countQuery := `SELECT COUNT(*) FROM users;`
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total user count: %w", err)
	}

	// Получаем пользователей с использованием offset и limit
	query := `SELECT id, login FROM users ORDER BY id LIMIT $1 OFFSET $2;`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute user list query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Login); err != nil {
			return nil, 0, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, &u)
	}

	// Проверяем на ошибки итерации
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("row iteration error: %w", err)
	}

	return users, totalCount, nil
}

func (r *UserRepository) Login(ctx context.Context, login, password string) (bool, error) {
	user := &entity.User{}
	// Запрос на выборку только id и password (без других данных)
	query := `SELECT id, password FROM users WHERE login = $1;`
	row := r.db.QueryRowContext(ctx, query, login)

	// Проверяем на ошибки при сканировании строки
	if err := row.Scan(&user.ID, &user.Password); err != nil {
		// Если строка не найдена, возвращаем false, а не ошибку
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	// Если пароль совпадает, возвращаем true
	if user.Password == password {
		return true, nil
	}

	// Если пароль не совпадает, возвращаем false
	return false, nil
}
