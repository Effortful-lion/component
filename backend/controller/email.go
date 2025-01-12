package controller

import (
    "fmt"
    "math/rand"
    "net/http"
    "time"

    "component/backend/service"
    "component/email"

    "github.com/gin-gonic/gin"
)

// 生成随机验证码
func generateVerificationCode() string {
    rand.Seed(time.Now().UnixNano())
    code := fmt.Sprintf("%06d", rand.Intn(1000000))
    fmt.Println("Generated verification code:", code)
    return code // 6位数字验证码
}

// 发送验证码
func SendEmailCode(c *gin.Context) {
    type RequestBody struct {
        Email string `json:"email" binding:"required,email"`
    }
    var req RequestBody
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 生成验证码
    code := generateVerificationCode()

    // 保存验证码到 Redis（或者数据库）
    if err := service.SaveVerificationCode(req.Email, code, 15*time.Minute); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save verification code"})
        return
    }

    // 导入邮件发送模块发送邮件
    emailSender := email.NewEmailSender(
        "smtp.qq.com",               // SMTP 服务器
        465,                         // 端口号（QQ邮箱通常使用 SSL 加密）
        "1106764332@qq.com",         // 发件邮箱
        "gdgwjsrtqklmhjhi",          // 授权码（非邮箱密码）
        "1106764332@qq.com",         // 发件邮箱（默认发件人）
    )
    err := emailSender.SendEmail(req.Email, "Your Verification Code", "Your code is: "+code)
    if err != nil {
        fmt.Println("Failed to send email:", err) // 添加错误日志
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Verification code sent"})
}


