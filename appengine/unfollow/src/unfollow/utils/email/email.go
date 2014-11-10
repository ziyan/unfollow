package email

import (
    "appengine"
    "appengine/delay"
    "appengine/urlfetch"
    "bytes"
    "unfollow/settings"
    "encoding/json"
    "fmt"
    "net/http"
)

var send = delay.Func("email", func(context appengine.Context, sender, recipient, to, subject, text, html, tag string) error {

    context.Infof("email: sender = %s, recipient = %s, to = %s, subject = %s, tag = %s", sender, recipient, to, subject, tag)
    context.Infof("email: text = %s", text)
    context.Infof("email: html = %s", html)

    data := struct {
        From, To, Subject, Tag, HtmlBody, TextBody string
    }{
        fmt.Sprintf("%s <%s>", sender, settings.EMAIL_FROM),
        fmt.Sprintf("%s <%s>", recipient, to),
        subject,
        tag,
        html,
        text,
    }

    // encode to json
    buffer := new(bytes.Buffer)
    if err := json.NewEncoder(buffer).Encode(data); err != nil {
        return err
    }
    reader := bytes.NewReader(buffer.Bytes())

    // build the request
    client := urlfetch.Client(context)
    request, err := http.NewRequest("POST", "https://api.postmarkapp.com/email", reader)
    if err != nil {
        return err
    }

    request.Header.Set("X-Postmark-Server-Token", settings.EMAIL_API_KEY)
    request.Header.Set("Content-Type", "application/json")
    request.Header.Set("Accept", "applicaiton/json")

    // send the request
    response, err := client.Do(request)
    if err != nil {
        return err
    }
    defer response.Body.Close()

    // parse the response
    reply := struct {
        To, Message, MessageID string
        ErrorCode              int
    }{}
    if err := json.NewDecoder(response.Body).Decode(&reply); err != nil {
        return err
    }

    context.Infof("email: to = %s, error_code = %d, message = %s, message_id = %s", reply.To, reply.ErrorCode, reply.Message, reply.MessageID)
    return nil
})

// Send an email.
func Send(context appengine.Context, sender, recipient, to, subject, text, html, tag string) {
    send.Call(context, sender, recipient, to, subject, text, html, tag)
}
