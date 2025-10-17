// 代码生成时间: 2025-10-17 21:15:39
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "log"
    "github.com/gofiber/fiber/v2"
)

// DiskSpaceManager provides functionality to manage disk space
type DiskSpaceManager struct {
    // Holds the path to the disk or directory to manage
    Path string
}

// NewDiskSpaceManager creates a new instance of DiskSpaceManager
func NewDiskSpaceManager(path string) *DiskSpaceManager {
    return &DiskSpaceManager{Path: path}
}

// CheckSpaceAvailability checks if there is enough disk space available
func (d *DiskSpaceManager) CheckSpaceAvailability() (bool, error) {
    // Get filesystem stats for the given path
    fs := &fiber.Disk{} // Fiber's Disk handler
    stats, err := fs.Stats()
    if err != nil {
        return false, err
    }

    // Check if the available space is greater than 0
    if stats.Free > 0 {
        return true, nil
    }
    return false, nil
}

// GetDiskUsage retrieves disk usage statistics
func (d *DiskSpaceManager) GetDiskUsage() (*fiber.DiskStats, error) {
    fs := &fiber.Disk{} // Fiber's Disk handler
    stats, err := fs.Stats()
    if err != nil {
        return nil, err
    }
    return stats, nil
}

// GetDirectorySize calculates the total size of a directory
func (d *DiskSpaceManager) GetDirectorySize() (int64, error) {
    var size int64
    err := filepath.WalkDir(d.Path, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if !d.IsDir() {
            fi, err := d.Info()
            if err != nil {
                return err
            }
            size += fi.Size()
        }
        return nil
    })
    if err != nil {
        return 0, err
    }
    return size, nil
}

// StartServer starts the Fiber HTTP server with disk space management routes
func StartServer(manager *DiskSpaceManager) error {
    app := fiber.New()

    // Route to check if there is enough disk space available
    app.Get("/check", func(c *fiber.Ctx) error {
        available, err := manager.CheckSpaceAvailability()
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        return c.JSON(fiber.Map{
            "available": available,
        })
    })

    // Route to get disk usage statistics
    app.Get("/stats", func(c *fiber.Ctx) error {
        stats, err := manager.GetDiskUsage()
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        return c.JSON(fiber.Map{
            "stats": stats,
        })
    })

    // Route to get the total size of a directory
    app.Get("/size", func(c *fiber.Ctx) error {
        size, err := manager.GetDirectorySize()
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        return c.JSON(fiber.Map{
            "size": size,
        })
    })

    // Start the Fiber server
    return app.Listen(":3000")
}

func main() {
    manager := NewDiskSpaceManager(".")
    if err := StartServer(manager); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}