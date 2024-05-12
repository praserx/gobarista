// Copyright 2023 PraserX
package resources

import "embed"

//go:embed assets
var DirAssets embed.FS

//go:embed templates
var DirTemplates embed.FS

const DIR_TEMPLATES_EMAIL = "templates/email"
const DIR_TEMPLATES_WEB = "templates/web"

// General constants
const HTML_BILL_TEMPLATE = "bill.cs.html"
const HTML_BILL_TEMPLATE_FULL_PATH = DIR_TEMPLATES_EMAIL + "/" + HTML_BILL_TEMPLATE
const HTML_CONFIRM_TEMPLATE = "confirm.cs.html"
const HTML_CONFIRM_TEMPLATE_FULL_PATH = DIR_TEMPLATES_EMAIL + "/" + HTML_CONFIRM_TEMPLATE
const HTML_UNPAID_TEMPLATE = "unpaid.cs.html"
const HTML_UNPAID_TEMPLATE_FULL_PATH = DIR_TEMPLATES_EMAIL + "/" + HTML_UNPAID_TEMPLATE
