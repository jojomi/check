package check

import (
	"net/http"
	"strconv"
)

// HTTPStatusCheck is a http statuscode check.
type HTTPStatusCheck struct {
	BaseCheck

	url string
}

// Execute does the check
func (h *HTTPStatusCheck) Execute() Result {
	resp, err := http.Get(h.url)

	result := h.getEmptyResult()
	if err != nil {
		result.Status = StatusCritical
		result.Message = err.Error()
	} else {
		code := resp.StatusCode
		validStatusList := h.GetExpected().([]int)

		result.Status = StatusCritical
		result.Value = strconv.Itoa(code)

		for _, validCode := range validStatusList {
			if code == validCode {
				result.Status = StatusOK
			}
		}
	}
	h.SetResult(result)
	return result
}

func (c *HTTPStatusCheck) Parse(configMap ConfigMap) {
	url, err := configMap.GetString("url")
	if err == nil {
		c.url = url
	}

	validStatus, err := configMap.GetIntArray("expected")
	if err == nil {
		c.SetExpected(validStatus)
	}

	c.ParseBaseData(configMap)
}
