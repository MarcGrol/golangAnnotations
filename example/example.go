package example

import "time"

//go:generate astTools -input-dir .

// {"Action":"Event","Data":{"Aggregate":"Tour"}}
type TourCreated struct {
	Year int `json:"year"`
}

// {"Action":"Event","Data":{"Aggregate":"Tour"}}
type CyclistCreated struct {
	Year        int    `json:"year"`
	CyclistId   int    `json:"cyclistId"`
	CyclistName string `json:"cyclistName"`
	CyclistTeam string `json:"cyclistTeam"`
}

// {"Action":"Event","Data":{"Aggregate":"Tour"}}
type EtappeCreated struct {
	Year                 int       `json:"year"`
	EtappeId             int       `json:"etappeId"`
	EtappeDate           time.Time `json:"etappeDate"`
	EtappeStartLocation  string    `json:"etappeStartLocation"`
	EtappeFinishLocation string    `json:"etappeFinishLocation"`
	EtappeLength         int       `json:"etappeLength"`
	EtappeKind           int       `json:"etappeKind"`
}

// {"Action":"Event","Data":{"Aggregate":"Tour"}}
type EtappeResultsCreated struct {
	Year                     int   `json:"year"`
	LastEtappeId             int   `json:"lastEtappeId"`
	BestDayCyclistIds        []int `json:"bestDayCyclistIds"`
	BestAllrounderCyclistIds []int `json:"bestAllrounderCyclistIds"`
	BestSprinterCyclistIds   []int `json:"bestSprinterCyclistIds"`
	BestClimberCyclistIds    []int `json:"bestClimberCyclistIds"`
}

// {"Action":"Event","Data":{"Aggregate":"Gambler"}}
type GamblerCreated struct {
	GamblerUid       string `json:"gamblerUid"`
	GamblerName      string `json:"gamblerName"`
	GamblerEmail     string `json:"gamblerEmail"`
	GamblerImageIUrl string `json:"gamblerImageIUrl"`
}

// {"Action":"Event","Data":{"Aggregate":"Gambler"}}
type GamblerTeamCreated struct {
	GamblerUid      string `json:"gamblerUid"`
	Year            int    `json:"year"`
	GamblerCyclists []int  `json:"gamblerCyclists"`
}

// {"Action":"Event","Data":{"Aggregate":"News"}}
type NewsItemCreated struct {
	Year      int       `json:"year"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Sender    string    `json:"sender"`
}
