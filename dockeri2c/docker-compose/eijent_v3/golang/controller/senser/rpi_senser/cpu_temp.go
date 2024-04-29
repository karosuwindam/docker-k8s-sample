package rpisenser

import (
	"io/ioutil"
	"strconv"

	"github.com/pkg/errors"
)

func useIoutilReadFile(fileName string) (string, error) {
	bytes, err := ioutil.ReadFile(fileName)

	return string(bytes), errors.Wrapf(err, "ioutil.ReadFile(%v)", fileName)
}

func cpu_temp_read() (float64, error) {
	var out float64 = -1
	if str, err := useIoutilReadFile(CPU_TMP_PASS); err != nil {
		return out, err
	} else {
		if str[len(str)-1] == 10 {
			str = str[:len(str)-1]
		}
		f, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return out, errors.Wrapf(err, "strconv.ParseFloat(%v)", str)
		}

		f = f / 1000
		out = f
	}
	return out, nil
}
