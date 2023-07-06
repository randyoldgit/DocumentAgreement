package postgres

import (
	"DocumentAgreement/internal/adapters/entities"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(ctx context.Context, userAuth entities.UserAuth) error {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (username, password_hash) values ($1, $2) returning id`,
		userTable)
	row := r.db.QueryRowContext(ctx, query, userAuth.UserName, userAuth.Password)
	if err := row.Scan(&id); err != nil {
		return err
	}
	//записать в ctx ID пользователя
	return nil
}
func (r *AuthPostgres) GetByCredentials(ctx context.Context, userAuth entities.UserAuth) (int, error) {
	var id int
	query := fmt.Sprintf(`SELECT id FROM %s WHERE username=$1 and password_hash=$2`, userTable)
	err := r.db.GetContext(ctx, &id, query, userAuth.UserName, userAuth.Password)
	if errors.Is(err, sql.ErrNoRows) {
		//return 0, еррорф ноуроу + запрашиваемый айди
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (r *AuthPostgres) IsLoginFree(ctx context.Context, userAuth entities.UserAuth) (bool, error) {
	var id int
	query := fmt.Sprintf(`SELECT id FROM %s WHERE username=$1`, userTable)
	err := r.db.GetContext(ctx, &id, query, userAuth.UserName)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if id > 0 {
		return true, nil
	}
	return false, nil
}
func (r *AuthPostgres) Verify(ctx context.Context, tokens entities.Tokens) (entities.Tokens, error) {
	return tokens, nil
}
