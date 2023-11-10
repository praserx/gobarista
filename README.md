# GoBarista

GoBarista is simple billing command line application.

## The story behind

Far away in Czechia there are many government institutions which not have free coffee. These clerks and officers working in these institutions has to buy coffee by them self. But there is one super officer which taking care of coffee for his colleagues. He buys the coffee for over time, and sometimes he has to do billing. This application is exactly for this.

## Installation

```
go install github.com/praserx/gobarista/pkg/cmd/gobarista@latest
```

### Configuration

The configuration is stored as standard INI file. See below for more information.

```INI
[paths]
database = /path/to/my/sqlite3/database/db.sqlite3

[spayd]
an = "1234567890/1234"
iban = "CZ0000000000000000001234"
bic = "YYXXCZPP"
currency = czk
recipient_name = "John Doe"
message = "Coffee"
custom_message = "Some additional message for bill if needed"

[messages]
subject_bill = "Coffee bill | John's Coffee Service"
subject_confirm = "Payment confirmation | John's Coffee Service"
company_name = "John's Coffee Service"
subtitle_msg = "Since 2021"
no_plaintext = "This message is only at HTML format"

[smtp]
host = smtp.example.com
port = 25
username = johndoe@example.com
password = password
name = "John's Coffe Service"
from = johndoe@example.com

```

## Example usage

At first you can use `gobarista help` for help message.

### Preparation and checks

```bash
gobarista -c .local/conf.ini help
gobarista -c .local/conf.ini database initialize
gobarista -c .local/conf.ini database migrate
gobarista -c .local/conf.ini users add 001001 John Doe johndoe@example.com Tokyo
```

### Billing

```bash
gobarista -c .local/conf.ini billing period-create 2023-10-01 2023-10-30 2023-11-01 500 250
gobarista -c .local/conf.ini billing add-bill 1 1 15
gobarista -c .local/conf.ini billing period-close 1
gobarista -c .local/conf.ini billing period-summary 1
gobarista -c .local/conf.ini billing issue-bills 1
gobarista -c .local/conf.ini billing confirm-payment 1
```
