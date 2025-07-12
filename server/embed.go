package server

import _ "embed"

//go:embed static/index.html
var indexHtml string

//go:embed static/vercel.html
var vercelHtml string

func VercelHtml() string {
	return vercelHtml
}
