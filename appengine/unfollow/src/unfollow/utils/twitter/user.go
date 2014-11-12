package twitter

import (
    "net/url"
    "strconv"
)

type User struct {
    ID                   int64  `json:"id"`
    ScreenName           string `json:"screen_name"`
    Name                 string `json:"name"`
    Description          string `json:"description"`
    ProfileImageUrlHttps string `json:"profile_image_url_https"`
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
