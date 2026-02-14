package data

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
    ID        int64     `json:"id"`
    Username  string    `json:"username"`
    CreatedAt time.Time `json:"created_at"`
}

var ErrDuplicateUsername = errors.New("duplicate username")
var ErrUserNotFound  = errors.New("user not Found")
type UserModel struct {
	DB *pgxpool.Pool
}

func (m UserModel) Insert(ctx context.Context, user *User) error {

	query :=  `INSERT INTO USERS (username) VALUES ($1) RETURNING id, username, created_at`

	 err := m.DB.QueryRow(ctx,query,user.Username).Scan(&user.ID,&user.Username,&user.CreatedAt)

	 if err != nil {
		var pgErr  *pgconn.PgError
		if errors.As(err, &pgErr){
			// 23505 is the error code for a unique violation in PostgreSQL
			if pgErr.Code == "23505" {
				return ErrDuplicateUsername
			}
		}
		return err
	 }

	 return nil
}

func (m UserModel) GetByUsername(ctx context.Context,username string) (*User, error){

	query := `SELECT id,username, created_at FROM users where username  = $1`

	var userData User

	err := m.DB.QueryRow(ctx,query,username).Scan(&userData.ID,&username,&userData.CreatedAt)

	if err != nil {
		if  errors.Is(err,pgx.ErrNoRows){
			return nil, ErrUserNotFound
		}
	}

	return  &userData, nil
}