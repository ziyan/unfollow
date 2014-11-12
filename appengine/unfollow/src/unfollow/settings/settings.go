package settings

import (
    "appengine"
    "crypto/sha1"
    "encoding/hex"
    "github.com/ziyan/oauth"
    "net/url"
    "os"
    "time"
)

var DEBUG = appengine.IsDevAppServer()

var SECURE = !DEBUG

var VERSION = func() string {
    version := os.Getenv("PWD")
    if DEBUG {
        version = time.Now().String()
    }
    hasher := sha1.New()
    hasher.Write([]byte(version))
    return hex.EncodeToString(hasher.Sum(nil))
}()

var STATIC = "/static"

var SITE = "https://www.unfollow.io/"
var URL = func() *url.URL {
    url, err := url.Parse(SITE)
    if err != nil {
        panic(err)
    }
    return url
}()

var ANALYTICS = "UA-53151707-1"

var SECRET = []byte("RAQAyy7AkR84rznnqDhOnMYZvnFDiHbe0HSF650lgVVkrketrEPmY2130GVxBNyADfS8eDFHNKf")

var EMAIL_FROM = "hello@unfollow.io"
var EMAIL_API_KEY = func() string {
    if DEBUG {
        return "POSTMARK_API_TEST"
    }
    return "54dd1294-03ea-48ce-b845-fb9267acd764"
}()

var DEFAULT_LOCALE = "en_US"
var LOCALES = map[string]string{
    "en_US": "English",
    "zh_CN": "简体中文",
    "ja_JP": "日本語",
}
var LOCALE_PATTERNS = map[string]string{
    "zh": "zh_CN",
    "ja": "ja_JP",
}

var TWITTER_CONSUMER = oauth.NewToken("BBGDVII6UxbTvD77pN6eKvHq9", "xeSZo1SHASytg52DXAt4XRowmZJma5X2ZHoHVuZuCpWHvVrO6o")
