package example

import "time"

//go:generate astTools -input-dir .

type Tour struct {
	Year     int
	Etappes  []Etappe
	Cyclists []Cyclist
}

type Cyclist struct {
	Name   string
	Points int
}

type Etappe struct {
	Day            time.Time
	StartLocation  string
	FinishLocation string
}

// {"Annotation":"RestService","With":{"Path":"/api"}}
type TourService struct {
}

// {"Annotation":"RestOperation","With":{"Method":"GET", "Path":"/tour/:year"}}
func (ts TourService) getTourOnYear(year int) (Tour, error) {
	return Tour{
		Year:     2016,
		Cyclists: []Cyclist{},
		Etappes:  []Etappe{},
	}, nil
}

// {"Annotation":"RestOperation","With":{"Method":"POST", "Path":"/tour/:year/etappe"}}
func (ts *TourService) addEtappe(year int, etappe Etappe) (Etappe, error) {
	dateString := "2016-07-14T11:45:26.371Z"
	day, _ := time.Parse(dateString, dateString)
	return Etappe{
		Day:            day,
		StartLocation:  "Paris",
		FinishLocation: "Roubaix",
	}, nil
}

// {"Annotation":"RestOperation","With":{"Method":"POST", "Path":"/tour/:year/cyclist"}}
func (ts *TourService) addCyclist(year int, cyclist Cyclist) (Cyclist, error) {
	return Cyclist{
		Name:   "Boogerd",
		Points: 120,
	}, nil
}
