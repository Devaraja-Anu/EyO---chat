package data

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IdentityKey struct {
	UserID    int64
	PublicKey []byte
	CreatedAt time.Time
}

type IdentityKeyModel struct {
	DB *pgxpool.Pool
}

var ErrIdentityKeyExists = errors.New("identity key already exists")
var ErrIdentityKeyNotFound = errors.New("identity key not found")

func (m IdentityKeyModel) Insert(ctx context.Context, id_key *IdentityKey) error {

	query := `INSERT INTO identity_keys (user_id,public_key) VALUES ($1,$2) RETURNING created_at`

	err := m.DB.QueryRow(ctx, query, id_key.UserID, id_key.PublicKey).Scan(&id_key.CreatedAt)
	
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// 23505 is the error code for a unique violation in PostgreSQL
			if pgErr.Code == "23505" {
				return ErrIdentityKeyExists
			}
		}
		return err
	}

	return nil
}


func(m IdentityKeyModel) getIdentityKey(ctx context.Context, userID int64) (*IdentityKey,error){

	query := `
	SELECT user_id,public_key, created_at
	 FROM identity_keys 
	 WHERE user_id = $1`

	 var idKey IdentityKey

	 err := m.DB.QueryRow(ctx,query,userID).Scan(&idKey.UserID, &idKey.PublicKey,&idKey.CreatedAt)

	 if err != nil {
		if(errors.Is(err,pgx.ErrNoRows)){
			return  nil,ErrIdentityKeyNotFound
		}
		return nil, err
	 }

	return &idKey,nil	 

}