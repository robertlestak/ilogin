package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/robertlestak/ilogin/pkg/auth"
	log "github.com/sirupsen/logrus"
)

var (
	serverMode bool
	authUrl    string
	serviceUrl string
	cookieName string
	outFile    string
)

func init() {
	ll, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		ll = log.InfoLevel
	}
	log.SetLevel(ll)
}

func readrcfile() {
	l := log.WithFields(log.Fields{
		"app": "ilogin",
		"fn":  "readrcfile",
	})
	l.Debug("start")
	if _, err := os.Stat(os.Getenv("HOME") + "/.iloginrc"); os.IsNotExist(err) {
		l.Debug(".iloginrc does not exist")
		return
	}
	file, err := os.Open(os.Getenv("HOME") + "/.iloginrc")
	if err != nil {
		l.WithError(err).Fatal("failed to open .iloginrc")
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ln := scanner.Text()
		if ln == "" {
			continue
		}
		if ln[0] == '#' {
			continue
		}
		// split on space
		// if there are more or less than 2, skip
		spl := strings.Split(ln, " ")
		if len(spl) != 2 {
			continue
		}
		switch spl[0] {
		case "auth":
			authUrl = spl[1]
		case "service":
			serviceUrl = spl[1]
		case "cookie":
			cookieName = spl[1]
		case "out_file":
			outFile = spl[1]
		}
	}
	if err := scanner.Err(); err != nil {
		l.Fatal(err)
		return
	}
}

func server() {
	l := log.WithFields(log.Fields{
		"app": "ilogin",
		"fn":  "server",
	})
	l.Debug("start")
	if err := auth.TokenServer(); err != nil {
		l.Fatal(err)
	}
}

func client() {
	l := log.WithFields(log.Fields{
		"app": "ilogin",
		"fn":  "client",
	})
	l.Debug("start")
	readrcfile()
	if authUrl == "" {
		l.Error("auth is empty")
		os.Exit(1)
	}
	if serviceUrl == "" {
		l.Error("service is empty")
		os.Exit(1)
	}
	if cookieName == "" {
		l.Error("cookie is empty")
		os.Exit(1)
	}
	go auth.OpenAuthWindow(serviceUrl, authUrl, cookieName)
	l.Debug("waiting for auth token")
	token, err := auth.ClientCallbackServer()
	if err != nil {
		l.Error(err)
	}
	if outFile != "" {
		if strings.HasPrefix(outFile, "~/") {
			homedir, err := os.UserHomeDir()
			if err != nil {
				l.Error(err)
				os.Exit(1)
			}
			outFile = filepath.Join(homedir, outFile[2:])
		}
		if err := ioutil.WriteFile(outFile, []byte(token), 0644); err != nil {
			l.Error(err)
			os.Exit(1)
		}
	}
	fmt.Print(token)
	l.Debug("end")
}

func main() {
	l := log.WithFields(log.Fields{
		"app": "ilogin",
		"fn":  "main",
	})
	l.Debug("start")
	flag.BoolVar(&serverMode, "server", false, "run server mode")
	flag.StringVar(&authUrl, "auth", "", "auth url")
	flag.StringVar(&serviceUrl, "service", "", "service url")
	flag.StringVar(&cookieName, "cookie", "", "cookie name")
	flag.StringVar(&outFile, "f", "", "output file")
	flag.Parse()
	if serverMode {
		server()
	} else {
		client()
	}
}
