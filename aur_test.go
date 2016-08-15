package aur

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kovetskiy/lorg"
	"github.com/stretchr/testify/assert"
)

func init() {
	log := lorg.NewLog()
	log.SetLevel(lorg.LevelDebug)
	log.SetIndentLines(true)
	SetLogger(log)

}

func TestGetInfo_RetrieveInformationAboutPackages(t *testing.T) {
	test := assert.New(t)

	server := httptest.NewServer(http.HandlerFunc(
		func(writer http.ResponseWriter, request *http.Request) {
			test.Equal(
				"GET",
				request.Method,
			)
			test.Equal(
				"/rpc/?version=5&type=info&arg%5B%5D=foo&arg%5B%5D=bar",
				request.RequestURI,
			)

			writer.Write([]byte(`{
  "version":5,
  "type":"multiinfo",
  "resultcount":1,
  "results":[{
      "ID":10,
      "Name":"package",
      "PackageBaseID":20,
      "PackageBase":"stacket",
      "Version":"1-1",
      "Description":"description",
      "URL":"sss",
      "NumVotes":30,
      "Popularity":40,
      "OutOfDate":null,
      "Maintainer":"kovetskiy",
      "FirstSubmitted":50,
      "LastModified":60,
      "URLPath":"/cgit/aur.git/snapshot/stacket.tar.gz",
      "MakeDepends":["go","git"],
      "Keywords":["x","y"],
      "License":["GPL"]
  }]
}`,
			))
		}),
	)
	defer server.Close()

	aurBaseURL = server.URL + "/rpc/"

	packages, err := GetPackages("foo", "bar")
	test.NoError(err)
	if test.Contains(packages, "package") {
		test.EqualValues(Package{
			ID:             10,
			Name:           "package",
			PackageBaseID:  20,
			PackageBase:    "stacket",
			Version:        "1-1",
			Description:    "description",
			URL:            "sss",
			NumVotes:       30,
			Popularity:     40,
			OutOfDate:      nil,
			Maintainer:     "kovetskiy",
			FirstSubmitted: 50,
			LastModified:   60,
			URLPath:        "/cgit/aur.git/snapshot/stacket.tar.gz",
			MakeDepends:    []string{"go", "git"},
			Licenses:       []string{"GPL"},
			Keywords:       []string{"x", "y"},
		}, packages["package"])
	}
}
