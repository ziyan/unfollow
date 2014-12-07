package twitter

import (
    "net/url"
    "strconv"
    "strings"
)

type User struct {
    ID                   int64  `json:"id"`
    ScreenName           string `json:"screen_name"`
    Name                 string `json:"name"`
    Description          string `json:"description"`
    Location             string `json:"location"`
    URL                  string `json:"url"`
    ProfileImageUrlHttps string `json:"profile_image_url_https"`
    StatusesCount        int64  `json:"statuses_count"`
    FriendsCount         int64  `json:"friends_count"`
    FollowersCount       int64  `json:"followers_count"`
    ListedCount          int64  `json:"listed_count"`
    Verified             bool   `json:"verified"`
    Protected            bool   `json:"protected"`
    ContributorsEnabled  bool   `json:"contributors_enabled"`
    DefaultProfile       bool   `json:"default_profile"`
    DefaultProfileImage  bool   `json:"default_profile_image"`
    CreatedAt            string `json:"created_at"`
}

func (twitter *Twitter) VerifyCredentials() (*User, error) {
    user := User{}
    values := url.Values{
        "skip_status":      {"true"},
        "include_entities": {"false"},
    }
    if err := twitter.Get("/1.1/account/verify_credentials.json", values, &user); err != nil {
        return nil, err
    }

    twitter.Context.Infof("twitter: user: %v", user)
    return &user, nil
}

func (twitter *Twitter) LookupUsers(ids []int64) ([]*User, error) {
    strs := make([]string, 0, len(ids))
    for _, id := range ids {
        strs = append(strs, strconv.FormatInt(id, 10))
    }

    values := url.Values{
        "include_entities": {"false"},
        "user_id": {strings.Join(strs, ",")},
    }

    users := make([]*User, 0)
    if err := twitter.Get("/1.1/users/lookup.json", values, &users); err != nil {
        return nil, err
    }

    twitter.Context.Infof("twitter: users: %v", users)
    return users, nil
}

func (twitter *Twitter) Followers(id int64) ([]*User, error) {
    result := struct {
        Users []*User `json:"users"`
    }{}
    values := url.Values{
        "count":   {"200"},
        "user_id": {strconv.FormatInt(id, 10)},
    }
    if err := twitter.Get("/1.1/followers/list.json", values, &result); err != nil {
        return nil, err
    }

    twitter.Context.Infof("twitter: followers: %v", result)
    return result.Users, nil
}

func (twitter *Twitter) Friends(id int64) ([]*User, error) {
    result := struct {
        Users []*User `json:"users"`
    }{}
    values := url.Values{
        "count":   {"200"},
        "user_id": {strconv.FormatInt(id, 10)},
    }
    if err := twitter.Get("/1.1/friends/list.json", values, &result); err != nil {
        return nil, err
    }

    twitter.Context.Infof("twitter: friends: %v", result)
    return result.Users, nil
}

func (twitter *Twitter) FollowersIDs(id int64) ([]int64, error) {
    result := struct {
        IDs []int64 `json:"ids"`
    }{}
    values := url.Values{
        "count":   {"5000"},
        "user_id": {strconv.FormatInt(id, 10)},
    }
    if err := twitter.Get("/1.1/followers/ids.json", values, &result); err != nil {
        return nil, err
    }

    twitter.Context.Infof("twitter: followers: %v", result)
    return result.IDs, nil
}

func (twitter *Twitter) FriendsIDs(id int64) ([]int64, error) {
    result := struct {
        IDs []int64 `json:"ids"`
    }{}
    values := url.Values{
        "count":   {"5000"},
        "user_id": {strconv.FormatInt(id, 10)},
    }
    if err := twitter.Get("/1.1/friends/ids.json", values, &result); err != nil {
        return nil, err
    }

    twitter.Context.Infof("twitter: friends: %v", result)
    return result.IDs, nil
}
