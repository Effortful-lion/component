package controller

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "component/backend/service"
)

// 定义 JWT 密钥
var jwtSecret = []byte("your_jwt_secret")

// 登录处理
func EmailLogin(c *gin.Context) {
    type RequestBody struct {
        Email string `json:"email" binding:"required,email"`
        Code  string `json:"code" binding:"required"`
    }
    var req RequestBody
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 验证验证码
    if !service.VerifyCode(req.Email, req.Code) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired verification code"})
        return
    }

    // 生成 JWT
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email": req.Email,
        "exp":   time.Now().Add(24 * time.Hour).Unix(),
    })
    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
