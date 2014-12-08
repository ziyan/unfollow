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

    Name        string `datastore:"name,noindex" json:"name"`
    Description string `datastore:"description,noindex" json:"description"`
    Location    string `datastore:"location,noindex" json:"location"`
    Website     string `datastore:"website,noindex" json:"website"`

    ScreenName string `datastore:"screen_name,noindex" json:"screen_name"`
    Avatar     string `datastore:"avatar,noindex" json:"avatar"`

    Verified      bool  `datastore:"verified,noindex" json:"verified"`
    Protected     bool  `datastore:"protected,noindex" json:"protected"`
    Contributed   bool  `datastore:"contributed,noindex" json:"contributed"`
    Default       bool  `datastore:"default,noindex" json:"default"`
    DefaultAvatar bool  `datastore:"default_avatar,noindex" json:"default_avatar"`
    Created       int64 `datastore:"created,noindex" json:"created"`

    TweetsCount    int64 `datastore:"tweets_count,noindex" json:"tweets_count"`
    ListsCount     int64 `datastore:"lists_count,noindex" json:"lists_count"`
    FavoritesCount int64 `datastore:"favorites_count,noindex" json:"favorites_count"`

    FriendsCount   int64 `datastore:"friends_count,noindex" json:"friends_count"`
    FollowersCount int64 `datastore:"followers_count,noindex" json:"followers_count"`

    FriendsIDs   []int64 `datastore:"friends_ids,noindex" json:"friends_ids"`
    FollowersIDs []int64 `datastore:"followers_ids,noindex" json:"followers_ids"`
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

func GetNodesByIDs(db *db.Database, ids []int64) ([]*Node, error) {
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

func PutNode(db *db.Database, node *Node) error {
    if err := db.Put(node, nil); err != nil {
        return err
    }

    return nil
}

func PutNodes(d *db.Database, nodes []*Node) error {
    entities := make([]db.Entity, 0, len(nodes))
    for _, node := range nodes {
        entities = append(entities, node)
    }

    if err := d.PutMulti(entities, nil); err != nil {
        return err
    }

    return nil
}
