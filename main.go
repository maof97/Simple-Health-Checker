package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var DEBUG_MODE bool

type Response_Generic struct {
	Id     string
	Result string
}

type healthStatus struct {
	Id               string
	Host             string
	Status           string
	StatusMessage    string
	SubStatus        string
	SubStatusMessage string
	TimeUnix         int
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	newLog(0, "API", "Test API was called (hello)")
	// Declare a new Person struct.
	Response := ""
	var p Person

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Got: "+Response+err.Error(), http.StatusBadRequest)
		return
	}

	// Do something with the Person struct...
	log(0, "API", "Got Person Name: "+p.Name)
	fmt.Fprintf(w, "Person: %+v", p)
}

func reportStatusHandler(w http.ResponseWriter, r *http.Request) {
	log(0, "API", "API was called to whitelist a domain by "+r.RemoteAddr)
	// Declare a new Person struct.
	var req Request_WL //...in JSON
	var resp Response_Generic
	_ = resp

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Got: "+err.Error(), http.StatusBadRequest)
		log(2, "API", "Got error decoding JSON from "+r.RemoteAddr+". Error: "+err.Error())
		return
	}

	if req.Id != "" {
		// Do something with the Person struct...
		log(1, "API", "Got new whitelist request with uid="+req.Id+" from "+r.RemoteAddr+" regarding domain: "+req.Domain_full+". Affected host: "+req.Host)
		//fmt.Fprintf(w, "Request Received: %+v", req)

		//Handle Request####
		if req.Is_verified {

		} else {
			//TODO
			log(2, "API", "Request "+req.Id+" denied because not verified.")
		}
	} else {
		log(2, "API", "Got new request but with missing uid. Ignoring.")
	}

}

func awaitHealthUpdate() {
	// Set up handlers
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/report_status", reportStatusHandler)

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := ioutil.ReadFile("lib/tls-certs/NIAN+WCP-API+Client.pem") //TODO implement path in config.yml
	if err != nil {
		log(3, "ERROR", "CertReadErr: "+err.Error())
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create the TLS Config with the CA pool and enable Client certificate validation
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	// Create a Server instance to listen on port 8443 with the TLS config
	server := &http.Server{
		Addr:      "api.nian.local:8443",
		TLSConfig: tlsConfig,
	}

	// Listen to HTTPS connections with the server certificate and wait
	log(3, "ERROR", "ListenAndServerTLSErr: "+(server.ListenAndServeTLS("lib/tls-certs/NIAN+WCP-API+Server.pem", "lib/tls-certs/NIAN+WCP-API+Server-key.pem")).Error()) //TODO implement path in config.yml
}

func main() {
	go awaitHealthUpdate()

}

func newLog(logType int, sub_logType string, msg string) {
	// TODO Improve logging
	if logType > 0 { //if _not_ INTERNAL

	}
	if DEBUG_MODE && sub_logType == "DEBUG" {
		print(msg)
	}
}

func logToElastic() {

}

func logToFile() {

}

func reportAsMail() {

}

func alertAsMail() {

}

func alertAsTelegramMessage() {

}

func alertToElastic() {

}

func doActiveHealthCheck() {

}
