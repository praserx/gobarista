// Copyright 2023 PraserX
package commands

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/praserx/gobarista/pkg/cmd/gobarista/helpers"
	"github.com/praserx/gobarista/pkg/database"
	"github.com/praserx/gobarista/pkg/logger"
	"github.com/praserx/gobarista/pkg/mail"
	"github.com/praserx/gobarista/pkg/models"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

var Billing = cli.Command{
	Name:  "billing",
	Usage: "Billing periods and bills management",
	Subcommands: []*cli.Command{
		&BillingNewPeriod,
		&BillingAddBill,
		&BillingClosePeriod,
		&BillingPeriodSummary,
		&BillingIssueBills,
	},
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		users, err := database.SelectAllUsers()
		if err != nil {
			return fmt.Errorf("error: cannot get users: %v", err.Error())
		}

		for _, user := range users {
			fmt.Printf("%s %s\t\t%s\t\t%s\n", user.Firstname, user.Lastname, user.Email, user.Location)
		}

		return nil
	},
}

var BillingNewPeriod = cli.Command{
	Name:      "create-period",
	Usage:     "Add new billing period",
	ArgsUsage: "[from(2006-01-02) to(2006-01-02) issued(2006-01-02) total_amount]",
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		if ctx.NArg() != 4 {
			return fmt.Errorf("error: too few arguments: requires (4), get (%d)", ctx.NArg())
		}

		DateFrom, err := time.Parse("2006-01-02", ctx.Args().Get(0))
		if err != nil {
			return fmt.Errorf("error: cannot parse date: %v", err)
		}

		DateTo, err := time.Parse("2006-01-02", ctx.Args().Get(1))
		if err != nil {
			return fmt.Errorf("error: cannot parse date: %v", err)
		}

		DateOfIssue, err := time.Parse("2006-01-02", ctx.Args().Get(2))
		if err != nil {
			return fmt.Errorf("error: cannot parse date: %v", err)
		}

		TotalAmount, err := strconv.ParseFloat(ctx.Args().Get(3), 32)
		if err != nil {
			return fmt.Errorf("error: cannot parse float: %v", err)
		}

		period := models.Period{
			DateFrom:    DateFrom,
			DateTo:      DateTo,
			DateOfIssue: DateOfIssue,
			TotalMonths: int(DateTo.Sub(DateFrom).Hours() / 24 / 30),
			TotalAmount: float32(TotalAmount),
		}

		id, err := database.InsertPeriod(period)
		if err != nil {
			return fmt.Errorf("error: cannot create new billing period: %v", err.Error())
		}
		logger.Info(fmt.Sprintf("new billing period successfully created: new period id: %d", id))

		return nil
	},
}

var BillingClosePeriod = cli.Command{
	Name:      "close-period",
	Usage:     "Close billing period and calculate remaining values such as total quantity or unit price",
	ArgsUsage: "[id]",
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		if ctx.NArg() != 1 {
			return fmt.Errorf("error: too few arguments: requires (1), get (%d)", ctx.NArg())
		}

		pid, err := strconv.ParseUint(ctx.Args().Get(0), 10, 64)
		if err != nil {
			return fmt.Errorf("error: cannot parse uint: %v", err)
		}

		period, err := database.SelectPeriodByID(uint(pid))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error: billing period not found")
		} else if err != nil {
			return fmt.Errorf("error: cannot get billing period: %v", err)
		}

		bills, err := database.SelectAllBillsForPeriod(uint(pid))
		if err != nil {
			return fmt.Errorf("error: cannot get bills for given period: %v", err)
		}

		totalQuantity := 0
		for _, bill := range bills {
			totalQuantity += bill.Quantity
		}

		unitPrice := period.TotalAmount / float32(totalQuantity)

		err = database.UpdatePeriodOnClose(uint(pid), totalQuantity, unitPrice)
		if err != nil {
			return fmt.Errorf("error: cannot create new billing period: %v", err.Error())
		}
		logger.Info(fmt.Sprintf("billing period updated and closed successfully: period_id=%d, total_quantity=%d, unit_price=%.2f", pid, totalQuantity, unitPrice))

		for _, bill := range bills {
			amount := float32(bill.Quantity) * unitPrice
			err = database.UpdateBillOnPeriodClose(bill.ID, amount)
			if err != nil {
				logger.Error(fmt.Sprintf("error: cannot update bill with bill_id=%d: %v", bill.ID, err.Error()))
			} else {
				logger.Info(fmt.Sprintf("bill has been updated successfully: bill_id=%d amount=%.2f", bill.ID, amount))
			}
		}

		return nil
	},
}

