// Copyright 2023 PraserX
package commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/praserx/gobarista/pkg/cmd/gobarista/flags"
	"github.com/praserx/gobarista/pkg/cmd/gobarista/helpers"
	"github.com/praserx/gobarista/pkg/database"
	"github.com/praserx/gobarista/pkg/logger"
	"github.com/praserx/gobarista/pkg/mail"
	"github.com/praserx/gobarista/pkg/models"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

var Billing = cli.Command{
	Name:    "billing",
	Aliases: []string{"b"},
	Usage:   "Billing periods and bills management",
	Flags: []cli.Flag{
		&flags.FlagPlainPrint,
	},
	Subcommands: []*cli.Command{
		&BillingPeriods,
		&BillingBills,
	},
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		bills, err := database.SelectAllBills()
		if err != nil {
			return fmt.Errorf("error: cannot get bills: %v", err.Error())
		}

		if !ctx.Bool("plain") {
			t := table.NewWriter()
			t.AppendHeader(table.Row{"Bill ID", "Period ID", "Amount", "Issued", "Paid", "Confirmed"})
			for _, bill := range bills {
				t.AppendRow(table.Row{bill.ID, bill.PeriodID, bill.Amount, bill.Issued, bill.Paid, bill.PaymentConfirmation})
			}
			fmt.Println(t.Render())
		} else {
			for _, bill := range bills {
				fmt.Printf("%d %d %.2f %t %t %t\n", bill.ID, bill.PeriodID, bill.Amount, bill.Issued, bill.Paid, bill.PaymentConfirmation)
			}
		}

		return nil
	},
}

var BillingPeriods = cli.Command{
	Name:    "periods",
	Aliases: []string{"p"},
	Usage:   "Billing periods management (create, close, summary)",
	Subcommands: []*cli.Command{
		&BillingPeriodsNew,
		&BillingPeriodsClose,
		&BillingPeriodsSummary,
	},
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		periods, err := database.SelectAllPeriods()
		if err != nil {
			return fmt.Errorf("error: cannot get billing periods: %v", err.Error())
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"ID", "From", "To", "UnitPrice", "Total Amount", "Total Quantity", "Avg. Package Price", "Cash"})
		for _, period := range periods {
			t.AppendRow(table.Row{period.ID, period.DateFrom.Format("2006-01-02"), period.DateTo.Format("2006-01-02"), fmt.Sprintf("%.2f", period.UnitPrice), fmt.Sprintf("%.2f", period.TotalAmount), fmt.Sprintf("%d", period.TotalQuantity), fmt.Sprintf("%.2f", period.AmountPerPackage), fmt.Sprintf("%.2f", period.Cash)})
		}
		fmt.Println(t.Render())

		return nil
	},
}

var BillingPeriodsNew = cli.Command{
	Name:      "new",
	Aliases:   []string{"n"},
	Usage:     "Add new billing period",
	ArgsUsage: "[from(2006-01-02) to(2006-01-02) issued(2006-01-02) total_amount amount_per_package]",
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		if ctx.NArg() != 5 {
			return fmt.Errorf("error: too few arguments: requires (5), get (%d)", ctx.NArg())
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

		AmountPerPackage, err := strconv.ParseFloat(ctx.Args().Get(4), 32)
		if err != nil {
			return fmt.Errorf("error: cannot parse float: %v", err)
		}

		totalMonths := 1
		if int(DateTo.Sub(DateFrom).Hours()/24/30) > 0 {
			totalMonths = int(DateTo.Sub(DateFrom).Hours() / 24 / 30)
		}

		period := models.Period{
			DateFrom:         DateFrom,
			DateTo:           DateTo,
			DateOfIssue:      DateOfIssue,
			TotalMonths:      totalMonths,
			TotalAmount:      float32(TotalAmount),
			AmountPerPackage: float32(AmountPerPackage),
		}

		id, err := database.InsertPeriod(period)
		if err != nil {
			return fmt.Errorf("error: cannot create new billing period: %v", err.Error())
		}
		logger.Info(fmt.Sprintf("new billing period successfully created: new period id: %d", id))

		return nil
	},
}

