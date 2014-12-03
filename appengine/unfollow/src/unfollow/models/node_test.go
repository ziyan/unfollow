package models

import (
    "appengine/aetest"
    "testing"
    "unfollow/utils/cache"
    "unfollow/utils/db"
)

func TestPutNode(t *testing.T) {
    context, err := aetest.NewContext(nil)
    if err != nil {
        t.Fatal(err)
    }
    defer context.Close()

    cache := cache.New(context)
    defer cache.Clear()

    db := db.New(context, cache)

    data := &Node{Name: "Name", Description: "Bio", ScreenName: "screenname", Avatar: ""}
    node1, err := PutNode(db, NodeKey(db, 1), data)
    if err != nil {
        t.Fatal(err)
    }

    if !NodeKey(db, 1).Equal(node1.Key()) {
        t.Fatal(node1.Key())
    }

    node, err := GetNodeByID(db, 1)
    if err != nil {
        t.Fatal(err)
    }

    if node == nil {
        t.Fatal(node)
    }

    if !NodeKey(db, 1).Equal(node.Key()) {
        t.Fatal(node.Key())
    }

    node2, err := PutNode(db, NodeKey(db, 2), data)
    if err != nil {
        t.Fatal(err)
    }

    if !NodeKey(db, 2).Equal(node2.Key()) {
        t.Fatal(node2.Key())
    }

    node, err = GetNodeByID(db, 2)
    if err != nil {
        t.Fatal(err)
    }

    if node == nil {
        t.Fatal(node)
    }

    if !NodeKey(db, 2).Equal(node.Key()) {
        t.Fatal(node.Key())
    }
}
