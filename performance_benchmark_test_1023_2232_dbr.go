// 代码生成时间: 2025-10-23 22:32:49
package main

import (
    "fmt"
    "net/http"
    "time"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/benchmark"
)

// Initialize Fiber app
func main() {
    app := fiber.New()

    // Register benchmark middleware
# TODO: 优化性能
    app.Use(benchmark.New())

    // Define a simple route
    app.Get("/test", func(c *fiber.Ctx) error {
        // Simulate some work
        time.Sleep(100 * time.Millisecond)
        return c.SendString("Hello, World!")
    })

    // Start the server
    addr := ":3000"
    fmt.Println("Server is running on:", addr)
    if err := app.Listen(addr); err != nil && err != fiber.ErrServerClosed {
        fmt.Println("Error starting server:", err)
    }
}
# 添加错误处理
