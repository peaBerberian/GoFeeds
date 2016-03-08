package main

import "fmt"
import "net/http"
import "io/ioutil"

// Website structure. Used as a config structure linking each
// website to its corresponding url
type websiteStruct struct {
	name string // name used as an id for this website
	url  string // url used for the request
}

var websites = []websiteStruct{
	{
		name: "pitchfork",
		url:  "http://www.google.fr",
	},
	{
		name: "stereogum",
		url:  "http://www.debian.org",
	},
}

func handlePitchforkBody(body []byte) (res string, err error) {
	return "", nil
}

func handleStereogumBody(body []byte) (res string, err error) {
	return "", nil
}

func main() {
	var c []chan httpResponse
	for i := 0; i < len(websites); i++ {
		c = append(c, make(chan httpResponse))
		go fetchUrl(websites[i].url, c[i])
	}

	for i := 0; i < len(websites); i++ {
		response := <-c[i]
		fmt.Printf("received for %s\n", websites[i].name)
		switch websites[i].name {
		case "pitchfork":
			handlePitchforkBody(response.body)
		case "stereogum":
			handleStereogumBody(response.body)
		}
	}
}

// Struct of a http request response returned by fetchUrl
type httpResponse struct {
	err  error  // error: http error, response parsing error
	body []byte // complete response body
}

// Fetch url and return to the given channel the http response once the
// request is finished.
// Close the channel when everything is done
func fetchUrl(url string, c chan httpResponse) {
	defer close(c)
	resp, err := http.Get(url)
	if err != nil {
		c <- httpResponse{err: err}
		return
	}
	defer resp.Body.Close()
	body, errRead := ioutil.ReadAll(resp.Body)
	if err != nil {
		c <- httpResponse{err: errRead, body: body}
		return
	}

	c <- httpResponse{err: nil, body: body}
}
