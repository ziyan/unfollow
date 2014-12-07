package network

import (
    "appengine"
    "appengine/taskqueue"
    "unfollow/urls"
    "unfollow/utils/twitter"
    "unfollow/models"
    "time"
)

func Schedule(context appengine.Context) error {
    url := urls.Reverse("task:network:discover")
    if _, err := taskqueue.Add(context, &taskqueue.Task{
        Path:   url.Path,
        Method: "POST",
    }, "default"); err != nil {
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