var BillingPeriodsClose = cli.Command{
	Name:      "close",
	Aliases:   []string{"c"},
	Usage:     "Close billing period and calculate remaining values such as total quantity or unit price",
	ArgsUsage: "[period_id]",
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		if ctx.NArg() != 1 {
			return fmt.Errorf("error: too few arguments: requires (1), get (%d)", ctx.NArg())
		}

		pid, err := strconv.ParseUint(ctx.Args().Get(0), 10, 32)
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

		unitPrice := (period.TotalAmount - period.Cash) / float32(totalQuantity)

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

var BillingPeriodsSummary = cli.Command{
	Name:      "summary",
	Aliases:   []string{"s"},
	Usage:     "Get numbers and statistics for given billing period",
	ArgsUsage: "[period_id]",
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		if ctx.NArg() != 1 {
			return fmt.Errorf("error: too few arguments: requires (1), get (%d)", ctx.NArg())
		}

		pid, err := strconv.ParseUint(ctx.Args().Get(0), 10, 32)
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

		fmt.Printf("ID:                %d\n", period.ID)
		fmt.Printf("Period:            %s - %s\n", period.DateFrom.Format("2006-01-02"), period.DateTo.Format("2006-01-02"))
		fmt.Printf("Unit price:        %.2f\n", period.UnitPrice)
		fmt.Printf("Total quantity:    %d\n", period.TotalQuantity)
		fmt.Printf("Total amount:      %.2f\n", period.TotalAmount)
		fmt.Printf("Total amount (wc): %.2f (total - cash)\n", period.TotalAmount)
		fmt.Printf("Total months:      %d\n", period.TotalMonths)
		fmt.Printf("Cash:              %.2f\n", period.Cash)
		fmt.Printf("Closed:            %t\n", period.Closed)
		fmt.Printf("Total bills:       %d\n", len(bills))

		return nil
	},
}

var BillingBills = cli.Command{
	Name:    "bills",
	Aliases: []string{"b"},
	Usage:   "Bills management section (add, issue, pay, confirm)",
	Subcommands: []*cli.Command{
		&BillingBillsAdd,
		&BillingBillsIssue,
		&BillingBillsPay,
		&BillingBillsPayCSV,
		&BillingBillsConfirmPayment,
		&BillingBillsUnpaidNotification,
	},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "print all bills overtime",
		},
		&cli.BoolFlag{
			Name:    "paid",
			Aliases: []string{"p"},
			Usage:   "print all bills which were already paid",
		},
		&cli.BoolFlag{
			Name:    "unpaid",
			Aliases: []string{"u"},
			Usage:   "print all bills which waiting for payment",
		},
		&cli.IntFlag{
			Name:    "sort",
			Aliases: []string{"s"},
			Usage:   "sort output by quantity (1: high-low, 2: low-high) or amount (3: high-low, 4:low-high)",
		},
	},
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		bills, err := database.SelectAllBills()
		if err != nil {
			return fmt.Errorf("error: cannot get billing periods: %v", err.Error())
		}

		var display []models.Bill
		var quantity int
		var amount float32

		for _, bill := range bills {
			if ctx.Bool("all") {
				quantity += bill.Quantity
				amount += bill.Amount
				display = append(display, bill)
			} else if ctx.Bool("paid") && bill.Paid {
				quantity += bill.Quantity
				amount += bill.Amount
				display = append(display, bill)
			} else if ctx.Bool("unpaid") && !bill.Paid {
				quantity += bill.Quantity
				amount += bill.Amount
				display = append(display, bill)
			}
		}

		switch variant := ctx.Int("sort"); variant {
		case 1:
			sort.Slice(display, func(i, j int) bool {
				return display[i].Quantity > display[j].Quantity
			})
		case 2:
			sort.Slice(display, func(i, j int) bool {
				return display[i].Quantity < display[j].Quantity
			})
		case 3:
			sort.Slice(display, func(i, j int) bool {
				return display[i].Amount > display[j].Amount
			})
		case 4:
			sort.Slice(display, func(i, j int) bool {
				return display[i].Amount < display[j].Amount
			})
		default:
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"ID", "PID", "UID", "Name", "Quantity", "Amount", "Issued", "Paid", "Payment Confirmed"})
		for _, bill := range display {
			user, err := database.SelectUserByID(bill.UserID)
			if err != nil {
				return fmt.Errorf("error: cannot get user: %v", err.Error())
			}

			t.AppendRow(table.Row{bill.ID, bill.PeriodID, bill.UserID, user.Firstname + " " + user.Lastname, fmt.Sprintf("%d", bill.Quantity), fmt.Sprintf("%.2f", bill.Amount), fmt.Sprintf("%t", bill.Issued), fmt.Sprintf("%t", bill.Paid), fmt.Sprintf("%t", bill.PaymentConfirmation)})
		}
		t.AppendFooter(table.Row{"", "", "", "", fmt.Sprintf("%d", quantity), fmt.Sprintf("%.2f", amount), "", "", ""})
		fmt.Println(t.Render())

		return nil
	},
}

