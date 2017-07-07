package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"strconv"
	"time"
	"encoding/json"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"
)

var sleepMilliseconds int
var blink bool
var lit bool
var led embd.DigitalPin

type LedStatus struct {
  Blinking    bool
  BlinkRate   int
  Lit	      bool
}

func SetUpLed(pin int, direction embd.Direction) (embd.DigitalPin){
	fmt.Println("Initializing Pin", pin, direction)
	led, err := embd.NewDigitalPin(pin)
	if err != nil {
		panic(err)
	}

	fmt.Println("Setting Direction", pin)
	if err := led.SetDirection(direction); err != nil {
		panic(err)
	}

	return led
}	


func main() {
	sleepMilliseconds = 500
	blink = false
	lit = false
	ledOn := false

    led = SetUpLed(68, embd.Out)
    defer led.Close()
    
	go func() {
		for {
			if(blink){
			    fmt.Println(ledOn)
			    if(ledOn) {
			        if err := led.Write(embd.High); err != nil {
		                panic(err)
                    } 
                } else {
                        if err := led.Write(embd.Low); err != nil {
	                    panic(err)
                    }
	            }
			}
		
			ledOn = !ledOn
			if(lit) {
			    if err := led.Write(embd.High); err != nil {
		        panic(err)
	        }
			}
			if(!blink && !lit) {
				if err := led.Write(embd.Low); err != nil {
		            panic(err)
	            }
			}
			time.Sleep(time.Duration(sleepMilliseconds) * time.Millisecond)
		}
    }()


	http.HandleFunc("/", handleMainPage)

	log.Println("Starting webserver on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("http.ListendAndServer() failed with %s\n", err)
	}
}

func pathToSlice(path string) []string {
	a := strings.SplitAfter(path, "/")
	var newSlice []string

	for _, s := range a {
		s = strings.Trim(s, "/")
		s = strings.ToLower(s)
		if(s != "") {
			newSlice = append(newSlice, s)
		}
	}	

	return newSlice
}
func handleApiCall(pathSlice []string, w http.ResponseWriter) {
	f := pathSlice[1]
	
	switch f {
		case "blink":
			blink = true
			lit = false
			if len(pathSlice) > 2 {
				bp, err := strconv.Atoi(pathSlice[2])
				if(err == nil) {
					sleepMilliseconds = bp
				}
			}
			sendOkMessage(w, "Led set to blink every " + strconv.Itoa(sleepMilliseconds) + " milliseconds")

		case "status":
			status := LedStatus{blink, sleepMilliseconds, lit}
			js, _ := json.Marshal(status)
			w.Header().Set("Content-Type", "application/json")
  			w.Write(js)
		case "on":
			blink = false
			lit = true
			sleepMilliseconds = 500
			sendOkMessage(w, "Led switched on")

		case "off":
			blink = false
			lit = false
			sleepMilliseconds = 500
			sendOkMessage(w, "Led switched off")
	}
}

func sendOkMessage(w http.ResponseWriter, message string) {
			w.Header().Set("Server", "REST Led")
			w.WriteHeader(200)
			w.Write([]byte(message))
}

func handleMainPage(w http.ResponseWriter, r *http.Request) {
	/*if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}*/
	

	if( r.URL.Path == "/" ) {
		fmt.Fprintf(w, "show the gui")
	} else {
		d := pathToSlice(r.URL.Path)
		m := d[0]
		switch m {
			case "api":
				handleApiCall(d, w)
		}
	}	
}
