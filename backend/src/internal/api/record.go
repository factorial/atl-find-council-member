package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type RecordRequest struct {
	Ref_ID int
}

type RecordRequestResult []Record

type Record struct {
	COUNCIL_DIST   string
	COUNCIL_MEMBER string
}

func ParseRecordRequest(reader io.Reader) (*RecordRequest, error) {
	if reader == nil {
		return nil, fmt.Errorf("Request Missing Body.")
	}

	var request RecordRequest
	if err := json.NewDecoder(reader).Decode(&request); err != nil {
		return nil, fmt.Errorf("Failed to parse request.")
	}

	return &request, nil
}

func WriteRecordRequestResponse(w http.ResponseWriter, result *RecordRequestResult) {
	if body, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	} else {
		http.Error(w, "Failed to marshal RecordRequestResult", http.StatusInternalServerError)
	}
}

func SubmitRecordRequest(request *RecordRequest) (*RecordRequestResult, error) {

	data := url.Values{}
	data.Add("wsparam[]", strconv.Itoa(request.Ref_ID))
	encodedData := data.Encode()

	req, err := http.NewRequest("POST", `http://egis.atlantaga.gov/app/home/php/egisws.php`, strings.NewReader(encodedData))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Response not OK: %v", resp.StatusCode)
	}

	result := make(RecordRequestResult, 0)

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	return &result, nil
}

func NewRecordRequestHandler(store Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := ParseRecordRequest(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := SubmitRecordRequest(request)
		WriteRecordRequestResponse(w, result)
	}
}
