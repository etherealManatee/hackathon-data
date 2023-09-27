package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	// assumes per page is always 9
	pages, remainder := calculatePages(Response.Meta)
	fmt.Println(pages, remainder)
	allHackathons := make([]Hackathon, 0)

	populateHackathons(Response.Hackathons, &allHackathons, pages, remainder)

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
	fmt.Println(allHackathons)

	// 	file, err := os.Create("devpost.csv")
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// 	defer file.Close()

	// csvWriter := csv.NewWriter(file)
	//
	//	for _, hackathon := range Response.Hackathons {
	//		row := []string{
	//			hackathon.Title,
	//			hackathon.OrganisationName,
	//			hackathon.Date,
	//			hackathon.DisplayedLocation.Location,
	//			hackathon.Url,
	//		}
	//		_ = csvWriter.Write(row)
	//	}
	//
	// csvWriter.Flush()
}

func calculatePages(meta Meta) (pages int, remainder int) {
	return meta.Count / meta.PerPage, meta.Count % meta.PerPage
}

func populateHackathons(rawHackathon []HackathonData, allHackathon *[]Hackathon, pages int, remainder int) {
	totalPages := pages
	if remainder != 0 {
		totalPages += 1
	}
	fmt.Println(totalPages)
	// apiUrl := fmt.Sprintf("https://devpost.com/api/hackathons?page=%v&status[]=upcoming", page)
}
