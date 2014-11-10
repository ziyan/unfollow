package cache

import (
    "appengine/aetest"
    "appengine/memcache"
    "testing"
)

var values = map[string]interface{}{
    "key":  "Hello world!",
    "key1": 1234,
}

func TestSingle(t *testing.T) {
    context, err := aetest.NewContext(nil)
    if err != nil {
        t.Fatal(err)
    }
    defer context.Close()

    memcache.Flush(context)
    defer memcache.Flush(context)

    cache := New(context)
    defer cache.Clear()

    for key, value := range values {
        object := value

        exist, err := cache.Get(key, &object, nil)
        if err != nil {
            t.Fatal(err)
        }

        if exist {
            t.Fatal("cache: should be missing")
        }

        if err := cache.Set(key, &object, nil); err != nil {
            t.Fatal(err)
        }

        exist, err = cache.Get(key, &object, nil)
        if err != nil {
            t.Fatal(err)
        }

        if !exist {
            t.Fatal("cache: should exist")
        }

        if object != value {
            t.Fatal("cache: mismatch", object)
        }

        cache.Clear()

        exist, err = cache.Get(key, &object, nil)
        if err != nil {
            t.Fatal(err)
        }

        if !exist {
            t.Fatal("cache: should exist")
        }

        if object != value {
            t.Fatal("cache: mismatch", object)
        }

        if err := cache.Delete(key, nil); err != nil {
            t.Fatal(err)
        }

        exist, err = cache.Get(key, &object, nil)
        if err != nil {
            t.Fatal(err)
        }

        if exist {
            t.Fatal("cache: should be missing")
        }

        added, err := cache.Add(key, &object, nil)
        if err != nil {
            t.Fatal(err)
        }

        if !added {
            t.Fatal("cache: should be added")
        }

        added, err = cache.Add(key, &object, nil)
        if err != nil {
            t.Fatal(err)
        }

        if added {
            t.Fatal("cache: should not be added")
        }
    }
}

func TestMulti(t *testing.T) {
    context, err := aetest.NewContext(nil)
    if err != nil {
        t.Fatal(err)
    }
    defer context.Close()

    memcache.Flush(context)
    defer memcache.Flush(context)

    cache := New(context)
    defer cache.Clear()

    mapping := make(map[string]interface{}, len(values))
    for key, value := range values {
        object := value
        mapping[key] = &object
    }

    if err := cache.SetMulti(mapping, nil); err != nil {
        t.Fatal(err)
    }

    keys := make([]string, 0, len(mapping))
    objects := make([]interface{}, 0, len(mapping))
    for key, object := range mapping {
        keys = append(keys, key)
        objects = append(objects, object)
    }

    if _, err := cache.GetMulti(keys, objects, nil); err != nil {
        t.Fatal(err)
    }

    cache.Clear()

    if _, err := cache.GetMulti(keys, objects, nil); err != nil {
        t.Fatal(err)
    }

    if err := cache.DeleteMulti(keys, nil); err != nil {
        t.Fatal(err)
    }
}

func BenchmarkSingle(b *testing.B) {
    b.StopTimer()

    context, err := aetest.NewContext(nil)
    if err != nil {
        b.Fatal(err)
    }
    defer context.Close()

    memcache.Flush(context)
    defer memcache.Flush(context)

    cache := New(context)
    defer cache.Clear()

    key := "test"
    object := "Hello world!"

    if err := cache.Set(key, &object, nil); err != nil {
        b.Fatal(err)
    }

    b.StartTimer()

    for i := 0; i < b.N; i++ {
        if _, err := cache.Get(key, &object, nil); err != nil {
            b.Fatal(err)
        }
    }

    b.StopTimer()
}
