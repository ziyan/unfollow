package network

import (
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

    // lookup the users
    t := twitter.New(handler.Context, nil)
    users, err := t.Users(ids)
    if err != nil {
        return nil, err
    }

    // convert to internal structure node
    nodes := make([]*models.Node, 0, len(users))
    for _, user := range users {
        node := TwitterUserToNode(user)
        node.SetKey(models.NodeKey(handler.Database, user.ID))
        nodes = append(nodes, node)
    }

    if err := UpdateNodes(handler.Database, nodes); err != nil {
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

    node, err := DiscoverNode(handler.Database, id)
    if err != nil {
        return nil, err
    }

    return node, nil
})
