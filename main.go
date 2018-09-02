package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"gopkg.in/alecthomas/kingpin.v2"
	"time"
)

var (
	exitAfterPeriod = kingpin.Flag("exitAfter", "Fail after seconds").Default("-1").Short('e').Int()
	readyAfterPeriod = kingpin.Flag("readyAfter", "Ready after seconds").Default("-1").Short('r').Int()
	liveAfterPeriod = kingpin.Flag("liveAfter", "Live after seconds").Default("-1").Short('l').Int()
)

type activeStates struct {
	Ready bool
	Live bool
}

var active = activeStates{
	Ready:false,
	Live:false,
}

func main() {
	kingpin.Parse()

	http.HandleFunc("/healthcheck", healthcheck)
	http.HandleFunc("/respond", respond())
	http.HandleFunc("/ready", readiness)
	http.HandleFunc("/live", liveness)

	if *readyAfterPeriod != -1 {
		go func() {
			fmt.Printf("Will be ready after %d seconds\n", *readyAfterPeriod)
			time.Sleep(time.Duration(*readyAfterPeriod)*time.Second)
			active.Ready = true

		}()
	}

	if *liveAfterPeriod != -1 {
		go func() {
			fmt.Printf("Will be live after %d seconds\n", *liveAfterPeriod)
			time.Sleep(time.Duration(*liveAfterPeriod)*time.Second)
			active.Live = true

		}()
	}

	if *exitAfterPeriod != -1 {
		go exitAfter(*exitAfterPeriod)
	}

	log.Println("Started server")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}


func exitAfter (waitTime int) {
	fmt.Printf("Will exit after %d seconds\n", waitTime)
	time.Sleep(time.Duration(waitTime + 1)*time.Second)
	fmt.Printf("Exiting after specified %d seconds\n", waitTime)
	os.Exit(0)
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func liveness(w http.ResponseWriter, r *http.Request) {
	if active.Live {
		fmt.Fprintf(w, "OK")
	} else {
		w.WriteHeader(500)
		fmt.Fprint(w, "Not yet")
	}
}

func readiness(w http.ResponseWriter, r *http.Request) {
	if active.Ready {
		fmt.Fprintf(w, "OK")
	} else {
		w.WriteHeader(500)
		fmt.Fprint(w, "Not yet")
	}
}

func respond() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg = "OK"
		if msgItems, ok := r.URL.Query()["msg"]; ok {
			if len(msgItems) > 0 {
				msg = msgItems[0]
			}
		}

		var code = 200
		if codeItems, ok := r.URL.Query()["code"]; ok {
			if len(codeItems) > 0 {
				var err error
				if code, err = strconv.Atoi(codeItems[0]); err != nil {
					code = 200
				}
			}
		}

		w.WriteHeader(code)
		fmt.Fprintf(w, msg)
	}
}
