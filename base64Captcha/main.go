package main

import (
	// 标准库，提供图像颜色处理
	//"image/color"
	"net/http"
	"time"

	// 三方库解决跨域问题
	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	// 三方库实现图片验证码  别名base64Captcha  可支持字母、数字的验证数据
	base64Captcha "github.com/mojocn/base64Captcha"
)

// Math 配置参数
// var (
// 	Height          = 70
// 	Width           = 240
// 	NoiseCount      = 0
// 	ShowLineOptions = base64Captcha.OptionShowHollowLine
// 	BgColor         = &color.RGBA{
// 		R: 144,
// 		G: 238,
// 		B: 144,
// 		A: 10,
// 	}
// 	FontsStorage base64Captcha.FontsStorage
// 	Fonts        []string
// )

// 全局变量  用于模拟存储用户信息的Mysql数据库
var users = make(map[string]string)

var store = base64Captcha.DefaultMemStore
// 音频的验证码自行探索

//var driver = base64Captcha.DefaultDriverDigit // 数字

//自定义验证码：
//var driver = base64Captcha.NewDriverMath(Height, Width, NoiseCount, ShowLineOptions, BgColor, FontsStorage, Fonts)	// 字母
var driver = base64Captcha.NewDriverString(
	80, 240, 6, 1, 4,
	"123456789abcdefghjkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ", // 包含数字和字母的字符集
	nil, nil, nil,
) // 字符串：数字+字母

// 生成验证码并返回图片
func generateCaptcha(c *gin.Context) {
	captcha := base64Captcha.NewCaptcha(driver, store)
	// 生成验证码（返回的 id：验证码id，base64编码的验证码图片，验证码答案，错误信息）
	id, b64, _, err := captcha.Generate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "验证码生成失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"captcha_id":  id,
		"captcha_img": b64,
	})
}

// 注册接口
func register(c *gin.Context) {
	var json struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		CaptchaId  string `json:"captcha_id"`
		CaptchaAns string `json:"captcha_ans"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式错误"})
		return
	}

	// 验证验证码
	if !store.Verify(json.CaptchaId, json.CaptchaAns, true) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}

	// 注册用户（简单模拟）
	if _, exists := users[json.Username]; exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户已存在"})
		return
	}

	// 存储用户数据
	users[json.Username] = json.Password

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

// 登录接口
func login(c *gin.Context) {
	var json struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		CaptchaId  string `json:"captcha_id"`
		CaptchaAns string `json:"captcha_ans"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据格式错误"})
		return
	}

	// 验证验证码（参数：id、answer不用说了；clear：使用后是否清除验证码 ）
	if !store.Verify(json.CaptchaId, json.CaptchaAns, true) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误"})
		return
	}

	// 验证用户名和密码
	storedPassword, exists := users[json.Username]
	if !exists || storedPassword != json.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "登录成功"})
}

func main() {
	r := gin.Default()

	// 配置 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500"}, // 前端的 URL
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 注册路由
	r.GET("/captcha", generateCaptcha)
	r.POST("/register", register)
	r.POST("/login", login)

	// 启动服务器
	r.Run("127.0.0.1:8080")
}
