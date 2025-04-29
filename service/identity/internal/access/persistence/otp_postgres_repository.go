package persistence

import (
	"context"
	"database/sql"
	"errors"

	"github.com/charmingruby/doris/lib/persistence/postgres"
	"github.com/charmingruby/doris/service/identity/internal/access/core/model"
	"github.com/jmoiron/sqlx"
)

const (
	findOTPByCorrelationID = "find otp by correlation id"
	createOTP              = "create otp"
)

func otpQueries() map[string]string {
	return map[string]string{
		findOTPByCorrelationID: `SELECT * FROM otps WHERE correlation_id = $1`,
		createOTP:              `INSERT INTO otps (id, correlation_id, code, purpose, expires_at) VALUES ($1, $2, $3, $4, $5)`,
	}
}

type OTPPostgresRepo struct {
	db    *sqlx.DB
	stmts map[string]*sqlx.Stmt
}

func NewOTPPostgresRepo(db *sqlx.DB) (*OTPPostgresRepo, error) {
	stmts := make(map[string]*sqlx.Stmt)

	for queryName, statement := range otpQueries() {
		stmt, err := db.Preparex(statement)
		if err != nil {
			return nil,
				postgres.NewPreparationErr(queryName, "api key", err)
		}

		stmts[queryName] = stmt
	}

	return &OTPPostgresRepo{
		db:    db,
		stmts: stmts,
	}, nil
}

func (r *OTPPostgresRepo) statement(queryName string) (*sqlx.Stmt, error) {
	stmt, ok := r.stmts[queryName]

	if !ok {
		return nil,
			postgres.NewStatementNotPreparedErr(queryName, "api key")
	}

	return stmt, nil
}

func (r *OTPPostgresRepo) Create(ctx context.Context, otp model.OTP) error {
	stmt, err := r.statement(createOTP)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(otp.ID, otp.CorrelationID, otp.Code, otp.Purpose, otp.ExpiresAt); err != nil {
		return err
	}

	return nil
}

func (r *OTPPostgresRepo) FindByCorrelationID(ctx context.Context, correlationID string) (model.OTP, error) {
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
