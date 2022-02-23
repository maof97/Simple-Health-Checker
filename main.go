package main

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Response_Generic struct {
	Id     string
	Result string
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Write "Hello, world!" to the response body
	log(0, "API", "HelloWorld of API was called")
	io.WriteString(w, "Hello, world!\n")
}

func testAPI(w http.ResponseWriter, r *http.Request) {
	log(0, "API", "Test API was called")
	// Declare a new Person struct.
	var p Person
	Response := ""

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

//#####

func API_wl(w http.ResponseWriter, r *http.Request) {
	log(0, "API", "API was called to whitelist a domain by "+r.RemoteAddr)
	// Declare a new Person struct.
	var req Request_WL //...in JSON
	var resp Response_Generic

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
			if req.Exp_time > 0 { //if not permanent...
				write_tmp_wl(req.Domain_full, req.Exp_time)
				log(1, "WHITELISTING", "API: Temp whitelisted domain: \""+req.Domain_full+"\" for "+fmt.Sprint(req.Exp_time)+" Minutes because of an API call with uid="+req.Id)
				log(1, "API", "Request "+req.Id+": Temp whitelisted domain: \""+req.Domain_full+"\" for "+fmt.Sprint(req.Exp_time)+" Minutes.")
				resp.Id = req.Id
				resp.Result = "Success"
				fmt.Fprintf(w, "%+v", resp)
			}

			if req.Exp_time == -1 { //if permanent
				//TODO
			}
		} else {
			//TODO
			log(2, "API", "Request "+req.Id+" denied because not verified.")
		}
	} else {
		log(2, "API", "Got new request but with missing uid. Ignoring.")
	}

}

//######
func HandleAPI() {
	// Set up handlers
	http.HandleFunc("/hello", helloHandler)
	//	http.HandleFunc("/testapi", testAPI)
	http.HandleFunc("/request_wl", API_wl)

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
	log(1, "INIT", "**  Starting Health Checker... **")

	// Open config file
	PWD := os.Getwd()
	_, err = ioutil.ReadFile(PWD+"/config.json") //Open Default DB
	if err != nil {
		log(2, "Fatal Error opening Config file")
		return
	} else {
		log(1,"" ,"STARTUP: Successful at opening config file.")
	

	//SQL CONNECT
	//sql_connect()

	go HandleAPI()



func log(_ , _, msg){
	// TODO Improve logging
	print(msg)
	return
}