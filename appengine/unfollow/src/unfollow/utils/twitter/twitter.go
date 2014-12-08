package twitter

import (
    "appengine"
    "appengine/taskqueue"
    "appengine/urlfetch"
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/ziyan/oauth"
    "io"
    "io/ioutil"
    "net/http"
    "net/url"
    "strconv"
    "time"
    "unfollow/settings"
    "unfollow/utils/security"
)

var (
    ErrCallbackUnconfirmed = errors.New("twitter: callback unconfirmed")
    ErrNotFound            = errors.New("twitter: not found")
    ErrRateLimitReached    = errors.New("twitter: rate limit reached")

    API_ACCOUNT_VERIFY_CREDENTIALS = "account_verify_credentials"
    API_USERS_SHOW                 = "users_show"
    API_USERS_LOOKUP               = "users_lookup"
    API_FOLLOWERS_LIST             = "followers_list"
    API_FRIENDS_LIST               = "friends_list"
    API_FOLLOWERS_IDS              = "followers_ids"
    API_FRIENDS_IDS                = "friends_ids"
    API_SERACH_TWEETS              = "search_tweets"
    API_STATUSES_MENTIONS_TIMELINE = "statuses_mentions_timeline"

    API_PATHS = map[string]string{
        API_ACCOUNT_VERIFY_CREDENTIALS: "/1.1/account/verify_credentials.json",
        API_USERS_SHOW:                 "/1.1/users/show.json",
        API_USERS_LOOKUP:               "/1.1/users/lookup.json",
        API_FOLLOWERS_LIST:             "/1.1/followers/list.json",
        API_FRIENDS_LIST:               "/1.1/friends/list.json",
        API_FOLLOWERS_IDS:              "/1.1/followers/ids.json",
        API_FRIENDS_IDS:                "/1.1/friends/ids.json",
        API_SERACH_TWEETS:              "/1.1/search/tweets.json",
        API_STATUSES_MENTIONS_TIMELINE: "/1.1/statuses/mentions_timeline.json",
    }

    API_POOL = []string{
        API_USERS_SHOW,
        API_USERS_LOOKUP,
        API_FOLLOWERS_LIST,
        API_FRIENDS_LIST,
        API_FOLLOWERS_IDS,
        API_FRIENDS_IDS,
        API_SERACH_TWEETS,
    }
)

func GetRequestToken(context appengine.Context, callback string) (*oauth.Token, error) {

    values := url.Values{
        "oauth_callback": {callback},
    }

    u, err := url.Parse("https://api.twitter.com/oauth/request_token")
    if err != nil {
        return nil, err
    }

    u.RawQuery = values.Encode()

    request, err := http.NewRequest("POST", u.String(), nil)
    if err != nil {
        return nil, err
    }

    request.Form = values

    if err := oauth.SignRequest(request, settings.TWITTER_CONSUMER, nil, security.GenerateRandomHexString(16)); err != nil {
        return nil, err
    }

    client := urlfetch.Client(context)
    response, err := client.Do(request)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    switch {
    case response.StatusCode == 200:
    default:
        return nil, errors.New(fmt.Sprintf("twitter: server response status code is %d", response.StatusCode))
    }

    buffer, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }

    context.Infof("twitter: request_token: %s", string(buffer))

    token, err := oauth.DecodeToken(string(buffer))
    if err != nil {
        return nil, err
    }

    query, err := url.ParseQuery(string(buffer))
    if err != nil {
        return nil, err
    }

    if query.Get("oauth_callback_confirmed") != "true" {
        return nil, ErrCallbackUnconfirmed
    }

    return token, nil
}

func CreateAuthorizeUrl(token *oauth.Token) (string, error) {

    u, err := url.Parse("https://api.twitter.com/oauth/authenticate")
    if err != nil {
        return "", err
    }

    u.RawQuery = url.Values{
        "oauth_token": {token.Key()},
    }.Encode()

    return u.String(), nil
}

