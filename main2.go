package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // 在创建之前，我们首先去初始化一个文件静态地址，Gitee 自带一个 StaticFS 函数可以满足我们的需求
    r.StaticFS("/kaiyuan", http.Dir("/opt/server/nginx-1.18/html/kaiyuan"))
    // 这个是放在系统初始化的时候去执行的

	r.LoadHTMLFiles("templates/index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

    r.POST("/upload", upload)

    r.Run("127.0.0.1:8080")
}

func upload(c *gin.Context) {
    // 接收上传的文件
    file, err := c.FormFile("file") // file 为表单中的 name 属性值
    if err!= nil {
        c.JSON(500, gin.H{
            "message": "上传失败",
        })
        return
    }

    // 接收到文件之后，我们需要对该文件做一个简单的解析，不是所有的文件都要接收，我们只接收后缀为.jpg.png.jpeg 格式的图片
    if file.Header.Get("Content-Type")!= "image/jpeg" && file.Header.Get("Content-Type")!= "image/png" && file.Header.Get("Content-Type")!= "image/jpg" {
        c.JSON(500, gin.H{
            "message": "上传文件格式不正确",
        })
        return
    }

    // 保存到本地文件
    // c.SaveUploadedFile(file,"./upload/"+file.Filename)

    // 假设要保存到的网络位置是 "http://example.com/upload/"
    networkLocation := "https://gitee.com/api/v5/repos/young-lion/picture-bed/contents/" + file.Filename

    // 保存文件到本地的临时目录
    savePath := "./upload/" + file.Filename
    if err := c.SaveUploadedFile(file, savePath); err!= nil {
        c.JSON(500, gin.H{
            "message": "保存文件失败",
        })
        return
    }

    // 假设你有一个函数 uploadToNetwork 可以将本地文件上传到网络位置
    uploadedURL, err := uploadToNetwork(savePath, networkLocation)
    if err!= nil {
        c.JSON(500, gin.H{
            "message": "上传文件到网络位置失败",
        })
        return
    }

    // 返回上传成功的 URL
    c.JSON(200, gin.H{
        "message": "文件上传成功",
        "url":     uploadedURL,
    })
}

func uploadToNetwork(localPath, networkLocation string) (string, error) {
    file, err := os.Open(localPath)
    if err!= nil {
		fmt.Println(1)
        return "", err
    }
    defer file.Close()

    body := &bytes.Buffer{}
    writer := multipart.NewWriter(body)
    part, err := writer.CreateFormFile("file", filepath.Base(localPath))
    if err!= nil {
		fmt.Println(2)
        return "", err
    }
    _, err = io.Copy(part, file)
    if err!= nil {
		fmt.Println(3)
        return "", err
    }
    writer.Close()

    req, err := http.NewRequest("POST", networkLocation, body)
    if err!= nil {
		fmt.Println(4)
        return "", err
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())
    // 使用 API 令牌进行认证，而不是 Cookie
    apiToken := "8c86a96e9c6aff896a6fb2ab242141ab"
    req.Header.Set("Authorization", "token "+ apiToken)
	//fmt.Println(req)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err!= nil {
		fmt.Println(5)
        return "", err
    }
    defer resp.Body.Close()

    // 处理 Gitee API 的错误响应
    if resp.StatusCode >= 400 {
		fmt.Println(6)
        return "", fmt.Errorf("Gitee API 响应错误: %s", resp.Status)
    }

    var result map[string]string
    if err := json.NewDecoder(resp.Body).Decode(&result); err!= nil {
		fmt.Println(7)
        return "", err
    }
    return result["url"], nil
}