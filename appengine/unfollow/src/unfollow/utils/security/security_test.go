package security

import "testing"

func TestGenerateRandom(t *testing.T) {
    if data := GenerateRandom(10); len(data) != 10 {
        t.Fatal()
    }
}

func TestGenerateRandomHexString(t *testing.T) {
    if data := GenerateRandomHexString(10); len(data) != 20 {
        t.Fatal()
    }
}

func TestGenerateRandomBase64String(t *testing.T) {
    GenerateRandomBase64String(1)
}

func TestGenerateRandomID(t *testing.T) {
    if id := GenerateRandomID(); len(id) != 32 {
        t.Fatal()
    }
}
