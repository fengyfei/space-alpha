package main

import (
	"github.com/gin-gonic/gin"
	handler "github.com/silverswords/quicksilver/yuque"
)

func main() {
	router := gin.Default()
	c := NewC()
	c.RegisterRouter(router)
	router.Run()
}

// Controller -
type Controller struct{}

// NewC -
func NewC() *Controller {
	return &Controller{}
}

// RegisterRouter -
func (c Controller) RegisterRouter(r gin.IRouter) {

	r.GET("/getlist", getlist)
	r.GET("/getdetails", getdetails)

}
func getlist(c *gin.Context) {

	handler.BookList(c.Writer, c.Request)

}

func getdetails(c *gin.Context) {

	handler.BookDetail(c.Writer, c.Request)

}
