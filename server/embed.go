package server

import _ "embed"

//go:embed static/favicon.ico
var favicon []byte

//go:embed static/index.html
var indexHtml string

//go:embed static/vercel.html
var vercelHtml string

func Favicon() []byte {
	return favicon
}

func VercelHtml() string {
	return vercelHtml
}
