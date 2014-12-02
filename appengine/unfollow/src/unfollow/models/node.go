package models

import (
    "appengine/datastore"
    "unfollow/utils/db"
    "reflect"
)

const (
    NODE_KIND = "node"
)

type Node struct {
    key *datastore.Key `datastore:"-"`
    ok  bool           `datastore:"-"`

    Name string `datastore:"name,noindex"`
    Bio  string `datastore:"bio,noindex"`

    Username string `datastore:"username,noindex"`
    Avatar   string `datastore:"avatar,noindex"`

    FriendCount int64 `datastore:"friend_count,noindex"`
    FollowerCount int64 `datastore:"follower_count,noindex"`

    FriendIDs []int64 `datastore:"friend_ids,noindex"`
    FollowerIDs []int64 `datastore:"follower_ids,noindex"`
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

func GetNodeByID(db *db.Database, id int64) (*Node, error) {
    return GetNode(db, NodeKey(db, id))
}

func PutNode(d *db.Database, key *datastore.Key, data *Node) (*Node, error) {
    node := &Node{}
    node.SetKey(key)
    if err := d.Transaction(func(db *db.Database) error {
        // get existing node
        if err := db.Get(node, nil); err != nil {
            return err
        }

        if node.Ok() {
            if node.Name == data.Name &&
                node.Bio == data.Bio &&
                node.Avatar == data.Avatar &&
                node.Username == data.Username &&
                node.FriendCount == data.FriendCount &&
                node.FollowerCount == data.FollowerCount &&
                reflect.DeepEqual(node.FriendIDs, data.FriendIDs) &&
                reflect.DeepEqual(node.FollowerIDs, data.FollowerIDs) {
                // nothing changed
                return nil
            }
        }

        // update node
        node.Name = data.Name
        node.Bio = data.Bio
        node.Avatar = data.Avatar
        node.Username = data.Username
        node.FriendCount = data.FriendCount
        node.FollowerCount = data.FollowerCount
        node.FriendIDs = data.FriendIDs
        node.FollowerIDs = data.FollowerIDs

        // save the node
        if err := db.Put(node, nil); err != nil {
            return err
        }

        return nil
    }, false); err != nil {
        return nil, err
    }

    return node, nil
}
