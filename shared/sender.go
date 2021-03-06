package shared

import (
	"fmt"
	"net/http"
	"net/url"
)

// SendToIFTTT sends the provided contents to IFTTT
func SendToIFTTT(iftttKey string, iftttTriggerName string, contents []string) {
	makerURL := fmt.Sprintf("https://maker.ifttt.com/trigger/%s/with/key/%s", iftttTriggerName, iftttKey)
	http.PostForm(makerURL, url.Values{"value1": contents})
}
