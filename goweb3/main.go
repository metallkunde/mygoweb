package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()

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

	err := r.Run(":19999")
	if err != nil {
		return
	}
}
