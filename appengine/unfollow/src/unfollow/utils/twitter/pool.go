package twitter

import (
    "appengine"
    "crypto/sha1"
    "encoding/hex"
    "appengine/taskqueue"
    "github.com/ziyan/oauth"
    "time"
)

func (twitter *Twitter) LeaseAccessToken(api string) (*oauth.Token, *taskqueue.Task, error) {
    // lease a token
    tasks, err := taskqueue.LeaseByTag(twitter.Context, 1, "twitter", 15 * 60, api)
    if err != nil {
        return nil, nil, err
    }

    // no more token available
    if len(tasks) == 0 {
        return nil, nil, ErrRateLimitReached
    }
    if len(tasks) != 1 {
        panic("twitter: more than one token leased")
    }

    // decode token
    task := tasks[0]
    accessToken, err := oauth.DecodeToken(string(task.Payload))
    if err != nil {
        return nil, nil, err
    }

    return accessToken, task, nil
}

func (twitter *Twitter) ReleaseAccessToken(task *taskqueue.Task, limit, remaining, reset int64) error {
    // calculate the least time
    lease := 0
    if remaining == 0 {
        // if the token is depleted, modify its lease
        // to expire 1 second after the 15min window
        lease = int(reset - time.Now().Unix() + 1)
        if lease < 0 {
            lease = 1
        }
    }

    if err := taskqueue.ModifyLease(twitter.Context, task, "twitter", lease); err != nil {
        return err
    }

    return nil
}

func PoolAccessToken(context appengine.Context, token *oauth.Token) error {
    encoded := token.Encode()
    payload := []byte(encoded)
    tasks := make([]*taskqueue.Task, 0, len(API_POOL))

    for _, api := range API_POOL {
        hasher := sha1.New()
        hasher.Write([]byte(encoded))
        hasher.Write([]byte(api))

        task := &taskqueue.Task{
            Payload: payload,
            Method:  "PULL",
            Name: hex.EncodeToString(hasher.Sum(nil)),
            Tag: api,
        }
        tasks = append(tasks, task)
    }

    if _, err := taskqueue.AddMulti(context, tasks, "twitter"); err != nil {
        errs, ok := err.(appengine.MultiError)
        if !ok {
            return err
        }

        for _, err := range errs {
            if err == taskqueue.ErrTaskAlreadyAdded {
                err = nil
            }
            if err != nil {
                return err
            }
        }
    }

    return nil
}
