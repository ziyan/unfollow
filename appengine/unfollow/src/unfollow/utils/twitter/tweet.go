package twitter

import (
    "net/url"
)

type Tweet struct {
    ID                  int64  `json:"id"`
    User                User   `json:"user"`
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
    if err := twitter.Get(API_SERACH_TWEETS, values, &result); err != nil {
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
    if err := twitter.Get(API_STATUSES_MENTIONS_TIMELINE, values, &tweets); err != nil {
        return nil, err
    }

    twitter.Context.Infof("twitter: mentions: %v", tweets)
    return tweets, nil
}
