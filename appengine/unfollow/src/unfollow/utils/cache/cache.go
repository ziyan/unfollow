package cache

import (
    "appengine"
    "appengine/memcache"
    "unfollow/settings"
    "time"
)

type Cache struct {
    Context appengine.Context
    Cache   map[string][]byte
    Debug   bool
}

func New(context appengine.Context) *Cache {
    return &Cache{
        Context: context,
        Cache:   make(map[string][]byte),
        Debug:   settings.DEBUG,
    }
}

type Options struct {
    Expiration time.Duration
}

var DefaultOptions = Options{}

func (cache *Cache) Size() int {
    return len(cache.Cache)
}

func (cache *Cache) Get(key string, object interface{}, options *Options) (bool, error) {

    if cache.Debug {
        cache.Context.Infof("cache: get: %v", key)
    }

    // check in memory cache
    if data, ok := cache.Cache[key]; ok {

        // negative cache
        if data == nil {
            return false, nil
        }

        // unmarshal
        if object != nil {
            if err := memcache.Gob.Unmarshal(data, object); err != nil {
                return false, err
            }
        }

        return true, nil
    }

    // check with memcache
    item, err := memcache.Get(cache.Context, key)

    // cache miss
    if err == memcache.ErrCacheMiss {
        cache.Cache[key] = nil
        return false, nil
    }

    if err != nil {
        return false, err
    }

    if object != nil {
        if err := memcache.Gob.Unmarshal(item.Value, object); err != nil {
            return false, err
        }
    }

    cache.Cache[key] = item.Value
    return true, nil
}

func (cache *Cache) GetMulti(keys []string, objects []interface{}, options *Options) ([]bool, error) {
    if objects != nil && len(keys) != len(objects) {
        panic("cache: keys and objects length must match")
    }

    exists := make([]bool, len(keys))
    missings := make(map[int]bool, 0)

    // first we look at in-memory cache
    for i := 0; i < len(keys); i++ {
        key := keys[i]

        // debug
        if cache.Debug {
            cache.Context.Infof("cache: get multi: %v", key)
        }

        data, ok := cache.Cache[key]

        // if item don't exist, add it to the list of keys to get
        if !ok {
            missings[i] = true
            continue
        }

        // if it exists, unmarshal it
        if data != nil {
            if objects != nil && objects[i] != nil {
                if err := memcache.Gob.Unmarshal(data, objects[i]); err != nil {
                    return nil, err
                }
            }
            exists[i] = true
        }
    }

    // everything is there, nothing else to get?
    if len(missings) == 0 {
        return exists, nil
    }

    missed := make([]string, 0, len(missings))
    for i, _ := range missings {
        missed = append(missed, keys[i])
    }

    // get from memcache
    items, err := memcache.GetMulti(cache.Context, missed)
    if err != nil {
        errs, ok := err.(appengine.MultiError)
        if !ok {
            return nil, err
        }
        for _, err := range errs {
            if err != nil {
                return nil, err
            }
        }
    }

    for i, _ := range missings {
        key := keys[i]
        item, ok := items[key]
        if !ok {
            cache.Cache[key] = nil
            continue
        }

        if objects != nil && objects[i] != nil {
            if err := memcache.Gob.Unmarshal(item.Value, objects[i]); err != nil {
                return nil, err
            }
        }

        // report it as exists
        exists[i] = true

        // cache it in memory
        cache.Cache[key] = item.Value
    }

    return exists, nil
}

func (cache *Cache) Set(key string, object interface{}, options *Options) error {

    if cache.Debug {
        cache.Context.Infof("cache: set: %v = %v", key, object)
    }

    value, err := memcache.Gob.Marshal(object)
    if err != nil {
        return err
    }

    item := &memcache.Item{
        Key:   key,
        Value: value,
    }

    if options != nil {
        item.Expiration = options.Expiration
    }

    if err := memcache.Set(cache.Context, item); err != nil {
        return err
    }

    cache.Cache[key] = value
    return nil
}

func (cache *Cache) SetMulti(objects map[string]interface{}, options *Options) error {

    if cache.Debug {
        cache.Context.Infof("cache: set multi: %v", objects)
    }

    items := make([]*memcache.Item, 0, len(objects))

    for key, object := range objects {
        value, err := memcache.Gob.Marshal(object)
        if err != nil {
            return err
        }

        item := &memcache.Item{
            Key:   key,
            Value: value,
        }

        if options != nil {
            item.Expiration = options.Expiration
        }

        items = append(items, item)
    }

    err := memcache.SetMulti(cache.Context, items)
    if err != nil {
        errs, ok := err.(appengine.MultiError)
        if !ok {
            return err
        }
        for _, err := range errs {
            if err != nil {
                return err
            }
        }
    }

    for _, item := range items {
        cache.Cache[item.Key] = item.Value
    }
    return nil
}

func (cache *Cache) Add(key string, object interface{}, options *Options) (bool, error) {

    if cache.Debug {
        cache.Context.Infof("cache: add: %v = %v", key, object)
    }

    value, err := memcache.Gob.Marshal(object)
    if err != nil {
        return false, err
    }

    item := &memcache.Item{
        Key:   key,
        Value: value,
    }

    if options != nil {
        item.Expiration = options.Expiration
    }

    if err := memcache.Add(cache.Context, item); err != nil {
        if err == memcache.ErrNotStored {
            return false, nil
        }
        return false, err
    }
    cache.Cache[key] = value
    return true, nil
}

func (cache *Cache) Delete(key string, options *Options) error {
    if cache.Debug {
        cache.Context.Infof("cache: delete: %v", key)
    }

    if err := memcache.Delete(cache.Context, key); err != nil && err != memcache.ErrCacheMiss {
        return err
    }
    cache.Cache[key] = nil
    return nil
}

func (cache *Cache) DeleteMulti(keys []string, options *Options) error {
    if cache.Debug {
        cache.Context.Infof("cache: delete multi: %v", keys)
    }

    err := memcache.DeleteMulti(cache.Context, keys)
    if err != nil {
        errs, ok := err.(appengine.MultiError)
        if !ok {
            return err
        }
        for _, err := range errs {
            if err != nil && err != memcache.ErrCacheMiss {
                return err
            }
        }
    }

    for _, key := range keys {
        cache.Cache[key] = nil
    }
    return nil
}

func (cache *Cache) Clear() {
    for key := range cache.Cache {
        if cache.Debug {
            cache.Context.Infof("cache: clear: %v", key)
        }
        delete(cache.Cache, key)
    }
}
