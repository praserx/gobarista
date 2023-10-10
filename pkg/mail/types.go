package mail

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
	Rank                 string
	PeriodFrom           string
	PeriodTo             string
	UnitPrice            string
	Quantity             string
	Amount               string
	PaymentAN            string
	PaymentVS            string
	PaymentCustomMessage string
	TotalMonths          string
	TotalQuantity        string
	TotalAverage         string
	TotalCustomers       string
	QRCode               string
}

// ConfirmationTemplateVars represents HTML template variables.
type ConfirmationTemplateVars struct {
	Title      string
	Subtitle   string
	BID        string
	PeriodFrom string
	PeriodTo   string
	Amount     string
}
