package security

import (
    "bytes"
    "crypto/rand"
    "encoding/base64"
    "encoding/hex"
    "io"
)

const (
    UpperAlpha   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
    LowerAlpha   = "abcdefghijklmnopqrstuvwxyz"
    Digits       = "0123456789"
    Alpha        = UpperAlpha + LowerAlpha
    AlphaNumeric = Alpha + Digits
)

// Generate a random binary string.
func GenerateRandom(length int) []byte {
    data := make([]byte, length)
    if _, err := io.ReadFull(rand.Reader, data); err != nil {
        panic(err)
    }
    return data
}

// Generate a random hex encoded string.
func GenerateRandomHexString(length int) string {
    return hex.EncodeToString(GenerateRandom(length))
}

// Generate a random base64 encoded string.
func GenerateRandomBase64String(length int) string {
    return base64.URLEncoding.EncodeToString(GenerateRandom(length))
}

// Generate random string from given alphabet.
func GenerateRandomString(length int, alphabet string) string {
    choices := []byte(alphabet)
    size := len(choices)
    buffer := new(bytes.Buffer)

    for _, data := range GenerateRandom(length) {
        buffer.WriteByte(choices[int(data)%size])
    }

    return buffer.String()
}

// Generate a random alpha numeric string of length 32.
func GenerateRandomID() string {
    return GenerateRandomString(32, AlphaNumeric)
}
