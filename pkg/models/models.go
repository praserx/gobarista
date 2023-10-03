// Copyright 2023 PraserX
package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`  // GORM default
	CreatedAt time.Time      `json:"timestamp"`             // GORM default
	UpdatedAt time.Time      `json:"-"`                     // GORM default
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`        // GORM default
	EID       string         `gorm:"column:eid" json:"eid"` // Note: Employee ID
	Email     string         `json:"email"`
	Firstname string         `json:"firstname"`
	Lastname  string         `json:"lastname"`
	Location  string         `json:"location"`
}

type Period struct {
	ID            uint           `gorm:"primarykey" json:"id"` // GORM default
	CreatedAt     time.Time      `json:"timestamp"`            // GORM default
	UpdatedAt     time.Time      `json:"-"`                    // GORM default
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`       // GORM default
	DateFrom      time.Time      `json:"date_from"`
	DateTo        time.Time      `json:"date_to"`
	DateOfIssue   time.Time      `json:"date_of_issue"`
	UnitPrice     float32        `json:"unit_price"`
	TotalMonths   int            `json:"total_months"`
	TotalQuantity int            `json:"total_quantity"`
	TotalAmount   float32        `json:"total_amount"`
	Closed        bool           `json:"closed"`
}

type Bill struct {
	ID        uint           `gorm:"primarykey" json:"id"` // GORM default
	CreatedAt time.Time      `json:"timestamp"`            // GORM default
	UpdatedAt time.Time      `json:"-"`                    // GORM default
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`       // GORM default
	Quantity  int            `json:"quantity"`
	Amount    float32        `json:"amount"`
	Issued    bool           `json:"issued"`
	Paid      bool           `json:"paid"`
	UserID    uint           `json:"-"` // GORM reference
	PeriodID  uint           `json:"-"` // GORM reference
}
