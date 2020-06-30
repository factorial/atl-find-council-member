package main

import (
	"backend/internal/api"

	"log"
	"os"
)

func main() {
	address := os.Args[1]
	searchAddressRequest := api.SearchAddressRequest{Address: address}
	searchAddressResult, err := api.SubmitSearchAddressRequest(&searchAddressRequest)
	if err != nil {
		log.Fatal(err)
	}

	for _, candidate := range searchAddressResult.Candidates {
		log.Printf("%f:%s:%d\n", candidate.Attributes.Score, candidate.Address, candidate.Attributes.Ref_ID)
		recordRequest := api.RecordRequest{Ref_ID: candidate.Attributes.Ref_ID}
		recordResult, err := api.SubmitRecordRequest(&recordRequest)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%v\n", recordResult)
	}
}
