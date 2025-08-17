package models

import (
    "time"

    "github.com/shopspring/decimal"
)

// User represents a user in the system.
type User struct {
    ID           int
    Username     string
    PasswordHash string
    Role         string
}

// Account represents a financial account owned by a user.
type Account struct {
    ID      int
    UserID  int
    Balance decimal.Decimal
}

// Transaction records fund transfers between accounts.
type Transaction struct {
    ID            int
    FromAccountID int
    ToAccountID   int
    Amount        decimal.Decimal
    CreatedAt     time.Time
}
