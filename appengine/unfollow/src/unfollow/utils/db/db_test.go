package db

import (
    "appengine/aetest"
    "appengine/datastore"
    "encoding/gob"
    "testing"
    "unfollow/utils/cache"
)

type TestModel struct {
    key     *datastore.Key `datastore:"-"`
    ok      bool           `datastore:"-"`
    Content string         `datastore:"content,noindex"`
}

func (model *TestModel) Key() *datastore.Key {
    return model.key
}

func (model *TestModel) SetKey(key *datastore.Key) {
    model.key = key
}

func (model *TestModel) Ok() bool {
    return model.ok
}

func (model *TestModel) SetOk(ok bool) {
    model.ok = ok
}

func init() {
    gob.Register(TestModel{})
}

func TestSingle(t *testing.T) {
    context, err := aetest.NewContext(nil)
    if err != nil {
        t.Fatal(err)
    }
    defer context.Close()

    cache := cache.New(context)
    defer cache.Clear()

    db := New(context, cache)

    key := datastore.NewKey(context, "test", "", 0, nil)
    entity := TestModel{key, false, "Test"}
    err = db.New(&entity, nil)
    if err != nil {
        t.Fatal(err)
    }
    key = entity.Key()

    entity1 := TestModel{}
    entity1.SetKey(key)
    if err := db.Get(&entity1, nil); err != nil {
        t.Fatal(err)
    }
    if !entity1.Ok() {
        t.Fatal("db: should exist")
    }
    if entity1.Content != entity.Content {
        t.Fatal("db: entity content mismatch", entity1.Content)
    }

    entity2 := TestModel{}
    entity2.SetKey(key)
    if err := db.Get(&entity2, nil); err != nil {
        t.Fatal(err)
    }
    if !entity2.Ok() {
        t.Fatal("db: should exist")
    }
    if entity2.Content != entity.Content {
        t.Fatal("db: entity content mismatch", entity2.Content)
    }

    entity3 := TestModel{}
    entity3.SetKey(datastore.NewKey(context, "test", "123", 0, nil))
    if err := db.Get(&entity3, nil); err != nil {
        t.Fatal(err)
    }
    if entity3.Ok() {
        t.Fatal("db: entity should not be found")
    }
}

func TestMulti(t *testing.T) {
    context, err := aetest.NewContext(nil)
    if err != nil {
        t.Fatal(err)
    }
    defer context.Close()

    cache := cache.New(context)
    defer cache.Clear()

    db := New(context, cache)

    entities := []Entity{
        &TestModel{datastore.NewKey(context, "test", "1", 0, nil), false, "Test1"},
        &TestModel{datastore.NewKey(context, "test", "2", 0, nil), false, "Test2"},
    }

    if err := db.PutMulti(entities, nil); err != nil {
        t.Fatal(err)
    }

    if cache.Size() != 2 {
        t.Fatal("db: cache size should be 2", cache.Size())
    }

    entities2 := []Entity{
        &TestModel{datastore.NewKey(context, "test", "1", 0, nil), false, ""},
        &TestModel{datastore.NewKey(context, "test", "2", 0, nil), false, ""},
        &TestModel{datastore.NewKey(context, "test", "3", 0, nil), false, ""},
        &TestModel{datastore.NewKey(context, "test", "4", 0, nil), false, ""},
    }

    if err := db.GetMulti(entities2, nil); err != nil {
        t.Fatal(err)
    }

    if !(entities2[0].Ok() && entities2[1].Ok() && !entities2[2].Ok() && !entities2[3].Ok()) {
        t.Fatal("db: two should exist")
    }

    if cache.Size() != 4 {
        t.Fatal("db: cache size should be 4", cache.Size())
    }
}
