package models

import (
    "appengine/datastore"
)

type LocalizedString struct {
    Locale string `datastore:"locale,noindex"`
    String string `datastore:"string,noindex"`
}

func MapLocalizedStrings(str []LocalizedString) map[string]string {
    m := make(map[string]string, len(str))
    for _, s := range str {
        m[s.Locale] = s.String
    }
    return m
}

func UnmapLocalizedStrings(str map[string]string) []LocalizedString {
    m := make([]LocalizedString, 0, len(str))
    for l, s := range str {
        m = append(m, LocalizedString{l, s})
    }
    return m
}

func KeysEqual(a, b []*datastore.Key) bool {
    if len(a) != len(b) {
        return false
    }
    for i := 0; i < len(a); i++ {
        if a[i] != nil && a[i].Equal(b[i]) {
            continue
        }
        return false
    }
    return true
}

func FilterKeys(keys []*datastore.Key, filter func(*datastore.Key) bool) []*datastore.Key {
    var filtered []*datastore.Key
    for _, key := range keys {
        if filter(key) {
            filtered = append(filtered, key)
        }
    }
    return filtered
}
