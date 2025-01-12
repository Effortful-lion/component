package service

import (
    "time"
	"component/redis"
)

// 保存验证码到 Redis
func SaveVerificationCode(email, code string, expiration time.Duration) error {
	rdb := redis.GetRedis()
    return rdb.Set("email_code:"+email, code, expiration).Err()
}

// 验证验证码
func VerifyCode(email, code string) bool {
	rdb := redis.GetRedis()
    storedCode, err := rdb.Get("email_code:"+email).Result()
    if err != nil || storedCode != code {
        return false
    }
    return true
}
