package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	print              = fmt.Println
	listenPortFlag     = flag.String("p", "", "-p Port to listen on")
	serveDirectoryFlag = flag.String("d", "", "(optional) -d Path to directory to serve")
	certChainPathFlag  = flag.String("c", "", "(optional) -c Path to cert chain")
	certPrivKeyFlag    = flag.String("k", "", "(optional) -k Path to cert private key")
	noCacheFlag        = flag.Bool("n", false, "(optional) -n No cache flag")
	isTLS              = false
)

func main() {

	checkFlags()

	http.Handle("/", mainHandler())

	if isTLS {
		log.Fatal(http.ListenAndServeTLS(":"+*listenPortFlag, *certChainPathFlag, *certPrivKeyFlag, nil))
	} else {
		log.Fatal(http.ListenAndServe(":"+*listenPortFlag, nil))
	}
}

func mainHandler() http.Handler {
	if *noCacheFlag {
		return noCache(http.FileServer(http.Dir(*serveDirectoryFlag)))
	} else {
		return http.FileServer(http.Dir(*serveDirectoryFlag))
	}
}

func noCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("X-Accel-Expires", "0")
		h.ServeHTTP(w, r)
	})
}

func checkFlags() error {

	flag.Parse()
	if *listenPortFlag == "" {
		print("[INFO] No listen port provided, setting listen port to 80")
		*listenPortFlag = "80"
	}

	if *certChainPathFlag != "" {
		_, err := os.Stat(*certChainPathFlag)
		if err != nil {
			return errors.New("[ERROR] Cert chain path invalid")
		}
	}

	if *certPrivKeyFlag != "" {
		_, err := os.Stat(*certPrivKeyFlag)
		if err != nil {
			return errors.New("[ERROR] Cert private key path ivalid")
		}
	}

	if *listenPortFlag == "443" && (*certChainPathFlag == "" || *certPrivKeyFlag == "") {
		return errors.New("[ERROR] Provided port 443 but no certificate!")
	}

	if *certChainPathFlag != "" && *certPrivKeyFlag != "" {
		isTLS = true
	}

	return nil
}
