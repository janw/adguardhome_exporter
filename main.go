package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"time"

	"github.com.sfragata/adguardhome_exporter/collector"
	"github.com.sfragata/adguardhome_exporter/server"
	"github.com/integrii/flaggy"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// These variables will be replaced by real values when do gorelease
var (
	version = "none"
	date    string
	commit  string
)

const prog_name = "adguardhome_exporter"

func main() {
	info := fmt.Sprintf(
		"%s\nDate: %s\nCommit: %s\nOS: %s\nArch: %s",
		version,
		date,
		commit,
		runtime.GOOS,
		runtime.GOARCH,
	)

	flaggy.SetName(prog_name)
	flaggy.SetDescription("Prometheus exporter for AdGuard Home")
	flaggy.SetVersion(info)

	var adguardUrl = "http://127.0.0.1"
	if envvalue, ok := os.LookupEnv("ADGUARD_HOME_URL"); ok {
		if len(envvalue) > 0 {
			adguardUrl = envvalue
		}
	}
	flaggy.String(&adguardUrl, "u", "url", "AdGuard Home URL (env var: ADGUARD_HOME_URL)")

	var adguardUsername = os.Getenv("ADGUARD_HOME_USERNAME")
	flaggy.String(&adguardUsername, "U", "username", "AdGuard Home username (env var: ADGUARD_HOME_USERNAME)")

	var adguardPassword = os.Getenv("ADGUARD_HOME_PASSWORD")
	flaggy.String(&adguardPassword, "P", "password", "AdGuard Home password (env var: ADGUARD_HOME_PASSWORD)")

	var metricsPort = "9311"
	flaggy.String(&metricsPort, "l", "listen-address", "Exporter metrics port")

	var adguardTlsNoVerify = false
	flaggy.Bool(&adguardTlsNoVerify, "", "tls-no-verify", "Disable TLS validation")

	var adguardTimeout = 2.0
	flaggy.Float64(&adguardTimeout, "", "timeout", "Request timeout in seconds")

	flaggy.Parse()

	u, err := url.Parse(adguardUrl)
	if err != nil || u.Host == "" || (u.Scheme != "" && u.Scheme != "http" && u.Scheme != "https") {
		log.Fatalf("Invalid url: %v", adguardUrl)
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	adguardUrl = u.String()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: adguardTlsNoVerify},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(adguardTimeout) * time.Second,
	}

	AdguardServer := server.AdguardServer{
		Url:        adguardUrl,
		Username:   adguardUsername,
		Password:   adguardPassword,
		HTTPClient: *client,
	}

	err = prometheus.Register(collector.NewAdguardCollector(AdguardServer, version))
	if err != nil {
		log.Fatalf("Failed to register collectors: %v", err)
	}

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting %s [:%s] for AdGuard Home at %s", prog_name, metricsPort, adguardUrl)
	err = http.ListenAndServe(":"+metricsPort, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
