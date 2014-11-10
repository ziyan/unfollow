package templates

import (
    "unfollow/urls"
    "html/template"
)

var _ = Register("url", func(name string, pairs ...string) string {
    return urls.Reverse(name, pairs...).String()
})

var _ = Register("safe", func(output string) template.HTML {
    return template.HTML(output)
})
