package network

import (
    "appengine"
    "appengine/taskqueue"
    "strconv"
    "unfollow/models"
    "unfollow/task"
    "unfollow/utils/twitter"
    "unfollow/web"
    "bytes"
    "encoding/binary"
)

// discover a list of nodes
var _ = task.Handle("network:discover:nodes", "/network/discover", func(handler *web.Handler) (interface{}, error) {

    // lease a bunch of tasks
    tasks, err := taskqueue.Lease(handler.Context, 100, "discover", 60)
    if err != nil {
        return nil, err
    }

    // nothing to do
    if len(tasks) == 0 {
        return nil, nil
    }

    // get the list of user ids
    ids := make([]int64, 0, len(tasks))
    for _, task := range tasks {
        // we use task name to deduplicate and store the user id
        var id int64
        buffer := bytes.NewReader(task.Payload)
        if err := binary.Read(buffer, binary.LittleEndian, &id); err != nil {
            return nil, err
        }
        ids = append(ids, id)
    }
    handler.Context.Infof("network: discover: ids = %v", ids)

    // lookup existing nodes
    existings, err := models.GetNodesByIDs(handler.Database, ids)
    if err != nil {
        return nil, err
    }

    index := make(map[int64]*models.Node)
    for _, node := range existings {
        index[node.ID()] = node
    }

    // lookup the users
    t := twitter.New(handler.Context, nil)
    users, err := t.LookupUsers(ids)
    if err != nil {
        return nil, err
    }

    // convert to internal structure node
    nodes := make([]*models.Node, 0, len(users))
    for _, user := range users {
        node := TwitterUserToNode(user)
        node.SetKey(models.NodeKey(handler.Database, user.ID))

        // preserve the edges if the node already exists
        existing := existing[user.ID]
        if existing != nil && existing.Ok() {
            node.FriendsIDs = existing.FriendsIDs
            node.FollowersIDs = existing.FollowersIDs
        }

        nodes = append(nodes, node)
    }

    // save them all
    if err := models.PutNodes(handler.Database, nodes); err != nil {
        return nil, err
    }

    // delete the tasks leased
    if err := taskqueue.DeleteMulti(handler.Context, tasks, "discover"); err != nil {
        return nil, err
    }

    // schedule a callback to this task
    if err := Schedule(handler.Context); err != nil {
        return nil, err
    }

    return nil, nil
})

// discover edges for a particular node
var _ = task.Handle("network:discover:node", "/network/discover/{id:[0-9]+}", func(handler *web.Handler) (interface{}, error) {

    // user id must be an integer
    id, err := strconv.ParseInt(handler.Variables["id"], 10, 64)
    if err != nil {
        panic(err)
    }

    t := twitter.New(handler.Context, nil)

    // check if the node already exists
    key := models.NodeKey(handler.Database, id)
    node, err := models.GetNode(handler.Database, key)
    if err != nil {
        return nil, err
    }

    // if node does not exist lookup it up
    if node == nil {
        users, err := t.LookupUsers([]int64{id})
        if err != nil {
            return nil, err
        }
        if len(users) == 0 {
            // user not found, we are done
            return nil, nil
        }
        if len(users) != 1 {
            panic("network: lookup users returned more than one user")
        }
        node = TwitterUserToNode(users[0])
    }

    friends, err := t.FriendsIDs(id)
    if err != nil {
        return nil, err
    }

    followers, err := t.FollowersIDs(id)
    if err != nil {
        return nil, err
    }

    // save the node
    node.FriendsIDs = friends
    node.FollowersIDs = followers

    if _, err := models.PutNode(handler.Database, key, node); err != nil {
        return nil, err
    }

    // queue discover
    tasks := make([]*taskqueue.Task, 0, len(friends)+len(followers))

    // discover friends
    for _, id := range friends {
        buffer := new(bytes.Buffer)
        if err := binary.Write(buffer, binary.LittleEndian, id); err != nil {
            return nil, err
        }

        task := &taskqueue.Task{
            Name:    strconv.FormatInt(id, 10),
            Method:  "PULL",
            Payload: buffer.Bytes(),
        }
        tasks = append(tasks, task)
    }

    // discover followers
    for _, id := range followers {
        buffer := new(bytes.Buffer)
        if err := binary.Write(buffer, binary.LittleEndian, id); err != nil {
            return nil, err
        }

        task := &taskqueue.Task{
            Name:    strconv.FormatInt(id, 10),
            Method:  "PULL",
            Payload: buffer.Bytes(),
        }
        tasks = append(tasks, task)
    }

    for len(tasks) > 0 {
        size := len(tasks)
        if size > 20 {
            size = 20
        }
        batch := tasks[:size]
        tasks = tasks[size:]

        if _, err := taskqueue.AddMulti(handler.Context, batch, "discover"); err != nil {
            errs, ok := err.(appengine.MultiError)
            if !ok {
                return nil, err
            }

            for _, err := range errs {
                if err == taskqueue.ErrTaskAlreadyAdded {
                    err = nil
                }
                if err != nil {
                    return nil, err
                }
            }
        }
    }

    // trigger a schedule
    if err := Schedule(handler.Context); err != nil {
        return nil, err
    }

    return nil, nil
})
