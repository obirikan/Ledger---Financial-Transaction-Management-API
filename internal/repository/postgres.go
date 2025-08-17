package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/obirikan/Ledger---Financial-Transaction-Management-API/internal/models"
	"github.com/shopspring/decimal"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) Repository {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) CreateUser(ctx context.Context, user *models.User) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO users(username, password_hash, role) VALUES ($1, $2, $3)`,
		user.Username, user.PasswordHash, user.Role)
	return err
}

func (r *PostgresRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	row := r.db.QueryRowContext(ctx, `SELECT id, username, password_hash, role FROM users WHERE username=$1`, username)
	if err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *PostgresRepo) CreateAccount(ctx context.Context, acc *models.Account) error {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO accounts(user_id, balance) VALUES ($1, $2) RETURNING id`,
		acc.UserID, acc.Balance.String(),
	).Scan(&acc.ID)
	return err
}

func (r *PostgresRepo) GetAccountsByUserID(ctx context.Context, userID int) ([]*models.Account, error) {

	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, balance FROM accounts WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*models.Account{}
	for rows.Next() {
		acc := &models.Account{}
		var balStr string
		if err := rows.Scan(&acc.ID, &acc.UserID, &balStr); err != nil {
			return nil, err
		}
		balance, err := decimal.NewFromString(balStr)
		if err != nil {
			return nil, err
		}
		acc.Balance = balance
		accounts = append(accounts, acc)
	}
	return accounts, nil
}

func (r *PostgresRepo) GetAccountByID(ctx context.Context, id int) (*models.Account, error) {
	acc := &models.Account{}
	var balStr string
	row := r.db.QueryRowContext(ctx, `SELECT id, user_id, balance FROM accounts WHERE id=$1`, id)
	if err := row.Scan(&acc.ID, &acc.UserID, &balStr); err != nil {
		return nil, err
	}
	balance, err := decimal.NewFromString(balStr)
	if err != nil {
		return nil, err
	}
	acc.Balance = balance
	return acc, nil
}

func (r *PostgresRepo) TransferFunds(ctx context.Context, fromID, toID int, amount decimal.Decimal) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Lock from account row for update
	var fromBalanceStr string
	err = tx.QueryRowContext(ctx, `SELECT balance FROM accounts WHERE id=$1 FOR UPDATE`, fromID).Scan(&fromBalanceStr)
	if err != nil {
		tx.Rollback()
		return err
	}
	fromBalance, err := decimal.NewFromString(fromBalanceStr)
	if err != nil {
		tx.Rollback()
		return err
	}
	if fromBalance.LessThan(amount) {
		tx.Rollback()
		return errors.New("insufficient funds")
	}

	// Subtract from source account
	_, err = tx.ExecContext(ctx, `UPDATE accounts SET balance = balance - $1 WHERE id = $2`, amount.String(), fromID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Add to destination account
	_, err = tx.ExecContext(ctx, `UPDATE accounts SET balance = balance + $1 WHERE id = $2`, amount.String(), toID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert transaction log
	_, err = tx.ExecContext(ctx,
		`INSERT INTO transactions (from_account_id, to_account_id, amount, created_at) VALUES ($1, $2, $3, $4)`,
		fromID, toID, amount.String(), time.Now().UTC())
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
