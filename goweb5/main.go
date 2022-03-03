package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.New()
	r.Use(gee.Logger())
	/*
		curl -i localhost:19999

		HTTP/1.1 200 OK
		Content-Type: text/html
		Date: Tue, 22 Feb 2022 09:03:43 GMT
		Content-Length: 20
	*/
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello There</h1>")
	})

	/*
		curl localhost:19999/hello?name=kunder

		hello kunder, you are at /hello
	*/
	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you are at %s\n", c.Query("name"), c.Path)
	})

	/*
		curl localhost:19999/login -X POST -d "username=kunder&password=114514"

		{"password":"114514","username":"kunder"}
	*/
	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	/*
		curl localhost:19999/assets/test.txt

		{"filepath":"test.txt"}
	*/
	r.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})

	r.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})
	/*
		v1 := r.Group("/v1")
		{
			v1.GET("/", func(c *gee.Context) {
				c.HTML(http.StatusOK, "<h1>Hello There</h1>")
			})

			v1.GET("/hello", func(c *gee.Context) {
				// expect /hello?name=kunder
				c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
			})
		}
	*/
	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/kunder
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}
	err := r.Run(":19999")
	if err != nil {
		return
	}
}
