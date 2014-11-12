package twitter

import (
    "net/url"
)

type Tweet struct {
    ID   int64 `json:"id"`
    User struct {
        ID           int64 `json:"id"`
        TweetsCount  int64 `json:"statuses_count"`
        FriendsCount int64 `json:"friends_count"`
    }   `json:"user"`
    InReplyToTweetID    int64  `json:"in_reply_to_status_id"`
    InReplyToScreenName string `json:"in_reply_to_screen_name"`
}

func (twitter *Twitter) Search(query string) ([]*Tweet, error) {
    result := struct {
        Tweets []*Tweet `json:"statuses"`
    }{}
    values := url.Values{
        "q":           {query},
        "result_type": {"recent"},
        "count":       {"100"},
    }
    if err := twitter.Get("/1.1/search/tweets.json", values, &result); err != nil {
        return nil, err
    }

    twitter.Context.Infof("twitter: search: %v", result)
    return result.Tweets, nil
}

func (twitter *Twitter) Mentions() ([]*Tweet, error) {
    tweets := make([]*Tweet, 0)
    values := url.Values{
        "count": {"200"},
    }
    if err := twitter.Get("/1.1/statuses/mentions_timeline.json", values, &tweets); err != nil {
        return nil, err
    }

    twitter.Context.Infof("twitter: mentions: %v", tweets)
    return tweets, nil
}
