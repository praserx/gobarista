# v1.3.0
- Added pay command
- Added pay-csv command
- Minor fixes and changes

# v1.2.0
- Added CI CodeQL and Staticcheck
- Added `AmountPerPackage` to models.Period
- Added `PaymentConfirmation` to models.Bill
- Added `Cash` to models.Period
- Changed TotalAmount computation
- Mail package reorganization
- Added payment confirmation e-mail and command
- Added command for all customers contacts print
- Fixed bulk add and switched firstname and lastname

# v1.1.1
- Version fix

# v1.1.0
- Added models.Schema table and its Version column
- Added credit field to models.User (do `gobarista database migrate`)
- More minor updates and mail layout update

# v1.0.0
- Initial implementation