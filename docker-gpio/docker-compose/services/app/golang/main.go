package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/stianeikeland/go-rpio"
)

func close(pin rpio.Pin) {
	pin.Low()
	rpio.Close()
}
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
func data(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	output := ""
	if upath == "/" {
		upath += "index.html"
	}
	upath = "." + upath
	if !Exists(upath) {
		w.WriteHeader(404)
		log.Printf("ERROR request:%v\n", r.URL.Path)
		return
	} else {
		fp, err := os.Open(upath)
		if err != nil {
			log.Panic(err)
			return
		}
		defer fp.Close()
		buf := make([]byte, 1024)
		for {
			n, err := fp.Read(buf)
			if err != nil {
				break
			}
			if n == 0 {
				break
			}
			output += string(buf[:n])
		}
		fmt.Fprintf(w, "%s", output)
	}
}
func chgpio(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		gpiodata, _ = strconv.Atoi(r.FormValue("gpio"))
	} else {
		fmt.Fprintf(w, "%v", gpiodata)
	}
}

var gpiodata int

func main() {
	port := "7080"
	if str := os.Getenv("WEB_PORT"); str != "" {
		port = str
	}
	gpiodata = 2
	gpiodata_tmp := false
	err := rpio.Open()
	if err != nil {
		return
	}
	pin := rpio.Pin(21)
	pin.Output()
	go func() {
		for {
			if !gpiodata_tmp {
				pin.High()
			} else {
				pin.Low()
			}
			if gpiodata == 0 {
				gpiodata_tmp = true
			} else if gpiodata == 1 {
				gpiodata_tmp = false
			} else {
				gpiodata_tmp = !gpiodata_tmp
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()

	fmt.Println("server start :" + port)
	http.HandleFunc("/chgpio", chgpio)
	http.HandleFunc("/", data)
	http.ListenAndServe(":"+port, nil)

	pin.Low()
	rpio.Close()
}
