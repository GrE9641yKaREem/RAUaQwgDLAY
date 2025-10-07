// 代码生成时间: 2025-10-08 03:52:24
package main

import (
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
    "encoding/hex"
    "errors"
    "fmt"
    "fiber/fiber/v2"
    "log"
    "strings"
)

// EncryptionKey is a 32-byte key for AES encryption.
// It should be kept secret and secure.
var EncryptionKey = []byte("your-encryption-key-32-bytes-long")

// PaddingPKCS7 pads the given []byte to match the AES block size (16 bytes).
func PaddingPKCS7(data []byte) []byte {
    blockSize := aes.BlockSize
    length := len(data)
    padding := blockSize - length%blockSize
    padText := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(data, padText...)
}

// UnPaddingPKCS7 removes the PKCS7 padding from the given []byte.
func UnPaddingPKCS7(data []byte) []byte {
    length := len(data)
    unpadding := int(data[length-1])
    return data[:(length - unpadding)]
}

// Encrypt encrypts the given plaintext using AES.
func Encrypt(plaintext string) (string, error) {
    data := PaddingPKCS7([]byte(plaintext))
    block, err := aes.NewCipher(EncryptionKey)
    if err != nil {
        return "", err
    }
    b := base64.StdEncoding.EncodeToString(pkcs5Padding(data, aes.BlockSize))
    mode := cipher.NewCBCEncrypter(block, EncryptionKey[:block.BlockSize()])
    ciphertext := make([]byte, len(data))
    mode.CryptBlocks(ciphertext, data)
    return b, nil
}

// Decrypt decrypts the given ciphertext using AES.
func Decrypt(ciphertext string) (string, error) {
    decodeBytes, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }
    block, err := aes.NewCipher(EncryptionKey)
    if err != nil {
        return "", err
    }
    if len(decodeBytes)%aes.BlockSize != 0 {
        return "", errors.New("ciphertext is not a multiple of the block size")
    }
    mode := cipher.NewCBCDecrypter(block, EncryptionKey[:block.BlockSize()])
    decrypted := make([]byte, len(decodeBytes))
    mode.CryptBlocks(decrypted, decodeBytes)
    decrypted = UnPaddingPKCS7(decrypted)
    return string(decrypted), nil
}

func main() {
    app := fiber.New()

    app.Post("/encrypt", func(c *fiber.Ctx) error {
        plaintext := c.Get("plaintext")
        encrypted, err := Encrypt(plaintext)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": fmt.Sprintf("Encryption failed: %s", err),
            })
        }
        return c.Status(fiber.StatusOK).JSON(fiber.Map{
            "plaintext": plaintext,
            "encrypted": encrypted,
        })
    })

    app.Post("/decrypt", func(c *fiber.Ctx) error {
        ciphertext := c.Get("ciphertext")
        decrypted, err := Decrypt(ciphertext)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": fmt.Sprintf("Decryption failed: %s", err),
            })
        }
        return c.Status(fiber.StatusOK).JSON(fiber.Map{
            "ciphertext": ciphertext,
            "decrypted": decrypted,
        })
    })

    log.Fatal(app.Listen(":3000"))
}