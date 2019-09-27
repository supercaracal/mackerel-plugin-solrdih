package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testServerHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/solr/test_core/dataimport":
		fmt.Fprintf(w, fetchJSON("test/stats.json"))
	default:
		fmt.Fprintf(w, "{}")
	}
})

func fetchJSON(path string) string {
	json, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(json)
}

func TestConvertElapsedTimeIntoSecond(t *testing.T) {
	matrix := map[string]uint64{
		"1:1:1.111": 3661,
		"0:0:0":     0,
		"0:0:19.25": 19,
	}

	for input, expected := range matrix {
		actual, err := convertElapsedTimeIntoSecond(input)

		if err != nil {
			t.Fatal(err)
		}

		if actual != expected {
			t.Fatalf("Expected: %d, Actual: %d", expected, actual)
		}

	}
}

func TestGraphDefinition(t *testing.T) {
	solrDIH := SolrDIHPlugin{}
	actual := solrDIH.GraphDefinition()

	if actual["total_count"].Label != "Solr DIH Total Count" {
		t.Fatalf("Expected: %s, Actual: %s", "Solr DIH Total Count", actual["total_count"].Label)
	}
}

func TestFetchMetrics(t *testing.T) {
	testServer := httptest.NewServer(testServerHandler)
	defer testServer.Close()

	solrDIH := SolrDIHPlugin{URL: testServer.URL, Core: "test_core"}
	actual, err := solrDIH.FetchMetrics()

	if err != nil {
		t.Fatal(err)
	}

	var expected uint64
	expected = 19
	if actual["elapsed_time.sec"] != expected {
		t.Fatalf("Expected: %d, Actual: %d", expected, actual["elapsed_time.sec"])
	}
}

func TestMetricKeyPrefix(t *testing.T) {
	solrDIH := SolrDIHPlugin{Prefix: "solrdih"}
	actual := solrDIH.MetricKeyPrefix()

	if actual != "solrdih" {
		t.Fatalf("Expected: %s, Actual: %s", "solrdih", actual)
	}
}
