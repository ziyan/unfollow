package ratelimit

import (
    "appengine"
    "appengine/memcache"
    "errors"
    "fmt"
    "strconv"
    "time"
)

var (
    ErrRateLimitExceeded = errors.New("ratelimit: exceeded rate limit")
)

func key(partition string, timestamp uint64) string {
    return fmt.Sprintf("ratelimit:%s:%x", partition, timestamp)
}

// Interval is in minutes
func Limit(context appengine.Context, partition string, interval uint8, limit uint64) (int64, error) {
    // get the minute timestamp
    timestamp := uint64(time.Now().Unix() / 60)

    // increment key
    value, err := memcache.Increment(context, key(partition, timestamp), 1, 0)
    if err != nil {
        return 0, err
    }

    // set expiration if we just created the key
    if value == 1 {
        if err := memcache.Set(context, &memcache.Item{
            Key:        "ratelimit:" + partition,
            Value:      []byte{'1'},
            Expiration: time.Duration(interval) * time.Minute,
        }); err != nil {
            return 0, err
        }
    }

    total := value

    // get previous values within the interval
    if interval > 1 {

        keys := make([]string, 0, interval-1)
        for i := uint8(1); i < interval; i++ {
            timestamp -= 1
            keys = append(keys, key(partition, timestamp))
        }

        items, err := memcache.GetMulti(context, keys)
        if err != nil {
            return 0, err
        }

        for _, item := range items {
            value, err := strconv.ParseUint(string(item.Value), 10, 64)
            if err != nil {
                return 0, err
            }
            total += value
        }
    }

    if total >= limit {
        return 0, ErrRateLimitExceeded
    }

    remaining := int64(limit) - int64(total)
    return remaining, nil
}
