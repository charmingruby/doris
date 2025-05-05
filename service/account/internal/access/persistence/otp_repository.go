package persistence

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/charmingruby/doris/lib/persistence/postgres"
	"github.com/charmingruby/doris/service/account/internal/access/core/model"
	"github.com/jmoiron/sqlx"
)

const (
	findOTPByCorrelationID = "find otp by correlation id"
	createOTP              = "create otp"
)

func otpQueries() map[string]string {
	return map[string]string{
		findOTPByCorrelationID: `SELECT * FROM otps WHERE correlation_id = $1 ORDER BY created_at DESC LIMIT 1`,
		createOTP:              `INSERT INTO otps (id, correlation_id, code, purpose, expires_at) VALUES ($1, $2, $3, $4, $5)`,
	}
}

type OTPRepo struct {
	db    postgres.Database
	stmts map[string]*sqlx.Stmt
}

func NewOTPRepo(db postgres.Database) (*OTPRepo, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range otpQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			return nil,
				postgres.NewPreparationErr(queryName, "api key", err)
		}

		stmts[queryName] = stmt
	}

	return &OTPRepo{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *OTPRepo) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		return nil,
			postgres.NewStatementNotPreparedErr(queryName, "api key")
	}

	return stmt, nil
}

func (r *OTPRepo) Create(ctx context.Context, otp model.OTP) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(createOTP)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx, otp.ID, otp.CorrelationID, otp.Code, otp.Purpose, otp.ExpiresAt); err != nil {
		return err
	}

	return nil
}

func (r *OTPRepo) FindMostRecentByCorrelationID(ctx context.Context, correlationID string) (model.OTP, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	stmt, err := r.statement(findOTPByCorrelationID)
	if err != nil {
		return model.OTP{}, err
	}

	var otp model.OTP
	if err := stmt.GetContext(ctx, &otp, correlationID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.OTP{}, nil
		}

		return model.OTP{}, err
	}

	return otp, nil
}
