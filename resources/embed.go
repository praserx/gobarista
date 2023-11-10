// Copyright 2023 PraserX
package resources

import "embed"

//go:embed templates
var DirTemplates embed.FS

const DIR_TEMPLATES = "templates"

// General constants
const HTML_BILL_TEMPLATE = "bill.cs.html"
const HTML_BILL_TEMPLATE_FULL_PATH = DIR_TEMPLATES + "/" + HTML_BILL_TEMPLATE
const HTML_CONFIRM_TEMPLATE = "confirm.cs.html"
const HTML_CONFIRM_TEMPLATE_FULL_PATH = DIR_TEMPLATES + "/" + HTML_CONFIRM_TEMPLATE
const HTML_UNPAID_TEMPLATE = "unpaid.cs.html"
const HTML_UNPAID_TEMPLATE_FULL_PATH = DIR_TEMPLATES + "/" + HTML_UNPAID_TEMPLATE
