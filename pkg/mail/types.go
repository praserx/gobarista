package mail

import "github.com/praserx/gobarista/pkg/stats"

// EmailSettings contains essential mail configuration and content.
type EmailSettings struct {
	Subject string
	Plain   string
	HTML    string
}

// BillTemplateVars represents HTML template variables.
type BillTemplateVars struct {
	Title                string
	Subtitle             string
	BID                  string
	UID                  string
	Name                 string
	Location             string
	Credit               string
	PeriodFrom           string
	PeriodTo             string
	UnitPrice            string
	Quantity             string
	Amount               string
	CreditPay            string
	Payment              string
	PaymentAN            string
	PaymentVS            string
	PaymentCustomMessage string
	Stats                stats.Stats
	QRCode               string
	AppVersion           string
}

// ConfirmationTemplateVars represents HTML template variables.
type ConfirmationTemplateVars struct {
	Title      string
	Subtitle   string
	BID        string
	PeriodFrom string
	PeriodTo   string
	Payment    string
	AppVersion string
}

// UnpaidTemplateVars represents HTML template variables.
type UnpaidTemplateVars struct {
	Title      string
	Subtitle   string
	BID        string
	PeriodFrom string
	PeriodTo   string
	Payment    string
	AppVersion string
}
