package notify

import (
    "log"
    "strings"

    "github.com/containrrr/shoutrrr"
    "github.com/sjdaws/dls/internal/global"
)

func Message(text string) {
    urls := strings.Split(global.NotificationUrls, " ")

    for _, url := range urls {
        url = strings.TrimSpace(url)
        if url == "" {
            continue
        }

        err := shoutrrr.Send(url, text)
        if err != nil {
            log.Printf("ERROR: notify: %v", err)
        }
    }
}
