package main

//bcrypt: 对于密码存储，bcrypt是一种专为此目的设计的算法，它通过增加计算成本来抵抗暴力破解。Go中可使用golang.org/x/crypto/bcrypt包。

import (
	"fmt"
	"io"
	//"golang.org/x/crypto/bcrypt"		// 专门针对密码进行加密的包
	"crypto/md5"						// md5加密包
)

func main() {
	// 加密前
	secret := "okokokook"			// 加盐操作：防止md5被逆向破解（可有可无，增加安全性）
	str := "123456"
	// 加密
	h := md5.New()
	io.WriteString(h,str)
	io.WriteString(h,secret)
	// 加密后
	entrystr := h.Sum(nil)
	//io.WriteString(h,string(entrystr))	// 接受字符串的参数
	io.Writer.Write(h,entrystr)		// 直接接受字节数组的参数
	entrystr_twice := h.Sum(nil)	// 二次加密
	fmt.Printf("%x\n",entrystr)
	fmt.Printf("%x",entrystr_twice)
}

// func main() {
// 	// 加密前
// 	str := "123456"
// 	salt := 25  // 最大盐值是31，但是盐值越大，加密速度越慢（响应就慢）【所以要适当设置盐值】
// 	// 加密
// 	// 默认盐值是4 ：entrystr,err := bcrypt.GenerateFromPassword([]byte(str),bcrypt.DefaultCost)  bcrypt.DefaultCost == 4
// 	entrystr,err := bcrypt.GenerateFromPassword([]byte(str),salt)
// 	if err != nil {
// 		fmt.Println("加密失败！")
// 		return
// 	}
// 	// 加密后
// 	fmt.Println(string(entrystr))
// }


