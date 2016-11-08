package check

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// HTTPContentCheck is a http content check.
type HTTPContentCheck struct {
	BaseCheck

	url string
}

// Execute does the check
func (h *HTTPContentCheck) Execute() Result {
	resp, err := http.Get(h.url)

	result := h.getEmptyResult()
	if err != nil {
		result.Status = StatusCritical
		result.Message = err.Error()
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			result.Status = StatusCritical
			result.Message = err.Error()
		} else {
			content := string(body)
			result.Value = content
			expectedContent := h.GetExpected().(string)

			if strings.Contains(content, expectedContent) {
				result.Status = StatusOK
			}
		}
	}
	h.SetResult(result)
	return result
}

func (c *HTTPContentCheck) Parse(configMap ConfigMap) {
	c.ParseBaseData(configMap)

	url, err := configMap.GetString("url")
	if err == nil {
		c.url = url
	}
}
