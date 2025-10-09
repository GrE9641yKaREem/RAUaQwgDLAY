// 代码生成时间: 2025-10-10 02:42:12
package main

import (
    "os"
    "path/filepath"
    "log"
    "github.com/gofiber/fiber/v2"
)

// 定义一个全局变量，用于存储文件操作的相关信息
var filesPath string

// 定义一个函数，用于列出目录中的所有文件
func listFiles(dir string) ([]string, error) {
    var files []string
    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            files = append(files, path)
        }
        return nil
    })
    return files, err
}

// 定义一个函数，用于执行批量文件操作
func batchFileOperation(files []string, operation string) error {
    for _, file := range files {
        // 根据操作类型执行不同的操作
        switch operation {
        case "delete":
            if err := os.Remove(file); err != nil {
                return err
            }
        default:
            return ErrUnsupportedOperation
        }
    }
    return nil
}

// 定义一个结构体，用于返回API响应
type ApiResponse struct {
    Message string `json:"message"`
}

// 定义一个错误，表示不支持的操作
var ErrUnsupportedOperation = errors.New("unsupported operation")

func main() {
    app := fiber.New()

    // 设置文件路径变量
    filesPath = "./files"

    // 定义一个API端点，用于执行批量文件删除操作
    app.Post("/delete", func(c *fiber.Ctx) error {
        files, err := listFiles(filesPath)
        if err != nil {
            log.Printf("Error listing files: %v", err)
            return c.Status(500).JSON(ApiResponse{Message: "Error listing files"})
        }
        err = batchFileOperation(files, "delete")
        if err != nil {
            log.Printf("Error deleting files: %v", err)
            return c.Status(500).JSON(ApiResponse{Message: "Error deleting files"})
        }
        return c.Status(200).JSON(ApiResponse{Message: "Files deleted successfully"})
    })

    // 启动Fiber服务器
    log.Fatal(app.Listen(":3000"))
}