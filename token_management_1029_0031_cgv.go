// 代码生成时间: 2025-10-29 00:31:26
package main

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
    "log"
)

// TokenManagement 结构体，用于治理代币系统
type TokenManagement struct {
    // 代币余额映射
    Balances map[string]float64
}

// NewTokenManagement 初始化 TokenManagement 结构体
func NewTokenManagement() *TokenManagement {
    return &TokenManagement{
        Balances: make(map[string]float64),
    }
}

// AddBalance 为指定账户增加代币余额
func (tm *TokenManagement) AddBalance(account string, amount float64) error {
    if amount < 0 {
        return fmt.Errorf("amount must be positive")
    }
    tm.Balances[account] += amount
    return nil
}

// SubtractBalance 从指定账户扣除代币余额
func (tm *TokenManagement) SubtractBalance(account string, amount float64) error {
    if amount < 0 {
        return fmt.Errorf("amount must be positive")
    }
    if tm.Balances[account] < amount {
        return fmt.Errorf("insufficient balance")
    }
    tm.Balances[account] -= amount
    return nil
}

// GetBalance 获取指定账户的代币余额
func (tm *TokenManagement) GetBalance(account string) (float64, error) {
    balance, exists := tm.Balances[account]
    if !exists {
        return 0, fmt.Errorf("account not found")
    }
    return balance, nil
}

// StartServer 启动治理代币系统的Fiber服务器
func StartServer(tm *TokenManagement) {
    app := fiber.New()

    // 添加代币余额
    app.Post("/add", func(c *fiber.Ctx) error {
        account := c.Query("account")
        amount := c.Query("amount")
        amountFloat, err := strconv.ParseFloat(amount, 64)
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "invalid amount format",
            })
        }
        if err := tm.AddBalance(account, amountFloat); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        return c.SendStatus(fiber.StatusOK)
    })

    // 扣除代币余额
    app.Delete("/sub", func(c *fiber.Ctx) error {
        account := c.Query("account")
        amount := c.Query("amount")
        amountFloat, err := strconv.ParseFloat(amount, 64)
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": "invalid amount format",
            })
        }
        if err := tm.SubtractBalance(account, amountFloat); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        return c.SendStatus(fiber.StatusOK)
    })

    // 获取代币余额
    app.Get("/balance", func(c *fiber.Ctx) error {
        account := c.Query("account")
        balance, err := tm.GetBalance(account)
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        return c.JSON(fiber.Map{
            "account": account,
            "balance": balance,
        })
    })

    log.Fatal(app.Listen(":3000"))
}

func main() {
    tm := NewTokenManagement()
    StartServer(tm)
}