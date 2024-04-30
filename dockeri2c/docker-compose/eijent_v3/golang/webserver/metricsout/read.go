package metricsout

import (
	"eijent/controller/senser"
	"fmt"
	"net/http"
)

func GetMetrics(w http.ResponseWriter, r *http.Request) {
	tmps := senser.ReadValue()
	lines := createLineData(tmps.Data)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	for _, line := range lines {
		fmt.Fprintf(w, "%s\n", line)
	}
}

func createLineData(datas []*senser.SenserData) []string {
	var output []string
	for _, data := range datas {
		output = append(output, createLineMetrics(data.Senser, data.Type, data.Data))
	}
	return output
}

func createLineMetrics(name, types, value string) string {
	if name == "" || name == "localhost" {
		return "senser_data{type=\"" + types + "\"} " + value

	}
	return "senser_data{type=\"" + types + "\",sennser=\"" + name + "\"} " + value
}