func GetAccessToken(context appengine.Context, token *oauth.Token, verifier string) (*oauth.Token, error) {

    values := url.Values{
        "oauth_verifier": {verifier},
    }

    u, err := url.Parse("https://api.twitter.com/oauth/access_token")
    if err != nil {
        return nil, err
    }

    u.RawQuery = values.Encode()

    request, err := http.NewRequest("POST", u.String(), nil)
    if err != nil {
        return nil, err
    }

    request.Form = values

    if err := oauth.SignRequest(request, settings.TWITTER_CONSUMER, token, security.GenerateRandomHexString(16)); err != nil {
        return nil, err
    }

    client := urlfetch.Client(context)
    response, err := client.Do(request)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    switch {
    case response.StatusCode == 200:
    default:
        return nil, errors.New(fmt.Sprintf("twitter: server response status code is %d", response.StatusCode))
    }

    buffer, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return nil, err
    }

    context.Infof("twitter: access_token: %s", string(buffer))
    return oauth.DecodeToken(string(buffer))
}

type Twitter struct {
    Context     appengine.Context
    Client      *http.Client
    AccessToken *oauth.Token
}

func New(context appengine.Context, accessToken *oauth.Token) *Twitter {
    return &Twitter{
        Context: context,
        Client: &http.Client{
            Transport: &urlfetch.Transport{
                Context:                       context,
                Deadline:                      time.Second * 60,
                AllowInvalidServerCertificate: false,
            },
        },
        AccessToken: accessToken,
    }
}

func (twitter *Twitter) Request(method, api string, values url.Values, payload, data interface{}) error {
    for {
        var task *taskqueue.Task

        accessToken := twitter.AccessToken
        if accessToken == nil {
            // get access token from queue
            var err error
            if accessToken, task, err = twitter.LeaseAccessToken(api); err != nil {
                return err
            }
        }

        url := url.URL{
            Scheme:   "https",
            Host:     "api.twitter.com",
            Path:     API_PATHS[api],
            RawQuery: values.Encode(),
        }
        twitter.Context.Infof("twitter: request: %s %s", method, url.String())

        var body io.Reader = nil
        if payload != nil {
            buffer, err := json.Marshal(payload)
            if err != nil {
                return err
            }
            body = bytes.NewReader(buffer)
        }

        request, err := http.NewRequest(method, url.String(), body)
        if err != nil {
            return err
        }

        request.Header.Set("Accept", "application/json")
        request.Form = values

        if err := oauth.SignRequest(request, settings.TWITTER_CONSUMER, accessToken, security.GenerateRandomHexString(16)); err != nil {
            return err
        }

        response, err := twitter.Client.Do(request)
        if err != nil {
            return err
        }
        defer response.Body.Close()

        // report rate limit
        limit, err := strconv.ParseInt(response.Header.Get("X-Rate-Limit-Limit"), 10, 64)
        if err != nil {
            limit = 15
        }

        remaining, err := strconv.ParseInt(response.Header.Get("X-Rate-Limit-Remaining"), 10, 64)
        if err != nil {
            remaining = 0
        }

        reset, err := strconv.ParseInt(response.Header.Get("X-Rate-Limit-Reset"), 10, 64)
        if err != nil {
            reset = time.Now().Unix() + 15 * 60
        }

        twitter.Context.Infof("twitter: ratelimit: limit = %d, remain = %d, reset = %d", limit, remaining, reset)

        // release token back to the pool
        if task != nil {
            if err := twitter.ReleaseAccessToken(task, limit, remaining, reset); err != nil {
                return err
            }
        }

        switch {
        case response.StatusCode >= 200 && response.StatusCode < 300:
        case response.StatusCode == 404:
            return ErrNotFound
        case response.StatusCode == 429:
            if twitter.AccessToken == nil {
                // try a different token
                continue
            }
            return ErrRateLimitReached
        default:
            if buffer, err := ioutil.ReadAll(response.Body); err == nil {
                twitter.Context.Errorf("twitter: error: %v", string(buffer))
            }
            return errors.New(fmt.Sprintf("twitter: server response status code is %d", response.StatusCode))
        }

        // parse the response
        if data != nil {
            if err := json.NewDecoder(response.Body).Decode(data); err != nil {
                return err
            }
        }

        return nil
    }
}

func (twitter *Twitter) Get(api string, values url.Values, data interface{}) error {
    return twitter.Request("GET", api, values, nil, data)
}

func (twitter *Twitter) Post(api string, values url.Values, payload, data interface{}) error {
    return twitter.Request("POST", api, values, payload, data)
}
