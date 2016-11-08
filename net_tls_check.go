package check

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"time"
)

// NetTLSCheck is doing a dns ip lookup.
type NetTLSCheck struct {
	BaseCheck

	host               string
	issuerOrganization string
	minValidDays       int
}

// Execute does a TLS/SSL check
func (n *NetTLSCheck) Execute() (result Result) {
	result = n.getEmptyResult()

	conf := &tls.Config{
		InsecureSkipVerify: false,
	}
	remote := n.host + ":443"
	conn, err := tls.Dial("tcp", remote, conf)
	if err != nil {
		return n.makeErrorResult(result, errors.New("No connection possible, remote: "+remote))
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates

	// check site cert
	siteCert := certs[0]

	issuerOrganization := siteCert.Issuer.Organization[0]
	result.Value = "Cert Issuer Organization: " + issuerOrganization + "\nSignature Algorithm: " + siteCert.SignatureAlgorithm.String() + "\nNot valid after: " + siteCert.NotAfter.Format("2006-01-02")

	if n.issuerOrganization != "" {
		if issuerOrganization != n.issuerOrganization {
			return n.makeErrorResult(result, errors.New("Bad Issuer Oranization: "+issuerOrganization))
		}

	}
	//fmt.Println(siteCert.Verify({}))

	// Check the signature algorithm, ignoring the root certificate.
	badCertAlgos := []x509.SignatureAlgorithm{x509.MD2WithRSA, x509.MD5WithRSA, x509.SHA1WithRSA, x509.DSAWithSHA1, x509.ECDSAWithSHA1}
	for _, badCertAlgo := range badCertAlgos {
		if siteCert.SignatureAlgorithm == badCertAlgo {
			return n.makeErrorResult(result, errors.New("Bad cert algorithm: "+badCertAlgo.String()))
		}
	}

	// Check the signature valid time
	if time.Now().AddDate(0, 0, n.minValidDays).After(siteCert.NotAfter) {
		return n.makeErrorResult(result, fmt.Errorf("Cert is not valid for more than %d days, ends %s", n.minValidDays, siteCert.NotAfter.Format("2006-01-02")))
	}

	result.Status = StatusOK
	n.SetResult(result)
	return
}

func (n *NetTLSCheck) Parse(configMap ConfigMap) {
	n.ParseBaseData(configMap)

	host, err := configMap.GetString("host")
	if err == nil {
		n.host = host
	}

	issuerOrganization, err := configMap.GetString("issuer_organization")
	if err == nil {
		n.issuerOrganization = issuerOrganization
	}

	minValidDays, err := configMap.GetInt("min_valid_days")
	if err == nil {
		n.minValidDays = minValidDays
	}
}

func (n *NetTLSCheck) makeErrorResult(result Result, err error) Result {
	result.Status = StatusCritical
	result.Message = err.Error()
	n.SetResult(result)
	return result
}
