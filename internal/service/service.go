package service

import (
	"context"
	"errors"

	"github.com/obirikan/Ledger---Financial-Transaction-Management-API/internal/models"
	"github.com/obirikan/Ledger---Financial-Transaction-Management-API/internal/repository"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)

// AuthService defines auth business logic.
type AuthService interface {
	Register(ctx context.Context, username, password string) error
	Login(ctx context.Context, username, password string) (string, error) // JWT token
}

// LedgerService defines ledger business logic.
type LedgerService interface {
	ListAccounts(ctx context.Context, userID int) ([]*models.Account, error)
	TransferFunds(ctx context.Context, userID int, fromAccountID, toAccountID int, amount decimal.Decimal) error
}

type authService struct {
	repo repository.Repository
}

func NewAuthService(repo repository.Repository) AuthService {
	return &authService{repo: repo}
}

func (a *authService) Register(ctx context.Context, username, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &models.User{
		Username:     username,
		PasswordHash: string(hashed),
		Role:         "user",
	}
	return a.repo.CreateUser(ctx, user)
}

func (a *authService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := a.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	token, err := GenerateJWT(user.ID, user.Role)
	return token, err
}

type ledgerService struct {
	repo repository.Repository
}

func NewLedgerService(repo repository.Repository) LedgerService {
	return &ledgerService{repo: repo}
}

func (l *ledgerService) ListAccounts(ctx context.Context, userID int) ([]*models.Account, error) {
	return l.repo.GetAccountsByUserID(ctx, userID)
}

func (l *ledgerService) TransferFunds(ctx context.Context, userID int, fromAccountID, toAccountID int, amount decimal.Decimal) error {
	fromAccount, err := l.repo.GetAccountByID(ctx, fromAccountID)
	if err != nil {
		return err
	}
	if fromAccount.UserID != userID {
		return errors.New("permission denied: source account does not belong to user")
	}
	if amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount must be positive")
	}
	// Check destination account existence
	if _, err := l.repo.GetAccountByID(ctx, toAccountID); err != nil {
		return err
	}
	return l.repo.TransferFunds(ctx, fromAccountID, toAccountID, amount)
}
