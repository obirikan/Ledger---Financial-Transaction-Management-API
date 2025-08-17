package repository

import (
	"context"

	"github.com/obirikan/Ledger---Financial-Transaction-Management-API/internal/models"
	"github.com/shopspring/decimal"
)

type Repository interface {
	// User
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)

	// Account
	CreateAccount(ctx context.Context, account *models.Account) error
	GetAccountsByUserID(ctx context.Context, userID int) ([]*models.Account, error)
	GetAccountByID(ctx context.Context, id int) (*models.Account, error)

	// Transactions & Ledger
	TransferFunds(ctx context.Context, fromAccountID, toAccountID int, amount decimal.Decimal) error
}
