package check

import (
	"net"
	"strings"

	"github.com/bogdanovich/dns_resolver"
)

// NetDNSCheck is doing a dns ip lookup.
type NetDNSCheck struct {
	BaseCheck

	host string
}

// Execute does a dns check
func (n *NetDNSCheck) Execute() Result {
	servers := []string{"8.8.8.8", "8.8.4.4"}

	resolver := dns_resolver.New(servers)
	resolver.RetryTimes = 1

	ips, err := resolver.LookupHost(n.host)
	result := n.getEmptyResult()

	if err != nil {
		result.Status = StatusCritical
		result.Message = err.Error()
	} else {
		for _, ip := range ips {
			if ipString := ip.String(); ipString == n.GetExpected() {
				result.Status = StatusOK
				result.Value = ipString
			}
		}
		if result.Status != StatusOK {
			result.Status = StatusCritical
			result.Value = ipListToString(ips)
		}
	}
	n.SetResult(result)
	return result
}

func (n *NetDNSCheck) Parse(configMap ConfigMap) {
	n.ParseBaseData(configMap)

	host, err := configMap.GetString("host")
	if err == nil {
		n.host = host
	}
}

func ipListToString(ips []net.IP) string {
	results := []string{}
	for _, ip := range ips {
		results = append(results, ip.String())
	}
	return strings.Join(results, ",")
}
