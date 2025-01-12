package main

// 雪花算法生成id

// 复习：雪花算法是由64位组成：符号位(1)、时间戳(41)、机器码(5[数据中心]+5[机器ID])、计数器(12)

import (
	"fmt"
	//"time"

	"github.com/bwmarrin/snowflake" // 雪花算法包: 基于twitter的snowflake算法实现
)

// 使用步骤：
// 1. 导包"github.com/bwmarrin/snowflake"
// 2. 使用NewNode() 创建一个节点，初始化雪花节点
// 3. 使用Generate() 生成id

func main() {
	// Node number must be between 0 and 1023 : 机器码只有10位 （2^10-1）
	node, err := snowflake.NewNode(1023)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate a snowflake ID.
	id := node.Generate()

	// Print out the ID in a few different ways.
	fmt.Printf("Int64  ID: %d\n", id)
	fmt.Printf("String ID: %s\n", id)

}

// 其中如果想用在一个项目中，其实这样也可以，但是还可以具备个性化，不是完全使用默认配置

// // Epoch is set to the twitter snowflake epoch of Nov 04 2010 01:42:54 UTC in milliseconds
// You may customize this to set a different epoch for your application.
//Epoch int64 = 1288834974657
// 这是默认配置，记录的是 一个时间的毫秒数，默认是 2010年11月04日 01:42:54 UTC 的时间
// 我们也可以记录自定义的时间，比如：2019年11月04日 01:42:54 UTC 的时间

// var (
// 	// 全局节点
// 	snowflakeNode *snowflake.Node
// )

// // 二次封装 初始化节点的函数
// func Init(startTime string, machineID int64) (node *snowflake.Node,err error){
// 	var st time.Time
// 	st, err = time.Parse("2006-01-02", startTime) // 将字符串转换成时间格式
// 	if err != nil {
// 		fmt.Println("startTime err:", err)
// 		return
// 	}
// 	snowflake.Epoch = st.UnixNano() / 1000000	// 设置自定义的时间
// 	node, err = snowflake.NewNode(machineID) // 传入机器码
// 	if err != nil {
// 		fmt.Println("NewNode err:", err)
// 		return
// 	}
// 	return 
// }

// // 二次封装 生成id的函数
// func GenID() int64{
// 	return snowflakeNode.Generate().Int64()
// }

// func main() {
// 	// 初始化节点
// 	node, err := Init("2019-11-04", 1)
// 	if err != nil {
// 		fmt.Println("Init err:", err)
// 		return
// 	}
// 	snowflakeNode = node

// 	// 生成id
// 	id := GenID()
// 	fmt.Println("id:", id)
// }
