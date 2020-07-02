package council

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Request struct {
	District string
}

type Result struct {
	Href     string `json="href"`
	District string `json="district"`
	Name     string `json="name"`
	Image    string `json="image"`
	Contact  string `json="contact"`
}

type Data map[string]Result

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
		http.Error(w, "Failed to marshal Council Result", http.StatusInternalServerError)
	}
}

func (data *Data) ProcessRequest(request *Request) (*Result, error) {
	if result, ok := (*data)[request.District]; ok {
		return &result, nil
	}

	return nil, nil
}

func NewData(datapath string) (*Data, error) {
	file, err := ioutil.ReadFile(datapath)
	if err != nil {
		log.Fatalf("Failed to open council data file: %v", err)
	}

	list := make([]Result, 0)

	if err = json.Unmarshal(file, &list); err != nil {
		log.Fatal("Failed to unmarshal council data file.")
	}

	result := make(Data)
	for idx, item := range list {
		result[item.District] = list[idx]
	}

	return &result, nil
}

func NewHandler(datapath string) http.HandlerFunc {
	data, err := NewData(datapath)
	if err != nil {
		log.Fatal("Failed to open council data source.")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		request, err := ParseRequest(r.Body)
		if err != nil {
			http.Error(w, "Council request failed.", http.StatusBadRequest)
			return
		}

		result, err := data.ProcessRequest(request)
		if err != nil {
			http.Error(w, "Council request failed.", http.StatusInternalServerError)
			return
		}

		if result != nil {
			WriteResponse(w, result)
		} else {
			http.Error(w, "Found nothing corresponding to that district key.", http.StatusNotFound)
		}
	}
}
