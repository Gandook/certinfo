package certinfo

import (
	"errors"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"
)

// subjectPattern, notBeforePattern, and notAfterPattern are regular expressions used for
// extracting the Subject, NotBefore, and NotAfter of an X.509 certificate.
var subjectPattern = regexp.MustCompile(`(?m)^Subject:\s*(.*?)\s*$`)
var notBeforePattern = regexp.MustCompile(`(?m)^NotBefore:\s*(.*)T(.*)Z\s*$`)
var notAfterPattern = regexp.MustCompile(`(?m)^NotAfter:\s*(.*)T(.*)Z\s*$`)

var errBadRequest = errors.New("bad request")
var errNotFound = errors.New("not found")

// CertInfo contains the useful retrieved information about a DigSig X.509 certificate.
type CertInfo struct {
	CID       string
	DAID      string
	Issuer    string
	Subject   string
	NotBefore string
	NotAfter  string
}

type InfoRetriever interface {
	// Retrieve retrieves the useful information about a DigSig X.509 certificate.
	// daid is the certificate's DAID.
	// cid is the certificate's CID.
	Retrieve(daid string, cid string) (CertInfo, error)
}

// defaultInfoRetriever implements the InfoRetriever interface.
type defaultInfoRetriever struct{}

// NewRetriever creates a new InfoRetriever instance.
func NewRetriever() InfoRetriever {
	return &defaultInfoRetriever{}
}

// constructURL constructs the URL for a specific DigSig X.509 certificate located
// in the IDeTRUST credential repository.
// daid is the certificate's DAID.
// cid is the certificate's CID.
func constructURL(daid string, cid string) string {
	return "https://idetrust.com/daid/" + daid + "/cid/" + cid
}

// getBody gets all the information from a certain url and returns it as a string.
func getBody(url string) (string, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, getErr := client.Get(url)
	if getErr != nil {
		return "", getErr
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	body, readErr := io.ReadAll(resp.Body)

	if string(body) == "Bad Request" {
		return "", errBadRequest
	}
	if string(body) == "Not Found" {
		return "", errNotFound
	}

	return string(body), readErr
}

// Retrieve retrieves the useful information about a DigSig X.509 certificate.
// daid is the certificate's DAID.
// cid is the certificate's CID.
func (d defaultInfoRetriever) Retrieve(daid string, cid string) (CertInfo, error) {
	url := constructURL(daid, cid)

	body, getErr := getBody(url)
	if getErr != nil {
		return CertInfo{}, getErr
	}

	subjectMatches := subjectPattern.FindAllStringSubmatch(body, 2)
	notBeforeMatch := notBeforePattern.FindStringSubmatch(body)
	notAfterMatch := notAfterPattern.FindStringSubmatch(body)

	issuer := subjectMatches[1][1]
	subject := subjectMatches[0][1]
	notBefore := notBeforeMatch[1] + " " + notBeforeMatch[2]
	notAfter := notAfterMatch[1] + " " + notAfterMatch[2]

	return CertInfo{
		CID:       cid,
		DAID:      daid,
		Issuer:    issuer,
		Subject:   subject,
		NotBefore: notBefore,
		NotAfter:  notAfter,
	}, nil
}
