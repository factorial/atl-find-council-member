package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type SearchAddressRequest struct {
	Address string
}

type SearchAddressCandidateAddress string

type SearchAddressCandidateLocation struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	M float64 `json:"m"`
}

type SearchAddressCandidateAttributes struct {
	Score  float64 `json:"score"`
	Ref_ID int     `json:"Ref_ID"`
}

type SearchAddressCandidate struct {
	Address    SearchAddressCandidateAddress    `json:"address"`
	Location   SearchAddressCandidateLocation   `json:"location"`
	Attributes SearchAddressCandidateAttributes `json:"attributes"`
}

type SearchAddressResult struct {
	Candidates []SearchAddressCandidate `json:"candidates"`
}

func ParseSearchAddressRequest(reader io.Reader) (*SearchAddressRequest, error) {
	if reader == nil {
		return nil, fmt.Errorf("Request Missing Body.")
	}

	var request SearchAddressRequest
	if err := json.NewDecoder(reader).Decode(&request); err != nil {
		return nil, fmt.Errorf("Failed to parse request.")
	}

	return &request, nil
}

func WriteSearchAddressResponse(w http.ResponseWriter, result *SearchAddressResult) {
	if body, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	} else {
		http.Error(w, "Failed to marshal SearchAddressResult", http.StatusInternalServerError)
	}
}

func SubmitSearchAddressRequest(request *SearchAddressRequest) (*SearchAddressResult, error) {

	req, err := http.NewRequest("GET", "https://egis.atlantaga.gov/arc/rest/services/WebLocators/TrAddrPointS/GeocodeServer/findAddressCandidates", nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Submitting request for address: %s\n", request.Address)

	q := req.URL.Query()
	q.Add("Single Line Input", request.Address)
	q.Add("f", "json")
	q.Add("outFields", "*")
	q.Add("outSR", `{"wkid":4326}`)
	q.Add("maxLocations", "6")

	req.URL.RawQuery = q.Encode()

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Status not OK: %v", resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	log.Printf("Raw body: %s\n", string(body))

	var result SearchAddressResult
	json.Unmarshal(body, &result)

	return &result, nil
}

func NewSearchAddressHandler(store Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := ParseSearchAddressRequest(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := SubmitSearchAddressRequest(request)
		WriteSearchAddressResponse(w, result)
	}
}
