package main

import (
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"strconv"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
	"github.com/mackerelio/golib/logging"
)

// SolrDIHPlugin for mackerelplugin
type SolrDIHPlugin struct {
	Prefix string
	URL    string
	Core   string
}

// SolrDIHStatus for JSON parsing
type SolrDIHStatus struct {
	Status         string `json:"status"`
	StatusMessages struct {
		DataSourceRequests string `json:"Total Requests made to DataSource"`
		FetchedRows        string `json:"Total Rows Fetched"`
		ProcessedDocuments string `json:"Total Documents Processed"`
		SkippedDocuments   string `json:"Total Documents Skipped"`
		ElapsedTime        string `json:"Time taken"`
	} `json:"statusMessages"`
}

var logger = logging.GetLogger("metrics.plugin.solrdih")

func fetchSolrDIHStatus(baseURL string, coreName string) (apiResp SolrDIHStatus, err error) {
	req, err := http.NewRequest(http.MethodGet, baseURL+"/solr/"+coreName+"/dataimport?command=status", nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", "mackerel-plugin-solrdih")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return
	}

	return
}

func convertElapsedTimeIntoSecond(elapsedTime string) (second uint64, err error) {
	if elapsedTime == "" {
		return
	}

	vals := strings.Split(strings.Split(elapsedTime, ".")[0], ":")

	h, err := strconv.ParseUint(vals[0], 10, 64)
	if err != nil {
		return
	}

	m, err := strconv.ParseUint(vals[1], 10, 64)
	if err != nil {
		return
	}

	s, err := strconv.ParseUint(vals[2], 10, 64)
	if err != nil {
		return
	}

	second = h*60*60 + m*60 + s

	return
}

// FetchMetrics interface for mackerelplugin
func (sd SolrDIHPlugin) FetchMetrics() (map[string]interface{}, error) {
	apiResp, err := fetchSolrDIHStatus(sd.URL, sd.Core)
	if err != nil {
		return nil, err
	}

	switch apiResp.Status {
	case "idle":
		var sec uint64
		sec, err = convertElapsedTimeIntoSecond(apiResp.StatusMessages.ElapsedTime)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"total_count.data_source_requests": apiResp.StatusMessages.DataSourceRequests,
			"total_count.fetched_rows":         apiResp.StatusMessages.FetchedRows,
			"total_count.processed_documents":  apiResp.StatusMessages.ProcessedDocuments,
			"total_count.skipped_documents":    apiResp.StatusMessages.SkippedDocuments,
			"elapsed_time.sec":                 sec,
		}, nil
	default:
		var zero uint64
		return map[string]interface{}{
			"total_count.data_source_requests": zero,
			"total_count.fetched_rows":         zero,
			"total_count.processed_documents":  zero,
			"total_count.skipped_documents":    zero,
			"elapsed_time.sec":                 zero,
		}, nil
	}
}

// GraphDefinition interface for mackerelplugin
func (sd SolrDIHPlugin) GraphDefinition() map[string]mp.Graphs {
	return map[string]mp.Graphs{
		"total_count": {
			Label: "Solr DIH Total Count",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "data_source_requests", Label: "DataSourceRequests", AbsoluteName: true},
				{Name: "fetched_rows", Label: "FetchedRows", AbsoluteName: true},
				{Name: "processed_documents", Label: "ProcessedDocuments", AbsoluteName: true},
				{Name: "skipped_documents", Label: "SkippedDocuments", AbsoluteName: true},
			},
		},
		"elapsed_time": {
			Label: "Solr DIH Elapsed Time",
			Unit:  "integer",
			Metrics: []mp.Metrics{
				{Name: "sec", Label: "ElapsedTime", AbsoluteName: true},
			},
		},
	}
}

// MetricKeyPrefix interface for mackerelplugin
func (sd SolrDIHPlugin) MetricKeyPrefix() string {
	return sd.Prefix
}

func main() {
	optURL := flag.String("url", "http://127.0.0.1:8983", "Solr URL")
	optCoreName := flag.String("core", "", "Solr core name")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	if *optCoreName == "" {
		logger.Errorf("Solr core name is required.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	sd := SolrDIHPlugin{Prefix: "solrdih", URL: *optURL, Core: *optCoreName}
	p := mp.NewMackerelPlugin(sd)
	p.Tempfile = *optTempfile
	p.Run()
}
