package address

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Request struct {
	Address string
}

type Result struct {
	Candidates []Candidate `json:"candidates"`
}

type Address string

type Location struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	M float64 `json:"m"`
}

type Attributes struct {
	Score  float64 `json:"score"`
	Ref_ID int     `json:"Ref_ID"`
}

type Candidate struct {
	Address    Address    `json:"address"`
	Location   Location   `json:"location"`
	Attributes Attributes `json:"attributes"`
}

func ParseRequest(reader io.Reader) (*Request, error) {
	if reader == nil {
		return nil, fmt.Errorf("Request Missing Body.")
	}

	var request Request
	if err := json.NewDecoder(reader).Decode(&request); err != nil {
		return nil, fmt.Errorf("Failed to parse request.")
	}

	return &request, nil
}

func WriteResponse(w http.ResponseWriter, result *Result) {
	if body, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	} else {
		http.Error(w, "Failed to marshal AddressResult", http.StatusInternalServerError)
	}
}

func BuildHttpRequest(request *Request) (*http.Request, error) {
	url := "https://egis.atlantaga.gov/arc/rest/services/WebLocators/TrAddrPointS/GeocodeServer/findAddressCandidates"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to build HTTP Request: %v", err)
	}

	q := req.URL.Query()
	q.Add("Single Line Input", request.Address)
	q.Add("f", "json")
	q.Add("outFields", "*")
	q.Add("outSR", `{"wkid":4326}`)
	q.Add("maxLocations", "6")

	req.URL.RawQuery = q.Encode()

	return req, nil
}

func ParseHttpResponse(resp *http.Response) (*Result, error) {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Response status not OK: %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}

	var result Result
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Failed unmarshal response body: %v", err)
	}

	return &result, nil
}

func SubmitRequest(request *Request) (*Result, error) {
	httpRequest, err := BuildHttpRequest(request)
	if err != nil {
		return nil, fmt.Errorf("Failed to submit request: %v", err)
	}

	client := http.Client{}
	resp, err := client.Do(httpRequest)
	if err != nil {
		log.Fatal(err)
	}

	result, err := ParseHttpResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("Failed to get result from HTTP response: %v", err)
	}

	return result, nil
}

func NewHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := ParseRequest(r.Body)
		if err != nil {
			http.Error(w, "Address request failed.", http.StatusBadRequest)
			return
		}

		result, err := SubmitRequest(request)
		if err != nil {
			http.Error(w, "Address request failed.", http.StatusInternalServerError)
			return
		}

		WriteResponse(w, result)
	}
}
