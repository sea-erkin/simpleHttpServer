package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	print              = fmt.Println
	listenPortFlag     = flag.String("p", "", "-p Port to listen on")
	serveDirectoryFlag = flag.String("d", "", "(optional) -d Path to directory to serve")
	certChainPathFlag  = flag.String("c", "", "(optional) -c Path to cert chain")
	certPrivKeyFlag    = flag.String("k", "", "(optional) -k Path to cert private key")
	isTLS              = false
)

func main() {

	checkFlags()

	http.Handle("/", http.FileServer(http.Dir(*serveDirectoryFlag)))

	if isTLS {
		log.Fatal(http.ListenAndServeTLS(":"+*listenPortFlag, *certChainPathFlag, *certPrivKeyFlag, nil))
	} else {
		log.Fatal(http.ListenAndServe(":"+*listenPortFlag, nil))
	}
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
