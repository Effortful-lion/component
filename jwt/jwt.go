package jwt

import (
    "component/redis"
    "fmt"
    "time"
    "github.com/golang-jwt/jwt/v4"
)

// 用于签名的字符串(jwt的签名密钥)
var mySigningKey = []byte("jwt_secret")

// 自定义声明
type MyClaims struct {
    Username string `json:"username"`
    jwt.RegisteredClaims
}

// 生成自定义声明
func GenMyClaims() *MyClaims {
    return &MyClaims{
        // 自定义声明：
        Username: "username",
        // 内嵌标准声明：
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 20)), // jwt过期时间
            Issuer:    "username",                                         // 发布者
        },
    }
}

// 生成token
func GenMyToken() (string, error) {
    // 生成MyClaims
    claims := GenMyClaims()
    // 生成token对象：需要claims声明对象和签名方法：提供了 加密方法 和 特定的加密信息
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    // 生成签名字符串
    return token.SignedString(mySigningKey)
}

// 校验token
func ParseMytoken(tokenString string) bool {
    // 用jwt的函数：使用声明解析jwt令牌，得到token
    token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
        return mySigningKey, nil
    })
    if err!= nil {
        fmt.Println(err)
    }
    return token.Valid
}

// 为了在一段时间内不断的刷新token，保证token的安全性，我们设置一个refresh token
// refresh负责 每隔一段时间 进行 token 的刷新即生成新的token并同步到redis中
func GenMyRefresh_token() (string, error) {
    claims := GenMyClaims()
    refresh_token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
    return refresh_token.SignedString(mySigningKey)  // refresh的密钥和token的密钥不可以
}

// 定时任务（打点器[某个时间间隔重复执行]） + go（监听）
// 使用定时任务定时更新token，以refresh_token在redis中的过期时间作为期限，并每20秒进行一次token的刷新
func StartTokenRefreshWithTicker() {
    ticker := time.NewTicker(time.Second * 20) // 每隔20秒触发一次，可根据实际需求调整时间间隔

    // 生成Refresh_token并存储到redis，设置过期时间，同时处理可能的错误
    refresh_token, err := GenMyRefresh_token()
    if err!= nil {
        fmt.Printf("生成refresh_token失败: %v\n", err)
        return
    }
    fmt.Println("refresh_token:", refresh_token)
    err = redis.Rdb.Set("refresh_token", refresh_token, time.Second*40).Err()
    if err!= nil {
        fmt.Printf("将refresh_token存入redis失败: %v\n", err)
        return
    }

    go func() {
        for {
            // 先检查refresh_token是否过期，如果过期则退出循环，结束刷新操作
            isExpired, err := checkRefreshTokenExpired()
            if err!= nil {
                fmt.Printf("检查refresh_token过期情况失败: %v\n", err)
                continue
            }
            if isExpired {
                fmt.Println("refresh_token过期，账密重新校验信息")
                return
            }

            newToken, err := GenMyToken()
            if err!= nil {
                fmt.Printf("生成新token失败: %v\n", err)
                continue
            }
            fmt.Printf("成功刷新token，新token: %s\n", newToken)
            // 这里可以添加将新token同步到需要的地方的逻辑，比如存储到某个全局变量或者发送给客户端等
            // 重置redis中的token，在web程序中（不要忘记还要及时向前端返回token，否则前端无法使用）
            err = redis.Rdb.Set("token", newToken, time.Second*20).Err()
            if err!= nil {
                fmt.Printf("将新token存入redis失败: %v\n", err)
                continue
            }

            <-ticker.C // 等待下一个时间间隔到达
        }
    }()
}

// 检查refresh_token是否过期
func checkRefreshTokenExpired() (bool, error) {
    remainingTime, err := redis.Rdb.TTL("refresh_token").Result()
    if err!= nil {
        return false, err
    }
    return remainingTime <= 0, nil
}