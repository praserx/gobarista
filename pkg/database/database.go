// Copyright 2023 PraserX
package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/praserx/gobarista/pkg/models"
)

// Global database instance
var gdb *gorm.DB

// SetupDatabase set ups database with defined arguments and initialize global
// instance.
func SetupDatabase(opts ...Option) {
	if gdb == nil {
		gdb, _ = newDatabase(opts...)
	}
}

// newDatabase create new database instance, make connection and check newly
// established connection to database. If anything goes wrong, error will be
// returned.
func newDatabase(opts ...Option) (*gorm.DB, error) {
	var err error

	var options = &DatabaseOptions{
		Path: "",
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.Path == "" {
		return nil, fmt.Errorf("database: missing path to database file")
	}

	db, err := gorm.Open(sqlite.Open(options.Path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	return db, nil
}

// Get returns database instance.
func Get() *gorm.DB {
	return gdb
}

// Close gracefully end up connection to database.
func Close() {
	dbi, err := gdb.DB()
	if err != nil {
		dbi.Close()
	}
}

// RunAutoMigration creates all required tables in database.
func RunAutoMigration() error {
	return gdb.AutoMigrate(
		&models.Schema{},
		&models.User{},
		&models.Period{},
		&models.Bill{},
	)
}

func InsertSchema(schema models.Schema) (int, error) {
	obj := schema
	result := gdb.Create(&obj)
	return int(obj.Version), result.Error
}

func SelectVersion() (schema models.Schema, err error) {
	result := gdb.First(&schema)
	return schema, result.Error
}

func UpdateVersion(version uint) error {
	var schema models.Schema
	gdb.First(&schema)
	return gdb.
		Model(&schema).
		Update("Version", version).Error
}

func SelectAllUsers() (users []models.User, err error) {
	result := gdb.Find(&users)
	return users, result.Error
}

func SelectUserByID(id uint) (user models.User, err error) {
	result := gdb.First(&user, id)
	return user, result.Error
}

func SelectUserByEmail(email string) (user models.User, err error) {
	result := gdb.Where("email = ?", email).First(&user)
	return user, result.Error
}

func SelectUserByEID(eid string) (user models.User, err error) {
	result := gdb.Where("eid = ?", eid).First(&user)
	return user, result.Error
}

// InsertUser inserts user into the database and returns ID of the inserted
// record and error if any occurs.
func InsertUser(user models.User) (int, error) {
	obj := user
	result := gdb.Create(&obj)
	return int(obj.ID), result.Error
}

func SelectPeriodByID(id uint) (period models.Period, err error) {
	result := gdb.First(&period, id)
	return period, result.Error
}

// InsertUser inserts user into the database and returns ID of the inserted
// record and error if any occurs.
func InsertPeriod(period models.Period) (int, error) {
	obj := period
	result := gdb.Create(&obj)
	return int(obj.ID), result.Error
}

func UpdatePeriodOnClose(pid uint, totalQuantity int, unitPrice float32) error {
	var period models.Period
	gdb.First(&period, pid)
	return gdb.
		Model(&period).
		Updates(models.Period{TotalQuantity: totalQuantity, UnitPrice: unitPrice, Closed: true}).Error
}

func SelectBillByID(id uint) (bill models.Bill, err error) {
	result := gdb.First(&bill, id)
	return bill, result.Error
}

func SelectAllBills() (bills []models.Bill, err error) {
	result := gdb.Find(&bills)
	return bills, result.Error
}

func SelectAllBillsForPeriod(pid uint) ([]models.Bill, error) {
	var bill []models.Bill
	result := gdb.
		Where("period_id = ?", pid).
		Find(&bill)
	return bill, result.Error
}

func InsertBill(bill models.Bill) (int, error) {
	obj := bill
	result := gdb.Create(&obj)
	return int(obj.ID), result.Error
}

func UpdateBillOnPeriodClose(bid uint, amount float32) error {
	var bill models.Bill
	gdb.First(&bill, bid)
	return gdb.
		Model(&bill).
		Update("Amount", amount).Error
}

func UpdateBillOnIssued(bid uint) error {
	var bill models.Bill
	gdb.First(&bill, bid)
	return gdb.
		Model(&bill).
		Update("Issued", true).Error
}

func UpdateBillOnPaid(bid uint) error {
	var bill models.Bill
	gdb.First(&bill, bid)
	return gdb.
		Model(&bill).
		Update("Paid", true).Error
}