var BillingPeriodSummary = cli.Command{
	Name:      "period-summary",
	Usage:     "Get numbers and statistics for given billing period",
	ArgsUsage: "[id]",
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		if ctx.NArg() != 1 {
			return fmt.Errorf("error: too few arguments: requires (1), get (%d)", ctx.NArg())
		}

		pid, err := strconv.ParseUint(ctx.Args().Get(0), 10, 64)
		if err != nil {
			return fmt.Errorf("error: cannot parse uint: %v", err)
		}

		period, err := database.SelectPeriodByID(uint(pid))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error: billing period not found")
		} else if err != nil {
			return fmt.Errorf("error: cannot get billing period: %v", err)
		}

		bills, err := database.SelectAllBillsForPeriod(uint(pid))
		if err != nil {
			return fmt.Errorf("error: cannot get bills for given period: %v", err)
		}

		fmt.Printf("ID:     		%d\n", period.ID)
		fmt.Printf("Period: 		%s - %s\n", period.DateFrom.Format("2006-01-02"), period.DateTo.Format("2006-01-02"))
		fmt.Printf("Unit price: 	%.2f\n", period.UnitPrice)
		fmt.Printf("Total quantity: %d\n", period.TotalQuantity)
		fmt.Printf("Total amount: 	%.2f\n", period.TotalAmount)
		fmt.Printf("Total months:   %d\n", period.TotalMonths)
		fmt.Printf("Closed: 		%t\n", period.Closed)
		fmt.Printf("Total bills: 	%d\n", len(bills))

		return nil
	},
}

var BillingAddBill = cli.Command{
	Name:      "add-bill",
	Usage:     "Add bill for given user and period with specified quantity",
	ArgsUsage: "[user_id period_id quantity]",
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		if ctx.NArg() != 3 {
			return fmt.Errorf("error: too few arguments: requires (4), get (%d)", ctx.NArg())
		}

		uid, err := strconv.ParseUint(ctx.Args().Get(0), 10, 64)
		if err != nil {
			return fmt.Errorf("error: cannot parse uint: %v", err)
		}

		pid, err := strconv.ParseUint(ctx.Args().Get(1), 10, 64)
		if err != nil {
			return fmt.Errorf("error: cannot parse uint: %v", err)
		}

		quantity, err := strconv.ParseInt(ctx.Args().Get(2), 10, 64)
		if err != nil {
			return fmt.Errorf("error: cannot parse int: %v", err)
		}

		user, err := database.SelectUserByID(uint(uid))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error: user not found")
		} else if err != nil {
			return fmt.Errorf("error: cannot get user: %v", err)
		}

		period, err := database.SelectPeriodByID(uint(pid))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error: billing period not found")
		} else if err != nil {
			return fmt.Errorf("error: cannot get billing period: %v", err)
		}

		if period.Closed {
			return fmt.Errorf("error: billing period is already closed")
		}

		bill := models.Bill{
			Quantity: int(quantity),
			UserID:   uint(uid),
			PeriodID: uint(pid),
		}

		id, err := database.InsertBill(bill)
		if err != nil {
			return fmt.Errorf("error: cannot add new bill to billing period: %v", err.Error())
		}
		logger.Info(fmt.Sprintf("new bill successfully added to billing period for user_id=%d, user_name=%s: new bill id: %d", uid, user.Firstname+" "+user.Lastname, id))

		return nil
	},
}

var BillingIssueBills = cli.Command{
	Name:      "issue-bills",
	Usage:     "Issue all bills for giver billing period",
	ArgsUsage: "[id]",
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		if err = helpers.SetupMail(ctx); err != nil {
			return err
		}

		if ctx.NArg() != 1 {
			return fmt.Errorf("error: too few arguments: requires (1), get (%d)", ctx.NArg())
		}

		pid, err := strconv.ParseUint(ctx.Args().Get(0), 10, 64)
		if err != nil {
			return fmt.Errorf("error: cannot parse uint: %v", err)
		}

		period, err := database.SelectPeriodByID(uint(pid))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error: billing period not found")
		} else if err != nil {
			return fmt.Errorf("error: cannot get billing period: %v", err)
		}

		if !period.Closed {
			return fmt.Errorf("error: cannot issue bills: period is not closed")
		}

		bills, err := database.SelectAllBillsForPeriod(uint(pid))
		if err != nil {
			return fmt.Errorf("error: cannot get bills for given period: %v", err)
		}

		for _, bill := range bills {
			user, err := database.SelectUserByID(bill.UserID)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("error: user for giver bill is not found: user_id=%d", bill.UserID)
			} else if err != nil {
				return fmt.Errorf("error: cannot get user: user_id=%d: %v", bill.UserID, err)
			}

			mail.SendBill(user, period, bill, len(bills))

			err = database.UpdateBillOnIssued(bill.ID)
			if err != nil {
				return fmt.Errorf("error: cannot update bill: bill_id=%d: %v", bill.ID, err)
			}
		}

		return nil
	},
}
