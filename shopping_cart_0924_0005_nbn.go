// 代码生成时间: 2025-09-24 00:05:22
package main

import (
    "fmt"
    "github.com/gofiber/fiber/v2"
    "log"
)

// ShoppingCart represents a shopping cart with items.
type ShoppingCart struct {
    Items map[string]int
}

// AddItem adds an item to the shopping cart.
func (sc *ShoppingCart) AddItem(item string, quantity int) error {
    if quantity < 1 {
        return fmt.Errorf("quantity must be at least 1")
    }
    sc.Items[item] = sc.Items[item] + quantity
    return nil
}

// RemoveItem removes an item from the shopping cart.
func (sc *ShoppingCart) RemoveItem(item string) error {
    if _, exists := sc.Items[item]; !exists {
        return fmt.Errorf("item not found in the cart")
    }
    delete(sc.Items, item)
    return nil
}

// ClearCart clears all items from the shopping cart.
func (sc *ShoppingCart) ClearCart() {
    sc.Items = make(map[string]int)
}

// CartHandler handles cart-related requests.
func CartHandler(c *fiber.Ctx) error {
    cart := ShoppingCart{Items: make(map[string]int)}

    switch c.Method() {
    case fiber.MethodPost:
        // Handle adding an item to the cart
        item := c.Query("item")
        quantity, _ := c.QueryInt("quantity")
        if err := cart.AddItem(item, quantity); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        return c.JSON(cart.Items)

    case fiber.MethodDelete:
        // Handle removing an item from the cart
        item := c.Query("item")
        if err := cart.RemoveItem(item); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        return c.JSON(cart.Items)

    case fiber.MethodGet:
        // Handle retrieving the cart contents
        return c.JSON(cart.Items)

    default:
        return c.SendStatus(fiber.StatusMethodNotAllowed)
    }
}

func main() {
    app := fiber.New()
    app.Get("/cart", CartHandler)
    app.Post("/cart", CartHandler)
    app.Delete("/cart", CartHandler)

    log.Fatal(app.Listen(":3000"))
}
