package gin_test

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Server_Test1() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	//router.Run(":8080")

	//自定义服务器设置
	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
