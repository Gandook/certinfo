package certinfo

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
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
	// ReadFromURL retrieves useful information about a DigSig X.509 certificate from
	// a certain URL.
	// daid is the certificate's DAID.
	// cid is the certificate's CID.
	ReadFromURL(daid string, cid string) (CertInfo, error)

	// ReadFromFile retrieves useful information about a DigSig X.509 certificate from
	// a certificate file.
	// filename is the path to the certificate file.
	ReadFromFile(filename string) (CertInfo, error)
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

// readContent reads all the content from a reader and returns it as a string.
func readContent(r io.Reader) (string, error) {
	content, readErr := io.ReadAll(r)
	return string(content), readErr
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

	if resp.StatusCode == http.StatusBadRequest {
		return "", errBadRequest
	}
	if resp.StatusCode == http.StatusNotFound {
		return "", errNotFound
	}

	body, readErr := readContent(resp.Body)

	return body, readErr
}

// ReadFromURL retrieves useful information about a DigSig X.509 certificate from
// a certain URL.
// daid is the certificate's DAID.
// cid is the certificate's CID.
func (d defaultInfoRetriever) ReadFromURL(daid string, cid string) (CertInfo, error) {
	url := constructURL(daid, cid)

	body, getErr := getBody(url)
	if getErr != nil {
		return CertInfo{}, getErr
	}

	info, retrieveErr := retrieve(body)
	if retrieveErr != nil {
		return info, retrieveErr
	}

	info.CID = cid
	info.DAID = daid

	return info, nil
}

// ReadFromFile retrieves useful information about a DigSig X.509 certificate from
// a certificate file.
// filename is the path to the certificate file.
func (d defaultInfoRetriever) ReadFromFile(filename string) (CertInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return CertInfo{}, err
	}
	defer func(file *os.File) {
		closingErr := file.Close()
		if closingErr != nil {
			log.Fatal(closingErr)
		}
	}(file)

	content, readErr := io.ReadAll(file)
	if readErr != nil {
		return CertInfo{}, readErr
	}

	info, retrieveErr := retrieve(string(content))
	if retrieveErr != nil {
		return info, retrieveErr
	}

	info.CID = "Not available"
	info.DAID = "Not available"

	return info, nil
}

// retrieve retrieves the useful information about a DigSig X.509 certificate.
// content contains the (leaf) certificate along with its corresponding intermediate and
// root certificates.
func retrieve(content string) (CertInfo, error) {
	subjectMatches := subjectPattern.FindAllStringSubmatch(content, 2)
	notBeforeMatch := notBeforePattern.FindStringSubmatch(content)
	notAfterMatch := notAfterPattern.FindStringSubmatch(content)

	issuer := subjectMatches[1][1]
	subject := subjectMatches[0][1]
	notBefore := notBeforeMatch[1] + " " + notBeforeMatch[2]
	notAfter := notAfterMatch[1] + " " + notAfterMatch[2]

	return CertInfo{
		CID:       "", // It will be set in the ReadFromURL function.
		DAID:      "", // It will be set in the ReadFromURL function.
		Issuer:    issuer,
		Subject:   subject,
		NotBefore: notBefore,
		NotAfter:  notAfter,
	}, nil
}
