package network

import (
    "strconv"
    "unfollow/ajax"
    "unfollow/models"
    "unfollow/web"
)

var _ = ajax.Post("network:nodes", "/network/nodes", func(view *web.View) (interface{}, error) {
    if err := view.LoginRequired(); err != nil {
        return nil, err
    }

    data := struct {
        IDs []int64 `json:"ids"`
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

    if len(ids) > 0 {
        if err := QueueDiscoverNodes(view.Context, ids); err != nil {
            return nil, err
        }
    }

    return struct {
        Nodes map[string]*models.Node `json:"nodes"`
    }{results}, nil
})

var _ = ajax.Post("network:node", "/network/node", func(view *web.View) (interface{}, error) {
    if err := view.LoginRequired(); err != nil {
        return nil, err
    }

    data := struct {
        ID int64 `json:"id"`
    }{}

    if err := ajax.Decode(view, &data); err != nil {
        return nil, err
    }

    node, err := models.GetNodeByID(view.Database, data.ID)
    if err != nil {
        return nil, err
    }

    if node == nil ||
        (!node.Protected && node.FollowersCount > 0 && len(node.FollowersIDs) == 0) ||
        (!node.Protected && node.FriendsCount > 0 && len(node.FriendsIDs) == 0) {
        if err := ScheduleDiscoverNode(view.Context, data.ID); err != nil {
            return nil, err
        }
    }

    return struct {
        ID   int64        `json:"id"`
        Node *models.Node `json:"node"`
    }{data.ID, node}, nil
})
