package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"job4j_go_share_trip_contracts/internal/domain/contract/entity"
)

type UserContractRepository struct{}

func NewUserContractRepository() *UserContractRepository {
	return &UserContractRepository{}
}

func (r *UserContractRepository) GetOpenByUserIDForUpdate(
	ctx context.Context,
	tx pgx.Tx,
	userID uuid.UUID,
) (*entity.Contract, error) {
	query := `
		SELECT id, company_id, user_id, status, start_date, end_date, created_at
		FROM contracts
		WHERE user_id = $1 AND status IN ('active', 'draft', 'suspended')
		ORDER BY created_at DESC
		LIMIT 1
		FOR UPDATE
	`

	contract := &entity.Contract{}
	var companyID *uuid.UUID
	err := tx.QueryRow(ctx, query, userID).Scan(
		&contract.ID,
		&companyID,
		&contract.UserID,
		&contract.Status,
		&contract.StartDate,
		&contract.EndDate,
		&contract.CreatedAt,
	)
	if companyID != nil {
		contract.CompanyID = *companyID
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get open contract by user: %w", err)
	}

	return contract, nil
}

func (r *UserContractRepository) CreateTx(ctx context.Context, tx pgx.Tx, contract *entity.Contract) error {
	query := `
		INSERT INTO contracts (id, company_id, user_id, status, start_date, end_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := tx.Exec(
		ctx,
		query,
		contract.ID,
		contract.CompanyID,
		contract.UserID,
		contract.Status,
		contract.StartDate,
		contract.EndDate,
		contract.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create contract: %w", err)
	}
	return nil
}

func (r *UserContractRepository) UpdateTx(ctx context.Context, tx pgx.Tx, contract *entity.Contract) error {
	query := `
		UPDATE contracts
		SET status = $2, start_date = $3, end_date = $4
		WHERE id = $1
	`
	_, err := tx.Exec(
		ctx,
		query,
		contract.ID,
		contract.Status,
		contract.StartDate,
		contract.EndDate,
	)
	if err != nil {
		return fmt.Errorf("update contract: %w", err)
	}
	return nil
}
