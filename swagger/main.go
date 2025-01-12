package main

import (
	"github.com/gin-gonic/gin"

    swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
    _ "component/swagger/docs"

)

// @可以修改的swagger文档的标题
// @version 版本(默认1.0)
// @description  这里是swagger中整个项目的描述
// @termsOfService https://www.swagger.io/terms/
 
// @contact.name  维护者名字
// @contact.url http://www.swagger.io/support
// @contact.email 维护者邮件
 
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
 
// @host 127.0.0.1:8080
// @BasePath /api/v1
func main() {
    r := gin.Default()
    //在代码中处理请求的接口函数(通常位于controller层)，按如下方式注释。

    //注册swagger api相关路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // 设置路由组/api/v1
    g := r.Group("/api/v1")
    {
        g.GET("/ping",PingHandler)
        
        g.POST("/login",LoginHandler)
        
    }

    r.Run("127.0.0.1:8080")
}

// 模型（model）
type User struct{
    Username string    `json:"username" example:"张三"`
    Password string    `json:"password" example:"123456"`
}

// 响应参数模型（model）
type ResponseParam struct {
    Code int `json:"code"`
    Message string `json:"message"`
    Data *User  `json:"data"` 
}

// 接收参数模型（model）
type RequestParam struct {
    Username string    `json:"username" example:"李四"`
    Password string    `json:"password" example:"123456"`
}

// 接收参数模型2（model）
type RequestParam2 struct {
    Username string    `json:"username" example:"李四"`
    Password string    `json:"password" example:"123456"`
    UpdateUsername string `json:"updateusername" example:"张三"`
    UpdatePassword string `json:"updatepassword" example:"789100"`
}

// @Summary 测试ping
    // @Description ping
    // @Tags 测试1
    // @Accept application/json
    // @Success 200 {string} string "{"msg": "pong"}"
    // @Failure 400 {string} string "{"msg": "nonono"}"
    // @Router /ping [get]
func PingHandler(ctx *gin.Context){
    ctx.JSON(200,gin.H{
        "code": 1000,
       "msg": "pong",
    })
}

    // @Summary 测试login
    // @Description login
    // @Tags 测试2
    // @Accept application/json
    // @Produce application/json
    // @Param Object body RequestParam true "用户信息"
    // @Success 200 {object} ResponseParam
    // @Failure 400 {string} string "{"msg": "failed"}"
    // @Router /login [post]
func LoginHandler(c *gin.Context){
    user := new(User)
    err := c.ShouldBindJSON(user)
    if err != nil{
        c.JSON(200,gin.H{
            "code": 5000,
            "msg": user,
        })
    }
    c.JSON(200,gin.H{
        "code": 2000,
        "msg": user,
    })
}

    // @Summary 删除用户
    // @Description deleteuser
    // @Tags 测试3
    // @Accept application/json
    // @Produce application/json
    // @Param id path int true "用户id"
    // Success 200 {string} string "{"msg": "该id="{{id}}"用户已经被注销"}"
    // Failure 400 {string} string "{"msg": "该id不存在"}"
    // @Router /user [delete]
func DeleteUser(c *gin.Context){
    // 传来待删除用户的id(路径参数)
    id := c.Param("id")
    ok := (id == "1")
    if ok {
        c.JSON(200,gin.H{
            "code": 5000,
            "msg": ("该id="+id+"用户不存在"),
        })
    }
    c.JSON(200,gin.H{
        "code": 2000,
        "msg": ("该id="+id+"用户已经被注销"),
    })

}

    //@Summary 更新用户信息
    // @Description update
    // @Tags 测试4
    // @Accept application/json
    // @Produce application/json
    // @Param Object body RequestParam2 true "用户更新后信息"
    // @Success 200 {object} ResponseParam
    // @Failure 400 {string} string "{"msg": "failed"}"
    // @Router /update [put]
func UpdateUser(c *gin.Context){

}

