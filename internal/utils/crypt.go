package utils

import (
    "crypto/rand"
    "crypto/subtle"
    "encoding/base64"
    "errors"
    "fmt"
    "strings"

    "golang.org/x/crypto/argon2"
)

type Argon2Config struct {
    Time    uint32
    Memory  uint32
    Threads uint8
    KeyLen  uint32
    SaltLen uint32
}

var DefaultConfig = &Argon2Config{
    Time:    3,
    Memory:  64 * 1024,
    Threads: 4,
    KeyLen:  32,
    SaltLen: 16,
}

func Hash(password string) (string, error) {
    return HashWithConfig(password, DefaultConfig)
}

func HashWithConfig(password string, cfg *Argon2Config) (string, error) {
    salt := make([]byte, cfg.SaltLen)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }

    hash := argon2.IDKey([]byte(password), salt, cfg.Time, cfg.Memory, cfg.Threads, cfg.KeyLen)

    b64Salt := base64.RawStdEncoding.EncodeToString(salt)
    b64Hash := base64.RawStdEncoding.EncodeToString(hash)

    return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
        argon2.Version, cfg.Memory, cfg.Time, cfg.Threads, b64Salt, b64Hash), nil
}

func Verify(password, encodedHash string) (bool, error) {
    parts := strings.Split(encodedHash, "$")
    if len(parts) != 6 {
        return false, errors.New("invalid hash format")
    }

    var version int
    var memory, time uint32
    var threads uint8

    if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
        return false, err
    }
    if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
        return false, err
    }

    salt, err := base64.RawStdEncoding.DecodeString(parts[4])
    if err != nil {
        return false, err
    }

    hash, err := base64.RawStdEncoding.DecodeString(parts[5])
    if err != nil {
        return false, err
    }

    newHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(hash)))

    return subtle.ConstantTimeCompare(hash, newHash) == 1, nil
}