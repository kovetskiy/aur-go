package aur

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/kovetskiy/lorg"
	"github.com/seletskiy/hierr"
)

// Package represents information about uploaded package.
type Package struct {
	ID             int64       `json:"ID"`
	Name           string      `json:"Name"`
	PackageBaseID  int64       `json:"PackageBaseID"`
	PackageBase    string      `json:"PackageBase"`
	Version        string      `json:"Version"`
	Description    string      `json:"Description"`
	URL            string      `json:"URL"`
	NumVotes       int64       `json:"NumVotes"`
	Popularity     int64       `json:"Popularity"`
	OutOfDate      interface{} `json:"OutOfDate"`
	Maintainer     string      `json:"Maintainer"`
	FirstSubmitted int64       `json:"FirstSubmitted"`
	LastModified   int64       `json:"LastModified"`
	URLPath        string      `json:"URLpath"`
	MakeDepends    []string    `json:"MakeDepends"`
	Licenses       []string    `json:"License"`
	Keywords       []string    `json:"Keywords"`
}

type responseInfo struct {
	Packages []Package `json:"results"`
}

var (
	aurBaseURL = "https://aur.archlinux.org/rpc/"
	aurVersion = "5"

	useragent = "aur-go"
)

var (
	client = new(http.Client)
	logger = lorg.NewDiscarder()
	debug  = false
)

// SetUserAgent which will be used for requests to AUR API.
func SetUserAgent(value string) {
	useragent = value
}

// SetLogger that will be used for debug messages.
func SetLogger(new lorg.Logger) {
	logger = new

	debug = true
}

// GetPackages allows to retrieve information about specified packages.
func GetPackages(name ...string) (map[string]Package, error) {
	var response responseInfo

	err := call("info", url.Values{"arg[]": name}, &response)
	if err != nil {
		return nil, err
	}

	packages := map[string]Package{}
	for _, pkg := range response.Packages {
		packages[pkg.Name] = pkg
	}

	return packages, nil
}

func call(method string, value url.Values, response interface{}) error {
	debugf("~> %s %s", method, value)

	url := aurBaseURL + "?version=" + aurVersion +
		"&type=" + method +
		"&" + value.Encode()

	payload, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return hierr.Errorf(
			err,
			"can't create http request to AUR API",
		)
	}

	payload.Header.Add("User-Agent", useragent)

	resource, err := client.Do(payload)
	if err != nil {
		return hierr.Errorf(
			err, "request to AUR API failed: GET %s", url,
		)
	}

	body, err := ioutil.ReadAll(resource.Body)
	if err != nil {
		return hierr.Errorf(
			err,
			"can't read AUR API response",
		)
	}

	debugf("<~ %s", resource.Status)

	if debug {
		debugf("<~ %s", string(body))
	}

	err = json.Unmarshal(body, response)
	if err != nil {
		return hierr.Errorf(
			err, "can't decode AUR API response",
		)
	}

	return nil
}
