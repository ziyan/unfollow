package network

import (
    "appengine"
    "appengine/taskqueue"
    "bytes"
    "encoding/binary"
    "strconv"
    "time"
    "unfollow/models"
    "unfollow/urls"
    "unfollow/utils/db"
    "unfollow/utils/twitter"
)

func QueueDiscoverNodes(context appengine.Context, ids []int64) error {
    // queue discover
    tasks := make([]*taskqueue.Task, 0, len(ids))

    // discover friends
    for _, id := range ids {
        buffer := new(bytes.Buffer)
        if err := binary.Write(buffer, binary.LittleEndian, id); err != nil {
            return err
        }

        task := &taskqueue.Task{
            Name:    strconv.FormatInt(id, 10),
            Method:  "PULL",
            Payload: buffer.Bytes(),
        }
        tasks = append(tasks, task)
    }

    // batch add tasks
    for len(tasks) > 0 {
        size := len(tasks)
        if size > 20 {
            size = 20
        }
        batch := tasks[:size]
        tasks = tasks[size:]

        if _, err := taskqueue.AddMulti(context, batch, "network"); err != nil {
            errs, ok := err.(appengine.MultiError)
            if !ok {
                return err
            }

            for _, err := range errs {
                if err == taskqueue.ErrTaskAlreadyAdded {
                    err = nil
                }
                if err != nil {
                    return err
                }
            }
        }
    }

    // trigger a schedule
    if err := ScheduleDiscoverNodes(context); err != nil {
        return err
    }

    return nil
}

func ScheduleDiscoverNodes(context appengine.Context) error {
    url := urls.Reverse("task:network:nodes")
    if _, err := taskqueue.Add(context, &taskqueue.Task{Path: url.Path, Method: "POST"}, "default"); err != nil {
        return err
    }
    return nil
}

func ScheduleDiscoverNode(context appengine.Context, id int64) error {
    url := urls.Reverse("task:network:node", "id", strconv.FormatInt(id, 10))
    if _, err := taskqueue.Add(context, &taskqueue.Task{Path: url.Path, Method: "POST"}, "default"); err != nil {
        return err
    }
    return nil
}

func TwitterUserToNode(user *twitter.User) *models.Node {
    node := &models.Node{}
    node.Name = user.Name
    node.Description = user.Description
    node.Location = user.Location
    node.Website = user.URL
    node.ScreenName = user.ScreenName
    node.Avatar = user.ProfileImageUrlHttps
    node.FriendsCount = user.FriendsCount
    node.FollowersCount = user.FollowersCount
    node.ListsCount = user.ListedCount
    node.TweetsCount = user.StatusesCount
    node.Verified = user.Verified
    node.Protected = user.Protected
    node.Contributed = user.ContributorsEnabled
    node.Default = user.DefaultProfile
    node.DefaultAvatar = user.DefaultProfileImage

    t, err := time.Parse(time.RubyDate, user.CreatedAt)
    if err == nil {
        node.Created = t.Unix()
    }

    return node
}

func UpdateNodes(db *db.Database, nodes []*models.Node) error {

    index := make(map[int64]*models.Node)
    ids := make([]int64, 0, len(nodes))
    for _, node := range nodes {
        if index[node.ID()] != nil {
            continue
        }

        index[node.ID()] = node
        ids = append(ids, node.ID())
    }

    existings, err := models.GetNodesByIDs(db, ids)
    if err != nil {
        return err
    }

    for _, existing := range existings {
        node := index[existing.ID()]
        if node != nil && existing.Ok() {
            node.FriendsIDs = existing.FriendsIDs
            node.FollowersIDs = existing.FollowersIDs
        }
    }

    uniques := make([]*models.Node, 0, len(index))
    for _, node := range index {
        uniques = append(uniques, node)
    }

    // save them all
    if err := models.PutNodes(db, uniques); err != nil {
        return err
    }

    return nil
}

func DiscoverNodes(db *db.Database, ids []int64) ([]*models.Node, error) {
    // lookup the users
    t := twitter.New(db.Context, nil)
    users, err := t.Users(ids)
    if err != nil {
        return nil, err
    }
    if len(users) == 0 {
        return nil, nil
    }

    // convert to internal structure node
    nodes := make([]*models.Node, 0, len(users))
    for _, user := range users {
        node := TwitterUserToNode(user)
        node.SetKey(models.NodeKey(db, user.ID))
        nodes = append(nodes, node)
    }

    if err := UpdateNodes(db, nodes); err != nil {
        return nil, err
    }

    return nodes, nil
}

func DiscoverNode(db *db.Database, id int64) (*models.Node, error) {
    t := twitter.New(db.Context, nil)

    user, err := t.User(id)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, nil
    }
    if user.ID != id {
        panic("network: twitter gave us a different user")
    }

    node := TwitterUserToNode(user)
    node.SetKey(models.NodeKey(db, id))

    node.FriendsIDs, err = t.FriendsIDs(id)
    if err != nil {
        return nil, err
    }

    node.FollowersIDs, err = t.FollowersIDs(id)
    if err != nil {
        return nil, err
    }

    if err := models.PutNode(db, node); err != nil {
        return nil, err
    }

    // discover recent friends and followers
    friends, err := t.Friends(id)
    if err != nil {
        return nil, err
    }

    followers, err := t.Followers(id)
    if err != nil {
        return nil, err
    }

    nodes := make([]*models.Node, 0, len(friends)+len(followers))
    for _, user := range append(friends, followers...) {
        node := TwitterUserToNode(user)
        node.SetKey(models.NodeKey(db, user.ID))
        nodes = append(nodes, node)
    }

    if err := UpdateNodes(db, nodes); err != nil {
        return nil, err
    }

    return node, nil
}
