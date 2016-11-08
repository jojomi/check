package check

import "github.com/paulstuart/ping"

// NetPingCheck is doing a ping.
type NetPingCheck struct {
	BaseCheck

	host    string
	timeout int
}

// Execute does a net check
// See https://github.com/jwendel/ping/blob/master/README.md for appropriate access rights
func (n *NetPingCheck) Execute() Result {
	timeout := 10
	err := ping.Pinger(n.host, timeout)

	result := Result{}

	if err != nil {
		result.Status = StatusCritical
		result.Message = err.Error()
	} else {
		result.Status = StatusOK
	}
	n.SetResult(result)
	return result
}

func (n *NetPingCheck) Parse(configMap ConfigMap) {
	n.ParseBaseData(configMap)

	host, err := configMap.GetString("host")
	if err == nil {
		n.host = host
	}

	n.timeout = 10
	timeout, err := configMap.GetInt("timeout")
	if err == nil {
		n.timeout = timeout
	}
}
