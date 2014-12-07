package network

import (
    "appengine/taskqueue"
    "unfollow/models"
    "unfollow/task"
    "unfollow/web"
    "unfollow/utils/twitter"
    "strconv"
)

// discover a list of nodes
var _ = task.Handle("network:discover", "/network/discover", func(handler *web.Handler) (interface{}, error) {

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
        id, err := strconv.ParseInt(task.Name, 10, 64)
        if err != nil {
            panic(err)
        }

        ids = append(ids, id)
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

    // check if the node already exists
    node, err := models.GetNodeByID(handler.Database, id)
    if err != nil {
        return nil, err
    }

    // node already exists, no work need to be done here
    if node != nil {
        return nil, nil
    }

    t := twitter.New(handler.Context, nil)
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

    if _, err := models.PutNode(handler.Database, node.Key(), node); err != nil {
        return nil, err
    }

    return nil, nil
})
