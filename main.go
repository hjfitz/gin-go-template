package main

import (
	"os"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func serveStaticPages(r *gin.Engine, dir string) {
	reg, _ := regexp.Compile(".html")
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			filename := reg.ReplaceAllString(info.Name(), "")
			filePath := "/" + filename
			if filename == "index" {
				filePath = "/"
			}
			r.StaticFile(filePath, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("some middleware invoked")
	}
}

func main() {
	err := godotenv.Load()

	r := gin.Default()

	// set up session middleware
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	// static file hosting
	r.Static("/public", "./public")

	// because r.StaticFS on / means that you can't mount any other MW
	serveStaticPages(r, "./pages")

	// anything we need to hide goes here
	r.Use(middleware())

	r.GET("/incr", func(c *gin.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.JSON(200, gin.H{"count": count})
	})


	r.Run(":5000")
}