var BillingBillsAdd = cli.Command{
	Name:      "add",
	Aliases:   []string{"a"},
	Usage:     "Add bill for given user and period with specified quantity",
	ArgsUsage: "[period_id user_id quantity]",
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		if ctx.NArg() != 3 {
			return fmt.Errorf("error: too few arguments: requires (4), get (%d)", ctx.NArg())
		}

		pid, err := strconv.ParseUint(ctx.Args().Get(0), 10, 32)
		if err != nil {
			return fmt.Errorf("error: cannot parse uint: %v", err)
		}

		uid, err := strconv.ParseUint(ctx.Args().Get(1), 10, 64)
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

var BillingBillsIssue = cli.Command{
	Name:      "issue",
	Aliases:   []string{"i"},
	Usage:     "Issue all bills for giver billing period",
	ArgsUsage: "[period_id]",
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

		pid, err := strconv.ParseUint(ctx.Args().Get(0), 10, 32)
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

			if err = mail.SendBill(user, period, bill, len(bills)); err != nil {
				logger.Error(fmt.Sprintf("error: billing e-mail has not been sent for bill_id=%d user_id=%d user_name='%s' user_email:'%s'", bill.ID, user.ID, user.Firstname+" "+user.Lastname, user.Email))
			} else {
				logger.Info(fmt.Sprintf("billing e-mail has been sent for bill_id=%d user_id=%d user_name='%s' user_email:'%s'", bill.ID, user.ID, user.Firstname+" "+user.Lastname, user.Email))

				err = database.UpdateBillOnIssued(bill.ID)
				if err != nil {
					return fmt.Errorf("error: cannot update bill: bill_id=%d: %v", bill.ID, err)
				}
			}
		}

		return nil
	},
}

var BillingBillsPay = cli.Command{
	Name:      "pay",
	Aliases:   []string{"p"},
	Usage:     "Pay (add payment tag) bills for giver billing period",
	ArgsUsage: "[bill_id]",
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		if ctx.NArg() != 1 {
			return fmt.Errorf("error: too few arguments: requires (1), get (%d)", ctx.NArg())
		}

		bid, err := strconv.ParseUint(ctx.Args().Get(0), 10, 32)
		if err != nil {
			return fmt.Errorf("error: cannot parse uint: %v", err)
		}

		err = database.UpdateBillOnPaid(uint(bid))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("error: bill not found: bill_id=%d", bid)
		} else if err != nil {
			return fmt.Errorf("error: cannot update bill_id=%d: %v", bid, err)
		}

		logger.Info(fmt.Sprintf("bill successfully marked as paid: bill_id=%d", bid))

		return nil
	},
}

