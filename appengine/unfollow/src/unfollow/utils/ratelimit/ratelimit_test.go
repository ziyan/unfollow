package ratelimit

import (
    "appengine/aetest"
    "testing"
)

func TestLimit(t *testing.T) {
    context, err := aetest.NewContext(nil)
    if err != nil {
        t.Fatal(err)
    }
    defer context.Close()

    limit := 10
    for i := 1; i < limit; i++ {
        remaining, err := Limit(context, "test", 5, uint64(limit))
        if err != nil {
            t.Fatal(err)
        }

        if remaining != int64(limit-i) {
            t.Fatal("ratelimit: remaining incorrect", remaining)
        }
    }

    remaining, err := Limit(context, "test", 5, uint64(limit))
    if err != ErrRateLimitExceeded {
        t.Fatal("ratelimit: rate limit should have been exceeded", err)
    }

    if remaining != 0 {
        t.Fatal("ratelimit: remaining should be zero", remaining)
    }
}
