package main

import (
	"component/redis"

	controller "component/backend/controller"
	"github.com/gin-gonic/gin"
)

func main() {

	redis.InitRedis()

	//redis.SetKey("key", "value")
	

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.NoRoute(func(c *gin.Context){
		c.JSON(404, gin.H{
			"msg": "404,您的也页面好像找不到了~",
		})
	})

	r.POST("/send-code", controller.SendEmailCode)
    r.POST("/login", controller.EmailLogin)

	r.Run("127.0.0.1:8080")
	
}