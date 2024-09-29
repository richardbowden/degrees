package models

import "time"

type Account struct {
	id int

	firstName   string
	middleName  string
	surname     string
	Email       string
	AccountType int
	SignUpStage int
	Enabled     bool
	CreateedOn  time.Time
	UpdatedAt   time.Time
}

type AccountWithAuth struct {
	Account
	PasswordHash string
}

type EmailAddress struct {
	Id          int
	Email       string
	Verified    bool
	VerifitedOn time.Time
	UpdatedOn   time.Time
}

type EmailAddresses []EmailAddress

type AccountFull struct {
	Account
	EmailAddresses
}

type EmailVerificationCode struct {
	Id   int
	Code string
}
