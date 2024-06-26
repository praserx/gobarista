// Copyright 2023 PraserX
package models

import (
	"time"

	"gorm.io/gorm"
)

// VERSION of database schema
const VERSION = uint(3)

// Transaction types
const (
	UNKNOWN  = 0
	DEPOSIT  = 1
	WITHDRAW = 2
)

type Schema struct {
	Version uint `gorm:"primarykey" json:"version"`
}

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`  // GORM default
	CreatedAt time.Time      `json:"timestamp"`             // GORM default
	UpdatedAt time.Time      `json:"-"`                     // GORM default
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`        // GORM default
	EID       string         `gorm:"column:eid" json:"eid"` // Note: Employee ID
	Email     string         `json:"email"`                 // User e-mail
	Firstname string         `json:"firstname"`             // User firstname
	Lastname  string         `json:"lastname"`              // User lastname
	Location  string         `json:"location"`              // User workplace location
	Credit    int            `json:"credit"`                // User credit
}

type Transaction struct {
	ID        uint           `gorm:"primarykey" json:"id"` // GORM default, Transactions ID
	CreatedAt time.Time      `json:"timestamp"`            // GORM default
	UpdatedAt time.Time      `json:"-"`                    // GORM default
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`       // GORM default
	Type      int            `json:"type"`                 // Transaction type (e.g. deposit or withdraw)
	Amount    float32        `json:"amount"`               // Total amount for transaction
	UserID    uint           `json:"-"`                    // GORM reference
}

type Period struct {
	ID               uint           `gorm:"primarykey" json:"id"` // GORM default
	CreatedAt        time.Time      `json:"timestamp"`            // GORM default
	UpdatedAt        time.Time      `json:"-"`                    // GORM default
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`       // GORM default
	DateFrom         time.Time      `json:"date_from"`            // Billing period from
	DateTo           time.Time      `json:"date_to"`              // Billing period to
	DateOfIssue      time.Time      `json:"date_of_issue"`        // Date of issue (used at bill date of issue)
	UnitPrice        float32        `json:"unit_price"`           // Unit price of coffee
	TotalMonths      int            `json:"total_months"`         // Total months of billing period
	TotalQuantity    int            `json:"total_quantity"`       // Total quantity of coffees
	TotalAmount      float32        `json:"total_amount"`         // Total amount (cost) of coffee packages
	AmountPerPackage float32        `json:"amount_per_package"`   // Average price of coffee package
	Cash             float32        `json:"cash"`                 // Some coffees were paid by cash
	Closed           bool           `json:"closed"`               // Is billing period close/finished?
}

type Bill struct {
	ID                  uint           `gorm:"primarykey" json:"id"` // GORM default
	CreatedAt           time.Time      `json:"timestamp"`            // GORM default
	UpdatedAt           time.Time      `json:"-"`                    // GORM default
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`       // GORM default
	Quantity            int            `json:"quantity"`             // How many coffees user drank
	Amount              float32        `json:"amount"`               // Total amount for user
	Payment             float32        `json:"payment"`              // Amount of money to pay (Bill.Amount - User.Credit)
	Issued              bool           `json:"issued"`               // Is bill issued/send to user?
	Paid                bool           `json:"paid"`                 // Is bill paid?
	PaymentConfirmation bool           `json:"payment_confirmation"` // Has confirmation of payment been sent?
	UserID              uint           `json:"-"`                    // GORM reference
	PeriodID            uint           `json:"-"`                    // GORM reference
}
