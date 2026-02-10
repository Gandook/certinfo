package certinfo

import (
	"errors"
	"os"
	"strconv"
	"testing"
)

func TestConstructURL(t *testing.T) {
	var url string

	url = constructURL("QCDEMO", "3")
	if url != "https://idetrust.com/daid/QCDEMO/cid/3" {
		t.Errorf("Expected URL to be 'https://idetrust.com/daid/QCDEMO/cid/3,' got '%s.'", url)
	}
}

func TestGetBody(t *testing.T) {
	var url string
	var body string
	var expected string
	var getErr error

	var daid string = "QCDEMO"
	var cid string

	t.Run("Valid", func(t *testing.T) {
		for i := 1; i <= 3; i++ {
			cid = strconv.Itoa(i)
			url = constructURL(daid, cid)

			body, getErr = getBody(url)
			if getErr != nil {
				t.Errorf("Unexpected error: %s", getErr)
			}

			content, readErr := os.ReadFile("testdata/QC-DEMO-" + cid + ".pem")
			if readErr != nil {
				t.Errorf("Unexpected error: %s", readErr)
			}

			expected = string(content)
			if body != expected {
				t.Errorf("Mismatch for DAID: %s, CID: %s.", daid, cid)
			}
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		url = constructURL(daid, "0")

		_, getErr = getBody(url)
		if getErr == nil {
			t.Errorf("Expected a 'bad request' error, got no error.")
		}
		if getErr != nil && !errors.Is(getErr, errBadRequest) {
			t.Errorf("Expected a 'bad request' error, got this error instead: %s", getErr)
		}

		url = constructURL(daid, "4")

		_, getErr = getBody(url)
		if getErr == nil {
			t.Errorf("Expected a 'not found' error, got no error.")
		}
		if getErr != nil && !errors.Is(getErr, errNotFound) {
			t.Errorf("Expected a 'not found' error, got this error instead: %s", getErr)
		}
	})
}

func TestRetrieve(t *testing.T) {
	var daid string = "QCDEMO"
	var info, expected CertInfo
	var err error

	retriever := NewRetriever()

	info, err = retriever.Retrieve(daid, "1")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	expected = CertInfo{
		CID:       "1",
		DAID:      "QCDEMO",
		Issuer:    "CN=QC DigSig Demo QC-DEMO https://www.idetrust.io 2026,O=QC DigSig Demo Inc.,L=Delmenhorst,C=DE",
		Subject:   "CN=https://idetrust.com/daid/QC%20DEMO/cid/1,O=IDeTRUST GmbH,C=DE",
		NotBefore: "2026-01-08 07:34:23",
		NotAfter:  "2027-01-08 07:34:23",
	}

	if info != expected {
		t.Errorf("The retrieved info does not match for DAID: %s, CID: 1.", daid)
	}

	info, err = retriever.Retrieve(daid, "2")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	expected = CertInfo{
		CID:       "2",
		DAID:      "QCDEMO",
		Issuer:    "CN=QC DigSig Demo QC-DEMO https://www.idetrust.io 2026,O=QC DigSig Demo Inc.,L=Delmenhorst,C=DE",
		Subject:   "CN=https://idetrust.com/daid/QC%20DEMO/cid/2,O=IDeTRUST GmbH,C=DE",
		NotBefore: "2026-01-09 06:55:35",
		NotAfter:  "2027-01-09 06:55:35",
	}

	if info != expected {
		t.Errorf("The retrieved info does not match for DAID: %s, CID: 2.", daid)
	}

	info, err = retriever.Retrieve(daid, "3")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	expected = CertInfo{
		CID:       "3",
		DAID:      "QCDEMO",
		Issuer:    "CN=QC DigSig Demo QC-DEMO https://www.idetrust.io 2026,O=QC DigSig Demo Inc.,L=Delmenhorst,C=DE",
		Subject:   "CN=https://idetrust.com/daid/QC%20DEMO/cid/3,O=IDeTRUST GmbH,C=DE",
		NotBefore: "2026-01-20 05:12:21",
		NotAfter:  "2027-01-20 05:12:21",
	}

	if info != expected {
		t.Errorf("The retrieved info does not match for DAID: %s, CID: 3.", daid)
	}
}
