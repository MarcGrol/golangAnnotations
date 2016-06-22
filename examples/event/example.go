package event

import "time"

//go:generate astTools -input-dir .

// @Event( aggregate = "Tour")
type TourCreated struct {
	Year int `json:"year"`
}

// @Event(aggregate="Tour")
type CyclistCreated struct {
	Year        int    `json:"year"`
	CyclistId   int    `json:"cyclistId"`
	CyclistName string `json:"cyclistName"`
	CyclistTeam string `json:"cyclistTeam"`
}

// @Event(aggregate="Tour")
type EtappeCreated struct {
	Year                 int       `json:"year"`
	EtappeId             int       `json:"etappeId"`
	EtappeDate           time.Time `json:"etappeDate"`
	EtappeStartLocation  string    `json:"etappeStartLocation"`
	EtappeFinishLocation string    `json:"etappeFinishLocation"`
	EtappeLength         int       `json:"etappeLength"`
	EtappeKind           int       `json:"etappeKind"`
}

// @Event(aggregate="Tour")
type EtappeResultsCreated struct {
	Year                     int   `json:"year"`
	LastEtappeId             int   `json:"lastEtappeId"`
	BestDayCyclistIds        []int `json:"bestDayCyclistIds"`
	BestAllrounderCyclistIds []int `json:"bestAllrounderCyclistIds"`
	BestSprinterCyclistIds   []int `json:"bestSprinterCyclistIds"`
	BestClimberCyclistIds    []int `json:"bestClimberCyclistIds"`
}

// @Event(aggregate="Gambler")
type GamblerCreated struct {
	GamblerUid       string `json:"gamblerUid"`
	GamblerName      string `json:"gamblerName"`
	GamblerEmail     string `json:"gamblerEmail"`
	GamblerImageIUrl string `json:"gamblerImageIUrl"`
}

// @Event(aggregate="Gambler")
type GamblerTeamCreated struct {
	GamblerUid      string `json:"gamblerUid"`
	Year            int    `json:"year"`
	GamblerCyclists []int  `json:"gamblerCyclists"`
}

// @Event(aggregate="News")
type NewsItemCreated struct {
	Year      int       `json:"year"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Sender    string    `json:"sender"`
}
