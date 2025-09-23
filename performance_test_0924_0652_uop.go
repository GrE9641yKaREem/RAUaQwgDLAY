// 代码生成时间: 2025-09-24 06:52:47
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/gofiber/fiber/v2"
)

// setupFiberApp sets up the Fiber application with a simple route.
func setupFiberApp() *fiber.App {
    app := fiber.New()
    app.Get("/test", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })
    return app
}

// runPerformanceTest runs the performance test by simulating concurrent requests.
func runPerformanceTest(app *fiber.App, duration time.Duration) {
    start := time.Now()
    var wg sync.WaitGroup
    const numConcurrentRequests = 100
    
    // Start the Fiber server in a separate goroutine.
    go func() {
        if err := app.Listen(":3000"); err != nil {
            log.Fatalf("Error starting Fiber server: %v", err)
        }
    }()

    // Simulate concurrent requests.
    for i := 0; i < numConcurrentRequests; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            client := &http.Client{Timeout: 10 * time.Second}
            for time.Since(start) < duration {
                _, err := client.Get("http://localhost:3000/test")
                if err != nil {
                    log.Printf("Error making request: %v", err)
                }
            }
        }()
    }
    wg.Wait()
    fmt.Printf("Performance test completed in %v
", time.Since(start))
}

func main() {
    app := setupFiberApp()
    duration := 10 * time.Second // Test duration
    runPerformanceTest(app, duration)
}
