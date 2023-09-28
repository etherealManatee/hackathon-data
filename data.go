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
	Hackathons []HackathonData
	Meta       Meta
}

type Meta struct {
	Count   int `json:"total_count"`
	PerPage int `json:"per_page"`
}

type HackathonData struct {
	Title             string            `json:"title"`
	DisplayedLocation DisplayedLocation `json:"displayed_location"`
	Date              string            `json:"submission_period_dates"`
	OrganisationName  string            `json:"organization_name"`
	Url               string            `json:"url"`
}

type DisplayedLocation struct {
	Location string `json:"location"`
}

type Hackathon struct {
	Title            string
	Location         string
	Date             string
	OrganisationName string
	Url              string
}

// https://devpost.com/api/hackathons?page=2&status[]=upcoming
// https://devpost.com/api/hackathons?page=2&status[]=upcoming
func main() {
	fmt.Print("Retrieving data from DEVPOST..\n")
	// get first request, mainly to get the number of pages
	response, err := http.Get("https://devpost.com/api/hackathons?status[]=upcoming")
	if err != nil {
		panic(err)
	}
	var Response Response
	if err := json.NewDecoder(response.Body).Decode(&Response); err != nil {
		log.Println(err)
		return
	}

	pages, remainder := calculatePages(Response.Meta)

	allHackathons := make([]Hackathon, 0)

	for _, hackathonRaw := range Response.Hackathons {
		h := &Hackathon{
			Title:            hackathonRaw.Title,
			Location:         hackathonRaw.DisplayedLocation.Location,
			Date:             hackathonRaw.Date,
			OrganisationName: hackathonRaw.OrganisationName,
			Url:              hackathonRaw.Url,
		}
		allHackathons = append(allHackathons, *h)
	}

	populateHackathons(Response.Hackathons, &allHackathons, pages, remainder)

	// check that hackathons retrieved is same as the number stated in their total_count
	if len(allHackathons) != Response.Meta.Count {
		log.Panicln("Hackathons retrieved does not match total_count, please check!")
	}

	// write to csv
	file, err := os.Create("devpost.csv")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)

	for _, hackathon := range allHackathons {
		row := []string{
			hackathon.Title,
			hackathon.OrganisationName,
			hackathon.Date,
			hackathon.Location,
			hackathon.Url,
		}
		_ = csvWriter.Write(row)
	}

	csvWriter.Flush()
}

func calculatePages(meta Meta) (pages int, remainder int) {
	return meta.Count / meta.PerPage, meta.Count % meta.PerPage
}

func populateHackathons(rawHackathon []HackathonData, allHackathons *[]Hackathon, pages int, remainder int) {
	totalPages := pages
	if remainder != 0 {
		totalPages += 1
	}

	// start with page 2 since page 1 was already done
	for i := 2; i < totalPages+1; i++ {
		apiUrl := fmt.Sprintf("https://devpost.com/api/hackathons?page=%v&status[]=upcoming", i)
		response, err := http.Get(apiUrl)
		if err != nil {
			panic(err)
		}
		var Response Response
		if err := json.NewDecoder(response.Body).Decode(&Response); err != nil {
			log.Println(err)
			return
		}
		for _, hackathonRaw := range Response.Hackathons {
			h := &Hackathon{
				Title:            hackathonRaw.Title,
				Location:         hackathonRaw.DisplayedLocation.Location,
				Date:             hackathonRaw.Date,
				OrganisationName: hackathonRaw.OrganisationName,
				Url:              hackathonRaw.Url,
			}
			*allHackathons = append(*allHackathons, *h)
		}
	}
}
