package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Response struct {
	Hackathons []Hackathon
	Meta       Meta
}

type Meta struct {
	Count   int `json:"total_count"`
	PerPage int `json:"per_page"`
}

type Hackathon struct {
	Title             string            `json:"title"`
	DisplayedLocation DisplayedLocation `json:"displayed_location"`
	Date              string            `json:"submission_period_dates"`
	OrganisationName  string            `json:"organization_name"`
	Url               string            `json:"url"`
}

type DisplayedLocation struct {
	Location string `json:"location"`
}

// https://devpost.com/api/hackathons?page=2&status[]=upcoming
// https://devpost.com/api/hackathons?page=2&status[]=upcoming
func main() {
	fmt.Print("Retrieving data from DEVPOST..\n")
	response, err := http.Get("https://devpost.com/api/hackathons?status[]=upcoming")
	if err != nil {
		panic(err)
	}
	var Response Response
	if err := json.NewDecoder(response.Body).Decode(&Response); err != nil {
		log.Println(err)
		return
	}
	calculatePages(Response.Meta)

	file, err := os.Create("devpost.csv")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)
	for _, hackathon := range Response.Hackathons {
		row := []string{
			hackathon.Title,
			hackathon.OrganisationName,
			hackathon.Date,
			hackathon.DisplayedLocation.Location,
			hackathon.Url,
		}
		_ = csvWriter.Write(row)
	}
	csvWriter.Flush()
}

func calculatePages(meta Meta) {
	total := meta.Count
	perPage := meta.PerPage
	remainder := total % perPage
	fmt.Print(remainder)
}