var BillingBillsPayCSV = cli.Command{
	Name:      "pay-csv",
	Usage:     "Bulk pay (add payment tag) by CSV with appropriate variable symbol",
	ArgsUsage: "[csv]",
	Action: func(ctx *cli.Context) (err error) {
		if err = helpers.SetupDatabase(ctx); err != nil {
			return err
		}

		var fr *os.File

		if fr, err = os.Open(ctx.Args().Get(0)); err != nil {
			fmt.Fprintf(os.Stderr, "cannot read file")
			os.Exit(1)
		}

		reader := bufio.NewReader(fr)
		defer func() {
			fr.Close()
		}()

		csv := [][]string{}
		for {
			line, _, err := reader.ReadLine()
			if err != nil {
				break
			}
			csv = append(csv, strings.Split(string(line), ";"))
		}

		csv = csv[1:] // Throw away CSV header
		for i, line := range csv {
			unbid := strings.ReplaceAll(line[12], "\"", "")
			if len(unbid) == 0 {
				logger.Info(fmt.Sprintf("missing variable symbol, skipping: line=%d", i))
				continue
			}

			bid, err := strconv.ParseUint(unbid, 10, 64)
			if err != nil {
				return fmt.Errorf("error: cannot parse uint: %v", err)
			}
			bill, err := database.SelectBillByID(uint(bid))
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Error(fmt.Sprintf("error: bill not found: bill_id=%d", bid))
			} else if err != nil {
				return fmt.Errorf("error: get bill with bill_id=%d: %v", bid, err)
			}

			if !bill.Paid {
				logger.Info(fmt.Sprintf("bill is not paid: bill_id=%d", bid))
			}
		}

		return nil
	},
}

var BillingBillsConfirmPayment = cli.Command{
	Name:      "confirm",
	Aliases:   []string{"c"},
	Usage:     "Send payment confirmation for all paid bills for given period (e-mail will not be sent for already notified users)",
	ArgsUsage: "[period_id]",
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

		pid, err := strconv.ParseUint(ctx.Args().Get(0), 10, 32)
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
			return fmt.Errorf("error: cannot send payment confirmation: period is not closed")
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

			if !bill.PaymentConfirmation && bill.Paid {
				if err = mail.SendPaymentConfirmation(user, period, bill); err != nil {
					logger.Error(fmt.Sprintf("error: billing e-mail has not been sent for bill_id=%d user_id=%d user_name='%s' user_email:'%s'", bill.ID, user.ID, user.Firstname+" "+user.Lastname, user.Email))
				} else {
					logger.Info(fmt.Sprintf("payment confirmation has been sent for bill_id=%d user_id=%d user_name='%s' user_email:'%s'", bill.ID, user.ID, user.Firstname+" "+user.Lastname, user.Email))

					err = database.UpdateBillOnPaymentConfirmation(bill.ID)
					if err != nil {
						return fmt.Errorf("error: cannot update bill: bill_id=%d: %v", bill.ID, err)
					}
				}
			}
		}

		return nil
	},
}

var BillingBillsUnpaidNotification = cli.Command{
	Name:      "notify",
	Aliases:   []string{"n"},
	Usage:     "Send notification e-mail for unpaid debt for all unpaid bills for given period",
	ArgsUsage: "[period_id]",
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

		pid, err := strconv.ParseUint(ctx.Args().Get(0), 10, 32)
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
			return fmt.Errorf("error: cannot send notification: period is not closed")
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

			if !bill.Paid {
				if err = mail.SendUnpaidNotification(user, period, bill); err != nil {
					logger.Error(fmt.Sprintf("error: notification e-mail has not been sent for bill_id=%d user_id=%d user_name='%s' user_email:'%s'", bill.ID, user.ID, user.Firstname+" "+user.Lastname, user.Email))
				} else {
					logger.Info(fmt.Sprintf("debt notification has been sent for bill_id=%d user_id=%d user_name='%s' user_email:'%s'", bill.ID, user.ID, user.Firstname+" "+user.Lastname, user.Email))
				}
			}
		}

		return nil
	},
}
