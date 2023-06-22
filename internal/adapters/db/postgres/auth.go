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

/*
type newKey int

const (
	newKeyType newKey = iota
)



func (r *AuthPostgres) CreateUser(ctx context.Context, userAuth entities.UserAuth) (int, error) {
	var id int
	//ctx.Done()
	//ctx.Value("asd")
	//context.WithValue(ctx, newKeyType, "asd")
	userAuth.Password = hasher.GeneratePasswordHash(userAuth.Password)
	query := fmt.Sprintf(`INSERT INTO %s (username, password_hash) values ($1, $2) returning id`,
		userTable)
	//r.db.QueryRowContext()
	// Этот метод, если что выйдет из выполнения
	row := r.db.QueryRow(query, userAuth.UserName, userAuth.Password)
	//А этот не ориентируется на контекст и будет выполняться, пока БД не вернёт ошибку
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// название метода должно быть IsUserExists -> (bool, error)
// а если это GetUser то возвращаем сущность пользователя
func (r *AuthPostgres) GetUser(user entities.UserAuth) (int, error) {
	var userId int
	user.Password = hasher.GeneratePasswordHash(user.Password)
	query := fmt.Sprintf(`SELECT id FROM %s WHERE username=$1 and password_hash=$2`, userTable)
	err := r.db.Get(&userId, query, user.UserName, user.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (r *AuthPostgres) UpdateSession(refreshToken string, expireAt time.Time) (bool, error) {
	query := fmt.Sprintf(`INSERT INTO %s (refresh_token, expired_at) values ($1, $2)`, sessionsTable)
	r.db.QueryRow(query, refreshToken, expireAt)
	//Почитать про соединения, потому что их нужно закрывать
	return true, nil
}

func (r *AuthPostgres) GetSession(refreshToken string) (time.Time, error) {
	fmt.Println(refreshToken)
	var expiredAt time.Time
	var result string
	query := fmt.Sprintf(`SELECT expired_at FROM %s WHERE refresh_token = $1`, sessionsTable)
	err := r.db.Get(&result, query, refreshToken)
	fmt.Println(result)
	if err != nil {
		return time.Now().AddDate(2000, 01, 01), err
	}
	return expiredAt, nil
}
*/
