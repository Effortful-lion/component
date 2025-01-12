package ratelimit

// 基于gin 框架 编写限流中间件 并 使用
// 限流中间件：漏桶 和 令牌桶

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	lingpai "github.com/juju/ratelimit"
)

// 令牌桶中间件
// 限流中间件: 每 fillInterval 的时间，可最多处理 1 个请求(相当于对处理请求的 打点器)
// 特点： 如果一段时间内不请求，token会存储，直到桶的最大容量cap，应对突发多个请求
func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
    // 创建令牌桶 ： 参数： 1. 填充时间（单位时间） 2. 容量
    bucket := lingpai.NewBucket(fillInterval, cap)
    return func(c *gin.Context) {
        // 判断是否限流(如果取不到 token 就限流)
        if bucket.TakeAvailable(1) <= 0 {
            c.String(http.StatusOK, "你点击的太快了，请慢点点击~")
            // 终止执行
            c.Abort()
            return
        }
        // 取到令牌继续执行
        c.Next()
    }
}

