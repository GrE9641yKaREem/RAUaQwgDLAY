// 代码生成时间: 2025-09-23 12:21:07
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "log"
    "io"
    "io/ioutil"
    "github.com/gofiber/fiber/v2"
)

// BackupSyncService 结构体，用于文件备份和同步操作
type BackupSyncService struct {
    srcPath  string
    destPath string
}

// NewBackupSyncService 创建一个新的BackupSyncService实例
func NewBackupSyncService(srcPath, destPath string) *BackupSyncService {
    return &BackupSyncService{
        srcPath:  srcPath,
        destPath: destPath,
    }
}

// SyncFileSync 同步单个文件
func (s *BackupSyncService) SyncFileSync(src, dest string) error {
    srcInfo, err := os.Stat(src)
    if err != nil {
        return fmt.Errorf("failed to stat source file: %w", err)
    }
    if !srcInfo.Mode().IsRegular() {
        return fmt.Errorf("%s is not a regular file", src)
    }

    err = copyFile(src, dest)
    if err != nil {
        return fmt.Errorf("failed to copy file: %w", err)
    }
    return nil
}

// SyncDir syncs the contents of the source directory to the destination directory
func (s *BackupSyncService) SyncDir() error {
    return filepath.WalkDir(s.srcPath, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return fmt.Errorf("failed to walk dir: %w", err)
        }

        // Skip if the directory entry is the destination directory itself
        if d.Name() == filepath.Base(s.destPath) {
            return filepath.SkipDir
        }

        destPath := filepath.Join(s.destPath, strings.TrimPrefix(path, s.srcPath+string(filepath.Separator)))
        return s.syncPath(path, destPath)
    })
}

// syncPath 同步单个路径，无论是文件还是目录
func (s *BackupSyncService) syncPath(src, dest string) error {
    srcInfo, err := os.Stat(src)
    if err != nil {
        return fmt.Errorf("failed to stat source path: %w", err)
    }

    if srcInfo.IsDir() {
        // Create destination directory if it does not exist
        if _, err := os.Stat(dest); os.IsNotExist(err) {
            if err := os.MkdirAll(dest, 0755); err != nil {
                return fmt.Errorf("failed to create destination directory: %w", err)
            }
        }
    } else {
        // Ensure destination directory exists
        destDir := filepath.Dir(dest)
        if _, err := os.Stat(destDir); os.IsNotExist(err) {
            if err := os.MkdirAll(destDir, 0755); err != nil {
                return fmt.Errorf("failed to create destination directory: %w", err)
            }
        }

        // Synchronize the file
        return s.SyncFileSync(src, dest)
    }
    return nil
}

// copyFile copies a file from src to dest
func copyFile(src, dest string) error {
    sourceFileStat, err := os.Stat(src)
    if err != nil {
        return err
    }

    if !sourceFileStat.Mode().IsRegular() {
        return fmt.Errorf("%s is not a regular file", src)
    }

    source, err := os.Open(src)
    if err != nil {
        return err
    }
    defer source.Close()

    destination, err := os.Create(dest)
    if err != nil {
        return err
    }
    defer destination.Close()

    numberOfBytesCopied, err := io.Copy(destination, source)
    if err != nil {
        return err
    }

    if numberOfBytesCopied != sourceFileStat.Size() {
        return fmt.Errorf("failed to copy file completely")
    }

    return destination.Sync()
}

// main function to run the backup and sync tool
func main() {
    service := NewBackupSyncService("./src", "./dest")
    err := service.SyncDir()
    if err != nil {
        log.Fatalf("backup and sync failed: %v", err)
    }
    fmt.Println("Backup and sync completed successfully")

    // Set up Fiber web server to handle requests
    app := fiber.New()

    // Define a route to trigger backup and sync
    app.Get("/backup", func(c *fiber.Ctx) error {
        if err := service.SyncDir(); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": err.Error(),
            })
        }
        return c.SendString("Backup and sync completed successfully")
    })

    // Start the server
    log.Fatal(app.Listen(":3000"))
}
