package network

import (
    "strconv"
    "unfollow/ajax"
    "unfollow/models"
    "unfollow/web"
)

var _ = ajax.Get("network:nodes", "/network/nodes", func(view *web.View) (interface{}, error) {
    if err := view.LoginRequired(); err != nil {
        return nil, err
    }

    data := struct {
        IDs   []int64 `json:"ids"`
        Queue bool    `json:"queue"`
    }{}

    if err := ajax.Decode(view, &data); err != nil {
        return nil, err
    }

    nodes, err := models.GetNodesByIDs(view.Database, data.IDs)
    if err != nil {
        return nil, err
    }

    ids := make([]int64, 0, len(nodes))
    results := make(map[string]*models.Node)
    for _, node := range nodes {
        if !node.Ok() {
            ids = append(ids, node.ID())
            continue
        }
        results[strconv.FormatInt(node.ID(), 10)] = node
    }

    queued := false
    if data.Queue && len(ids) > 0 {
        if err := QueueDiscoverNodes(view.Context, ids); err != nil {
            return nil, err
        }
        queued = true
    }

    return struct {
        Nodes map[string]*models.Node `json:"nodes"`
        Queued bool `json:"queued"`
    }{results, queued}, nil
})

var _ = ajax.Get("network:node", "/network/node", func(view *web.View) (interface{}, error) {
    if err := view.LoginRequired(); err != nil {
        return nil, err
    }

    data := struct {
        ID    int64 `json:"id"`
        Queue bool  `json:"queue"`
    }{}

    if err := ajax.Decode(view, &data); err != nil {
        return nil, err
    }

    node, err := models.GetNodeByID(view.Database, data.ID)
    if err != nil {
        return nil, err
    }

    queued := false
    if data.Queue {
        if node == nil ||
            (node.FollowersCount > 0 && len(node.FollowersIDs) == 0) ||
            (node.FriendsCount > 0 && len(node.FriendsIDs) == 0) {
            if err := ScheduleDiscoverNode(view.Context, data.ID); err != nil {
                return nil, err
            }

            queued = true
        }
    }

    return struct {
        ID   int64        `json:"id"`
        Node *models.Node `json:"node"`
        Queued bool `json:"queued"`
    }{data.ID, node, queued}, nil
})
