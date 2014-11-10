package templates

import (
    "bytes"
    "html/template"
    "io"
    "log"
    "path"
    "strings"
    "sync"
)

var cache = make(map[string]*template.Template)
var registry = make(template.FuncMap)
var mutex sync.Mutex

func Render(writer io.Writer, data interface{}, filenames ...string) error {
    key := strings.Join(filenames, ",")

    mutex.Lock()
    t := cache[key]
    mutex.Unlock()

    if t == nil {
        t = template.New("").Funcs(registry)
        for _, filename := range filenames {
            t = template.Must(t.ParseFiles(path.Join("template", filename)))
        }
        t = template.Must(t.ParseGlob(path.Join("template", "inc", "*")))

        mutex.Lock()
        cache[key] = t
        mutex.Unlock()

        log.Println("web: cached template:", key)
    }
    return t.Execute(writer, data)
}

func RenderToString(data interface{}, filenames ...string) (string, error) {
    var buffer bytes.Buffer
    if err := Render(&buffer, data, filenames...); err != nil {
        return "", err
    }
    return strings.TrimSpace(buffer.String()), nil
}

func Register(name string, function interface{}) interface{} {
    registry[name] = function
    return function
}
