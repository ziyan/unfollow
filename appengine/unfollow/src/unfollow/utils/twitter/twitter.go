package twitter

import (
    "appengine"
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
    "unfollow/settings"
    "unfollow/utils/security"
)

var (
    ErrCallbackUnconfirmed = errors.New("twitter: callback unconfirmed")
    ErrNotFound            = errors.New("twitter: not found")
    ErrRateLimitReached    = errors.New("twitter: rate limit reached")
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
        Context:     context,
        Client:      urlfetch.Client(context),
        AccessToken: accessToken,
    }
}

func (twitter *Twitter) Request(method, path string, values url.Values, payload, data interface{}) error {

    url := url.URL{
        Scheme:   "https",
        Host:     "api.twitter.com",
        Path:     path,
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

    if err := oauth.SignRequest(request, settings.TWITTER_CONSUMER, twitter.AccessToken, security.GenerateRandomHexString(16)); err != nil {
        return err
    }

    response, err := twitter.Client.Do(request)
    if err != nil {
        return err
    }
    defer response.Body.Close()

    // report rate limit
    twitter.Context.Infof("twitter: ratelimit: limit = %s, remain = %s, reset = %s",
        response.Header.Get("X-Rate-Limit-Limit"),
        response.Header.Get("X-Rate-Limit-Remaining"),
        response.Header.Get("X-Rate-Limit-Reset"))

    switch {
    case response.StatusCode >= 200 && response.StatusCode < 300:
    case response.StatusCode == 404:
        return ErrNotFound
    case response.StatusCode == 429:
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

func (twitter *Twitter) Get(path string, values url.Values, data interface{}) error {
    return twitter.Request("GET", path, values, nil, data)
}

func (twitter *Twitter) Post(path string, values url.Values, payload, data interface{}) error {
    return twitter.Request("POST", path, values, payload, data)
}
