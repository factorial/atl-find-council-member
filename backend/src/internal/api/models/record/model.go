package record

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Request struct {
	Ref_ID int
}

type Record struct {
	COUNCIL_DIST   string
	COUNCIL_MEMBER string
}

type Result []Record

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

func WriteRecordResponse(w http.ResponseWriter, result *Result) {
	if body, err := json.Marshal(result); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	} else {
		http.Error(w, "Failed to marshal RecordResult", http.StatusInternalServerError)
	}
}

func BuildHttpRequest(request *Request) (*http.Request, error) {
	data := url.Values{}
	data.Add("wsparam[]", strconv.Itoa(request.Ref_ID))
	encodedData := data.Encode()

	url := `http://egis.atlantaga.gov/app/home/php/egisws.php`
	req, err := http.NewRequest("POST", url, strings.NewReader(encodedData))
	if err != nil {
		return nil, fmt.Errorf("Failed to build HTTP request: %v", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(encodedData)))

	return req, nil
}

func ParseHttpResponse(resp *http.Response) (*Result, error) {
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Response not OK: %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}

	result := make(Result, 0)
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Failed unmarshal response body: %v", err)
	}

	return &result, nil
}

func SubmitRequest(request *Request) (*Result, error) {
	req, err := BuildHttpRequest(request)
	if err != nil {
		return nil, fmt.Errorf("Failed to submit request: %v", err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
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
			http.Error(w, "Record request failed.", http.StatusBadRequest)
			return
		}

		result, err := SubmitRequest(request)
		if err != nil {
			http.Error(w, "Record request failed.", http.StatusInternalServerError)
			return
		}

		WriteRecordResponse(w, result)
	}
}
