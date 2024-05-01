package dhtsenser_test

import (
	"testing"

	"github.com/d2r2/go-dht"
	"github.com/d2r2/go-logger"
)

var lg = logger.NewPackageLogger("main",
	logger.DebugLevel,
	// logger.InfoLevel,
)

func TestDHT(t *testing.T) {
	defer logger.FinalizeLogger()

	lg.Notify("***************************************************************************************************")
	lg.Notify("*** You can change verbosity of output, to modify logging level of module \"dht\"")
	lg.Notify("*** Uncomment/comment corresponding lines with call to ChangePackageLogLevel(...)")
	lg.Notify("***************************************************************************************************")
	sensorType := dht.DHT11
	pin := 583
	temperature, humidity, retried, err :=
		dht.ReadDHTxxWithRetry(sensorType, pin, false, 10)
	if err != nil {
		lg.Fatal(err)
	}
	// print temperature and humidity
	lg.Infof("Sensor = %v: Temperature = %v*C, Humidity = %v%% (retried %d times)",
		sensorType, temperature, humidity, retried)

}
