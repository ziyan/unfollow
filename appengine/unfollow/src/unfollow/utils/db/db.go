package db

import (
    "appengine"
    "appengine/datastore"
    "strings"
    "time"
    "unfollow/utils/cache"
    "unfollow/utils/security"
)

const (
    MEMCACHE_PREFIX = "db"
)

type Entity interface {
    Key() *datastore.Key
    SetKey(*datastore.Key)

    Ok() bool
    SetOk(bool)
}

type Database struct {
    Context     appengine.Context
    Cache       *cache.Cache
    Transacting bool
}

func New(context appengine.Context, cache *cache.Cache) *Database {
    return &Database{
        Context:     context,
        Cache:       cache,
        Transacting: false,
    }
}

type Options struct {
    GenerateKey func(appengine.Context, *datastore.Key) *datastore.Key
}

var DefaultOptions = Options{
    GenerateKey: GenerateStringKey,
}

func GenerateStringKey(context appengine.Context, key *datastore.Key) *datastore.Key {
    return datastore.NewKey(context, key.Kind(), security.GenerateRandomID(), 0, key.Parent())
}

func GenerateIntegerKey(context appengine.Context, key *datastore.Key) *datastore.Key {
    return datastore.NewKey(context, key.Kind(), "", time.Now().UnixNano()/10, key.Parent())
}

func (db *Database) Key(kind, stringID string, intID int64, parent *datastore.Key) *datastore.Key {
    return datastore.NewKey(db.Context, kind, stringID, intID, parent)
}

func (db *Database) Transaction(transaction func(db *Database) error, xg bool) error {
    // if already in a transaction, run the function directly
    if db.Transacting {
        return transaction(db)
    }

    return datastore.RunInTransaction(db.Context, func(context appengine.Context) error {
        return transaction(&Database{
            Context:     context,
            Cache:       db.Cache,
            Transacting: true,
        })
    }, &datastore.TransactionOptions{
        XG: xg,
    })
}

func (db *Database) New(entity Entity, options *Options) error {
    entity.SetOk(false)

    if options == nil {
        options = &DefaultOptions
    }

    if err := db.Transaction(func(db *Database) error {

        if !entity.Key().Incomplete() {
            // caller is asking for a specific key
            if err := datastore.Get(db.Context, entity.Key(), entity); err == nil {
                return nil
            }
        } else {
            // caller wants a random key
            for {
                key := options.GenerateKey(db.Context, entity.Key())
                if key == nil {
                    panic("db: no key generated")
                }

                err := datastore.Get(db.Context, key, nil)
                if err == datastore.ErrNoSuchEntity {
                    entity.SetKey(key)
                    break
                }
                if err != nil {
                    return err
                }
            }
        }

        if _, err := datastore.Put(db.Context, entity.Key(), entity); err != nil {
            return err
        }

        return nil

    }, false); err != nil {
        return err
    }

    if err := db.Cache.Delete(strings.Join([]string{MEMCACHE_PREFIX, entity.Key().Encode()}, ":"), nil); err != nil {
        return err
    }

    entity.SetOk(true)
    return nil
}

func (db *Database) Get(entity Entity, options *Options) error {
    key := entity.Key()
    if key.Incomplete() {
        panic("db: key incomplete")
    }

    entity.SetOk(false)

    if !db.Transacting {
        exist, err := db.Cache.Get(strings.Join([]string{MEMCACHE_PREFIX, key.Encode()}, ":"), entity, nil)
        if exist {
            entity.SetKey(key)
            entity.SetOk(true)
            return nil
        }

        if err != nil {
            return err
        }
    }

    err := datastore.Get(db.Context, key, entity)
    if err == datastore.ErrNoSuchEntity {
        return nil
    }

    if err != nil {
        return err
    }

    if !db.Transacting {
        if err := db.Cache.Set(strings.Join([]string{MEMCACHE_PREFIX, key.Encode()}, ":"), entity, nil); err != nil {
            return err
        }
    }

    entity.SetKey(key)
    entity.SetOk(true)
    return nil
}

