// Copyright 2023 PraserX
package mail

import (
	"fmt"

	"github.com/praserx/gobarista/pkg/config"
	"github.com/praserx/gobarista/pkg/logger"
	"github.com/praserx/gobarista/pkg/models"
	"github.com/praserx/gobarista/pkg/qrgen"
	"github.com/praserx/gobarista/pkg/rank"
)

func SendBill(user models.User, period models.Period, bill models.Bill, totalCustomers int) error {
	var err error
	var pinfo qrgen.PaymentInfo
	var tvars TemplateVars

	pinfo.IBAN = config.Get().Section("spayd").Key("iban").String()
	pinfo.BIC = config.Get().Section("spayd").Key("bic").String()
	pinfo.Currency = config.Get().Section("spayd").Key("currency").String()
	pinfo.RecipientName = config.Get().Section("spayd").Key("recipient_name").String()
	pinfo.Message = config.Get().Section("spayd").Key("message").String()
	pinfo.Amount = fmt.Sprintf("%.2f", bill.Amount)
	pinfo.VS = fmt.Sprintf("%d", bill.ID)

	tvars.BID = fmt.Sprintf("%d", bill.ID)
	tvars.UID = fmt.Sprintf("%d", user.ID)
	tvars.Name = fmt.Sprintf("%s %s", user.Firstname, user.Lastname)
	tvars.Location = user.Location
	tvars.Rank = rank.ComputeRank(bill.Quantity / period.TotalMonths)
	tvars.PeriodFrom = period.DateFrom.Format("2. 1. 2006")
	tvars.PeriodTo = period.DateTo.Format("2. 1. 2006")
	tvars.UnitPrice = fmt.Sprintf("%.2f", period.UnitPrice)
	tvars.Quantity = fmt.Sprintf("%d", bill.Quantity)
	tvars.Amount = fmt.Sprintf("%.2f", bill.Amount)
	tvars.TotalMonths = fmt.Sprintf("%d", period.TotalMonths)
	tvars.TotalQuantity = fmt.Sprintf("%d", period.TotalQuantity)
	tvars.TotalAverage = fmt.Sprintf("%d", period.TotalQuantity/totalCustomers)
	tvars.TotalCustomers = fmt.Sprintf("%d", totalCustomers)
	tvars.PaymentAN = config.Get().Section("spayd").Key("an").String()
	tvars.PaymentVS = fmt.Sprintf("%d", bill.ID)
	tvars.PaymentCustomMessage = config.Get().Section("spayd").Key("custom_message").String()
	tvars.QRCode, err = qrgen.GetQRCodeImageBase64(pinfo)

	if err != nil {
		return fmt.Errorf("could not generate qr code: %v", err)
	}

	es := &EmailSettings{
		Subject: config.Get().Section("messages").Key("subject_bill").String(),
		Plain:   config.Get().Section("messages").Key("no_plaintext").String(),
		HTML:    GetBillHTMLTemplate(tvars),
	}

	if err = SendMail(user.Email, es); err != nil {
		logger.Error("cannot sent mail: " + err.Error())
	}

	return nil
}
