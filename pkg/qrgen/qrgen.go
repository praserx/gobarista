// Copyright 2023 PraserX
package qrgen

import (
	"encoding/base64"

	"github.com/dundee/qrpay"
)

// Mailer struct definition.
type PaymentInfo struct {
	IBAN          string
	BIC           string
	Currency      string
	RecipientName string
	Message       string
	VS            string
	Amount        string
}

func GetQRCodeImage(pi PaymentInfo) (qrcode []byte, err error) {
	payment := qrpay.NewSpaydPayment()

	payment.SetIBAN(pi.IBAN)
	payment.SetBIC(pi.BIC)
	payment.SetCurrency(pi.Currency)
	payment.SetRecipientName(pi.RecipientName)
	payment.SetMessage(pi.Message)
	payment.SetExtendedAttribute("vs", pi.VS)
	payment.SetAmount(pi.Amount)

	return qrpay.GetQRCodeImage(payment)
}

func GetQRCodeImageBase64(pi PaymentInfo) (b64 string, err error) {
	data, err := GetQRCodeImage(pi)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString([]byte(data)), nil
}
