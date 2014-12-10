package network

import (
    "strconv"
    "strings"
    "unfollow/api"
    "unfollow/models"
    "unfollow/web"
)

var _ = api.Get("network:nodes", "/network/nodes", func(handler *web.Handler) (interface{}, error) {
    strs := strings.Split(handler.Request.URL.Query().Get("ids"), ",")
    ids := make([]int64, 0, len(strs))
    for _, str := range strs {
        id, err := strconv.ParseInt(str, 10, 64)
        if err != nil {
            return nil, err
        }
        ids = append(ids, id)
    }

    nodes, err := models.GetNodesByIDs(handler.Database, ids)
    if err != nil {
        return nil, err
    }

    results := make(map[string]*models.Node)
    for _, node := range nodes {
        if !node.Ok() {
            continue
        }
        results[strconv.FormatInt(node.ID(), 10)] = node
    }

    return struct {
        Nodes map[string]*models.Node `json:"nodes"`
    }{results}, nil
})
