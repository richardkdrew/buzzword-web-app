package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	// Set the host port, either from env or a default
	hostPort := os.Getenv("HOST_PORT")
	if hostPort == "" {
		hostPort = "8080"
	}

	// Add default route
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/css/", cssHandler)

	// Start listening
	log.Printf(fmt.Sprintf("(Go) Web server listening on port %s...\n", hostPort))
	log.Fatal(http.ListenAndServe(":"+hostPort, nil))
}

func defaultHandler(w http.ResponseWriter, req *http.Request) {

	// ParseForm to make form fields available
	err1 := req.ParseForm()
	if err1 != nil {
		log.Fatal("Error occurred while attempting to parse the form")
	}

	// Get the supplied category value
	category := req.PostFormValue("category")

	// Set the API URL, either from env a default
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:5001"
	}

	indexPageTemplate := template.New("index")

	path := "public/index.html"

	// Load in the template file
	indexPageTemplate, err := template.ParseFiles(path)

	// If no file was found then 404
	if err != nil {
		log.Printf(fmt.Sprintf("Could not find file (%s)...\n", path))
		w.WriteHeader(404)
		return
	}

	log.Printf(fmt.Sprintf("Calling API (%s)...\n", apiURL))

	// Get a buzzword
	buzzword := getBuzzword(apiURL, category)

	// Render the template with the supplied data context (buzzword)
	log.Printf(fmt.Sprintf("Serving file (%s)...\n", path))
	indexPageTemplate.Execute(w, buzzword)
}

func cssHandler(w http.ResponseWriter, req *http.Request) {

	// Determine the filename
	path := "public" + req.URL.Path

	// Open up the file and buffer it and write to the response
	file, err := os.Open(path)
	if err != nil {
		log.Printf(fmt.Sprintf("Could not find CSS file (%s)...\n", path))
		w.WriteHeader(404)
		return
	}

	// We should close file when done reading from it
	// Defer the closing of the file
	defer file.Close()
	bufferedReader := bufio.NewReader(file)
	log.Printf(fmt.Sprintf("Buffering output for file (%s)...\n", path))

	w.Header().Add("Content-Type", "text/css")
	bufferedReader.WriteTo(w)
}

func getBuzzword(api string, category string) BuzzwordResult {
	url := fmt.Sprintf("%s/api/buzzword", api)

	// Tack on the category (if there is one)
	if category != "" {
		url = fmt.Sprintf("%s?category=%s", url, html.EscapeString(category))
	}

	// Build the request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf(fmt.Sprintf("Had problems communicating with the API (%s)...\n ", url))
		log.Fatal("Error occurred while attempting to communicate with the API")
	}

	if err != nil {
		log.Printf(fmt.Sprintf("Had problems creating the API request (%s)...\n ", url))
		log.Fatal("Error occurred while creating the API request")
	}
	// For control over HTTP client headers, redirect policy, and other settings, create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and returns an HTTP response
	response, err := client.Do(request)
	if err != nil {
		log.Printf(fmt.Sprintf("Had problems communicating with the API (%s)...\n ", url))
		log.Fatal("Error occurred while attempting to communicate with the API")
	}

	// Caller should close response.Body when done reading from it
	// Defer the closing of the body
	defer response.Body.Close()

	var result BuzzwordResult

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		log.Printf(fmt.Sprintf("Had problems parsing the API response (%s)...\n ", response.Body))
		log.Fatal("Error occurred while parson the API response")
	}

	return result
}

type BuzzwordResult struct {
	Category string `json:"category"`
	Buzzword string `json:"buzzword"`
	APIID    string `json:"apiId"`
}