func (db *Database) GetMulti(entities []Entity, options *Options) error {

    for _, entity := range entities {
        if entity.Key().Incomplete() {
            panic("db: key incomplete")
        }
        entity.SetOk(false)
    }

    missings := make([]*datastore.Key, 0, len(entities))
    destinations := make([]Entity, 0, len(entities))

    if !db.Transacting {

        // remember the keys because cache will override them
        keys := make([]*datastore.Key, 0, len(entities))
        for _, entity := range entities {
            keys = append(keys, entity.Key())
        }

        // build cache keys
        k := make([]string, 0, len(entities))
        objects := make([]interface{}, 0, len(entities))
        for _, entity := range entities {
            k = append(k, strings.Join([]string{MEMCACHE_PREFIX, entity.Key().Encode()}, ":"))
            objects = append(objects, entity)
        }

        // get from cache
        exists, err := db.Cache.GetMulti(k, objects, nil)
        if err != nil {
            return err
        }

        // mark existing ones
        for i, exist := range exists {
            if exist {
                entities[i].SetKey(keys[i])
                entities[i].SetOk(true)
            } else {
                missings = append(missings, keys[i])
                destinations = append(destinations, entities[i])
            }
        }

    } else {
        for _, entity := range entities {
            missings = append(missings, entity.Key())
            destinations = append(destinations, entity)
        }
    }

    // nothing else to get?
    if len(missings) == 0 {
        return nil
    }

    caching := make(map[string]interface{})

    // get the rest from datastore
    err := datastore.GetMulti(db.Context, missings, destinations)
    if err != nil {
        // figure out what is missing
        errs, ok := err.(appengine.MultiError)
        if !ok {
            return err
        }

        for i, err := range errs {

            // entity exists?
            if err == nil {
                entity := destinations[i]
                if !db.Transacting {
                    k := strings.Join([]string{MEMCACHE_PREFIX, entity.Key().Encode()}, ":")
                    caching[k] = entity
                }
                entity.SetOk(true)
                continue
            }

            // entity does not exist
            if err == datastore.ErrNoSuchEntity {
                continue
            }

            // all other error is fatal
            return err
        }
    } else {
        for _, entity := range destinations {
            if !db.Transacting {
                k := strings.Join([]string{MEMCACHE_PREFIX, entity.Key().Encode()}, ":")
                caching[k] = entity
            }
            entity.SetOk(true)
        }
    }

    if len(caching) > 0 {
        // save in cache
        if err := db.Cache.SetMulti(caching, nil); err != nil {
            return err
        }
    }

    return nil
}

func (db *Database) Put(entity Entity, options *Options) error {
    key := entity.Key()
    if key.Incomplete() {
        panic("db: key incomplete")
    }

    entity.SetOk(false)

    if _, err := datastore.Put(db.Context, key, entity); err != nil {
        return err
    }

    if err := db.Cache.Delete(strings.Join([]string{MEMCACHE_PREFIX, key.Encode()}, ":"), nil); err != nil {
        return err
    }

    entity.SetOk(true)
    return nil
}

func (db *Database) PutMulti(entities []Entity, options *Options) error {
    keys := make([]*datastore.Key, 0, len(entities))
    values := make([]Entity, 0, len(entities))
    caches := make([]string, 0, len(keys))

    for _, entity := range entities {
        key := entity.Key()
        if key.Incomplete() {
            panic("db: key incomplete")
        }

        entity.SetOk(false)

        keys = append(keys, key)
        values = append(values, entity)
        caches = append(caches, strings.Join([]string{MEMCACHE_PREFIX, key.Encode()}, ":"))
    }

    if _, err := datastore.PutMulti(db.Context, keys, values); err != nil {
        if errs, ok := err.(appengine.MultiError); ok {
            for _, err := range errs {
                if err != nil {
                    return err
                }
            }
        }
        return err
    }

    if err := db.Cache.DeleteMulti(caches, nil); err != nil {
        return err
    }

    for _, entity := range entities {
        entity.SetOk(true)
    }
    return nil
}

func (db *Database) Delete(key *datastore.Key, options *Options) error {
    if key.Incomplete() {
        panic("db: key incomplete")
    }

    if err := datastore.Delete(db.Context, key); err != nil {
        return err
    }

    if err := db.Cache.Delete(strings.Join([]string{MEMCACHE_PREFIX, key.Encode()}, ":"), nil); err != nil {
        return err
    }

    return nil
}

func (db *Database) DeleteMulti(keys []*datastore.Key, options *Options) error {
    caches := make([]string, 0, len(keys))
    for _, key := range keys {
        if key.Incomplete() {
            panic("db: key incomplete")
        }

        caches = append(caches, strings.Join([]string{MEMCACHE_PREFIX, key.Encode()}, ":"))
    }

    if err := datastore.DeleteMulti(db.Context, keys); err != nil {
        if errs, ok := err.(appengine.MultiError); ok {
            for _, err := range errs {
                if err != nil {
                    return err
                }
            }
        }
        return err
    }

    if err := db.Cache.DeleteMulti(caches, nil); err != nil {
        return err
    }

    return nil
}
