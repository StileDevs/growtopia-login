package main

import (
	"embed"
	b64 "encoding/base64"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

//go:embed all:static
var assets embed.FS

var Tmpl *template.Template
var Assets fs.FS
var StaticAssets static.ServeFileSystem

func GetAssets() {
	asset, _ := fs.Sub(assets, "static")
	Assets = asset

	tmpl := template.Must(template.ParseFS(Assets, "*.html"))
	Tmpl = tmpl

	StaticAssets = static.EmbedFolder(assets, "static")
}

func init() {
	GetAssets()
}

type LoginData struct {
	Token    string `json:"_token" form:"_token"`
	GrowID   string `json:"growId" form:"growId"`
	Password string `json:"password" form:"password"`
}

func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello world",
		})
	})

	r.Use(static.Serve("/", StaticAssets))
	r.SetHTMLTemplate(Tmpl)

	r.POST("/player/login/dashboard", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "dashboard.html", gin.H{})
	})

	r.POST("/player/growid/login/validate", func(ctx *gin.Context) {

		var data LoginData

		if err := ctx.Bind(&data); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "internal error when parsing json",
			})
			return
		}

		token := b64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("_token=%s&growId=%s&password=%s", data.Token, data.GrowID, data.Password)))

		ctx.Header("Content-Type", "text/html")
		ctx.JSON(http.StatusOK, gin.H{
			"status":      "success",
			"message":     "Account Validated.",
			"token":       token,
			"url":         "",
			"accountType": "growtopia",
		})
	})

	r.GET("/player/validate/close", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "<script>window.close();</script>", gin.H{})
	})

	r.Run(":80")
	// r.Run("localhost:8080")
	// r.RunTLS("localhost:8080", "assets/ssl/_wildcard.growserver.app.pem", "assets/ssl/_wildcard.growserver.app-key.pem")
}
