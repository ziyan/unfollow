package models

import (
    "appengine/datastore"
    "unfollow/utils/db"
)

const (
    NODE_KIND = "node"
)

type Node struct {
    key *datastore.Key `datastore:"-"`
    ok  bool           `datastore:"-"`

    Name        string `datastore:"name,noindex"`
    Description string `datastore:"description,noindex"`
    Location    string `datastore:"location,noindex"`
    Website     string `datastore:"website,noindex"`

    ScreenName string `datastore:"screen_name,noindex"`
    Avatar     string `datastore:"avatar,noindex"`

    FriendsCount   int64 `datastore:"friend_count,noindex"`
    FollowersCount int64 `datastore:"follower_count,noindex"`

    FriendsIDs   []int64 `datastore:"friends_ids,noindex"`
    FollowersIDs []int64 `datastore:"followers_ids,noindex"`
}

func (u *Node) Key() *datastore.Key {
    return u.key
}

func (u *Node) SetKey(key *datastore.Key) {
    u.key = key
}

func (u *Node) Ok() bool {
    return u.ok
}

func (u *Node) SetOk(ok bool) {
    u.ok = ok
}

func (u *Node) ID() int64 {
    return u.key.IntID()
}

func NodeKey(db *db.Database, id int64) *datastore.Key {
    return db.Key(NODE_KIND, "", id, nil)
}

func GetNode(db *db.Database, key *datastore.Key) (*Node, error) {
    node := &Node{}
    node.SetKey(key)
    if err := db.Get(node, nil); err != nil {
        return nil, err
    }
    if !node.Ok() {
        return nil, nil
    }
    return node, nil
}

func GetNodes(d *db.Database, keys []*datastore.Key) ([]*Node, error) {
    nodes := make([]*Node, 0, len(keys))
    entities := make([]db.Entity, 0, len(keys))
    for _, key := range keys {
        node := &Node{}
        node.SetKey(key)
        nodes = append(nodes, node)
        entities = append(entities, node)
    }

    if err := d.GetMulti(entities, nil); err != nil {
        return nil, err
    }

    return nodes, nil
}

func GetNodeByID(db *db.Database, id int64) (*Node, error) {
    return GetNode(db, NodeKey(db, id))
}

func GetNodesByID(db *db.Database, ids []int64) ([]*Node, error) {
    keys := make([]*datastore.Key, 0, len(ids))
    for _, id := range ids {
        keys = append(keys, NodeKey(db, id))
    }

    nodes, err := GetNodes(db, keys)
    if err != nil {
        return nil, err
    }

    return nodes, nil
}

func PutNode(db *db.Database, key *datastore.Key, data *Node) (*Node, error) {
    node := &Node{}
    node.SetKey(key)

    if err := db.Put(node, nil); err != nil {
        return nil, err
    }

    return node, nil
}
