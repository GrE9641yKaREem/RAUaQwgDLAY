// 代码生成时间: 2025-10-30 11:54:55
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
    "github.com/gofiber/fiber/v2"
)

// ABTest represents the structure for A/B testing scenarios
type ABTest struct {
    VariantA bool
    VariantB bool
}

// NewABTest creates a new ABTest instance
func NewABTest() *ABTest {
    return &ABTest{}
}

// RunABTest starts the A/B testing
func (ab *ABTest) RunABTest(c *fiber.Ctx) error {
    // Randomly assign users to either variant A or B
    random := time.Now().UnixNano() % 2
    if random == 0 {
        ab.VariantA = true
        ab.VariantB = false
    } else {
        ab.VariantA = false
        ab.VariantB = true
    }
    
    // Return the variant assigned to the user
    return c.JSON(fiber.Map{
        "variant": ab.VariantA,
    })
}

func main() {
    app := fiber.New()
    
    // Endpoint to start A/B testing
    app.Get("/start-ab-test", func(c *fiber.Ctx) error {
        abTest := NewABTest()
        return abTest.RunABTest(c)
    })
    
    // Start the Fiber server
    log.Println("Server started on :3000")
    if err := app.Listen(":3000"); err != nil {
        log.Fatal(err)
    }
}
