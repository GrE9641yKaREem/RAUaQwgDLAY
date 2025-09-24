// 代码生成时间: 2025-09-24 20:12:15
package main

import (
    "os/exec"
    "os"
    "fmt"
    "log"
    "strings"

    "github.com/gofiber/fiber/v2"
)

// ProcessManager struct to hold process information
type ProcessManager struct {
    process *exec.Cmd
}

// NewProcessManager creates a new instance of ProcessManager
func NewProcessManager() *ProcessManager {
    return &ProcessManager{}
}

// StartProcess starts the process
func (pm *ProcessManager) StartProcess(name string, arg ...string) error {
    // Construct the command
    cmd := exec.Command(name, arg...)
    pm.process = cmd

    // Start the process
    if err := cmd.Start(); err != nil {
        return err
    }

    fmt.Printf("Process %s started with PID: %d
", name, cmd.Process.Pid)
    return nil
}

// StopProcess stops the process
func (pm *ProcessManager) StopProcess() error {
    if pm.process == nil {
        return fmt.Errorf("no process is running")
    }

    // Stop the process
    if err := pm.process.Process.Kill(); err != nil {
        return err
    }

    fmt.Println("Process stopped")
    return nil
}

// StatusProcess checks the status of the process
func (pm *ProcessManager) StatusProcess() (string, error) {
    if pm.process == nil {
        return "", fmt.Errorf("no process is running")
    }

    // Check the status of the process
    status, err := pm.process.Process.Wait()
    if err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok {
            if status, ok := exitErr.Sys().(*os.ProcessState); ok {
                if status.Exited() {
                    return fmt.Sprintf("Process exited with status: %d", status.ExitCode()), nil
                }
            }
        }
        return "", err
    }
    return "Process is running", nil
}

func main() {
    app := fiber.New()
    pm := NewProcessManager()

    // Start a process
    app.Get("/start", func(c *fiber.Ctx) error {
        name := c.Query("name")
        arg := c.Query("arg")
        if name == "" {
            return c.Status(fiber.StatusBadRequest).SendString("Process name is required")
        }

        if err := pm.StartProcess(name, arg); err != nil {
            return c.Status(fiber.StatusInternalServerError).SendString("Failed to start process: " + err.Error())
        }
        return c.SendString("Process started successfully")
    })

    // Stop a process
    app.Get("/stop", func(c *fiber.Ctx) error {
        if err := pm.StopProcess(); err != nil {
            return c.Status(fiber.StatusInternalServerError).SendString("Failed to stop process: " + err.Error())
        }
        return c.SendString("Process stopped successfully")
    })

    // Get process status
    app.Get("/status", func(c *fiber.Ctx) error {
        status, err := pm.StatusProcess()
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).SendString("Failed to get status: " + err.Error())
        }
        return c.SendString(status)
    })

    // Log requests to stdout
    app.Use("/", func(c *fiber.Ctx) error {
        log.Printf("%s %s
", c.Method(), c.Path())
        return c.Next()
    })

    // Start the Fiber server
    log.Fatal(app.Listen(":3000"))
}
