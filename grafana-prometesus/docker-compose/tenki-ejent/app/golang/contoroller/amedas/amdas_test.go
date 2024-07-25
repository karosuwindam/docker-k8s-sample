package amedas

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func TestGetJSON(t *testing.T) {
	ctx := context.Background()
	if ltime, err := getLastTime(ctx); err == nil {
		fmt.Println(MAP_DATA_URL + ltime + ".json")
		fmt.Println(getJsonData(ctx, MAP_DATA_URL+ltime+".json"))

		if ds, err1 := getJsonDataMapData(ctx, ltime); err1 == nil {
			if td, err2 := getJsonTable(ctx); err2 == nil {
				da := covertToePrometesusDatas(td, ds)
				fmt.Println(da)
				for i, v := range da {
					if i < 10 {
						fmt.Println(v.ToJson())
						fmt.Println(v.ToPrometesus())
					} else {
						break
					}
				}
			}
		}
	}
	api := NewAmedasAPI()
	fmt.Println(api.GetPrometesusData())
	// fmt.Println(getJsonData(ctx, AMEDAS_ELEMNT_URL))
}

func getJsonData(ctx context.Context, urldata string) (interface{}, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urldata, nil)
	// client := new(http.Client)
	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil

}
