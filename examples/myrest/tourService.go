//go:build !ci
// +build !ci

package myrest

import (
	"context"
	"time"
)

//go:generate golangAnnotations -input-dir .

type Tour struct {
	Year     int       `json:"year"`
	Etappes  []Etappe  `json:"etappes"`
	Cyclists []Cyclist `json:"cyclists"`
}

type Cyclist struct {
	UID    string `json:"uid"`
	Name   string `json:"name"`
	Points int    `json:"points"`
}

type Etappe struct {
	UID            string        `json:"uid"`
	Day            time.Time     `json:"day"`
	StartLocation  string        `json:"startLocation"`
	FinishLocation string        `json:"finishLocation"`
	EtappeResult   *EtappeResult `json:"etappeResult"`
}

type EtappeResult struct {
	EtappeUID      string   `json:"etappeUid"`
	DayRankings    []string `json:"dayRankings"`
	YellowRankings []string `json:"yellowRankings"`
	ClimbRankings  []string `json:"climbRankings"`
	SprintRankings []string `json:"sprintRankings"`
}

// @RestService( path = "/api/tour", novalidation = "true" )
type TourService struct {
}

// @RestOperation( method = "GET", path = "/{year}", format = "JSON" )
func (ts TourService) getTourOnUID(c context.Context, year int) (*Tour, error) {
	return &Tour{
		Year:     2016,
		Cyclists: []Cyclist{},
		Etappes:  []Etappe{},
	}, nil
}

// @RestOperation( method = "POST", path = "/{year}/etappe", format = "JSON" )
func (ts *TourService) createEtappe(c context.Context, year int, etappe Etappe) (*Etappe, error) {
	layout := "2006-01-02"
	dateString := "2016-07-14"
	day, _ := time.Parse(layout, dateString)
	return &Etappe{
		UID:            "14",
		Day:            day,
		StartLocation:  "Paris",
		FinishLocation: "Roubaix",
	}, nil
}

// @RestOperation( method = "PUT", path = "/{year}/etappe/{etappeUID}", format = "JSON" )
func (ts *TourService) addEtappeResults(c context.Context, year int, etappeUID string, results EtappeResult) error {
	return nil
}

// @RestOperation( method = "POST", path = "/{year}/cyclist", format = "JSON" )
func (ts *TourService) createCyclist(c context.Context, year int, cyclist Cyclist) (*Cyclist, error) {
	return &Cyclist{
		UID:    "42",
		Name:   "Boogerd, Michael",
		Points: 180,
	}, nil
}

// @RestOperation( method = "DELETE", path = "/{year}/cyclist/{cyclistUID}", format = "JSON" )
func (ts *TourService) markCyclistAbondoned(c context.Context, year int, cyclistUID string) error {
	return nil
}
