// Package web embeds the built frontend files served by Gin.
package web

import "embed"

// Dist contains the production Vite build generated from web/frontend.
//
//go:embed dist
var Dist embed.FS
