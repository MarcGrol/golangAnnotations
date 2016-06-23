package event

import "time"

//go:generate golangAnnotations -input-dir .

// @Event( aggregate = "Tour")
type TourCreated struct {
	Year int `json:"year"`
}

// @Event(aggregate="Tour")
type CyclistCreated struct {
	Year        int    `json:"year"`
	CyclistUid  string `json:"cyclistUid"`
	CyclistName string `json:"cyclistName"`
	CyclistTeam string `json:"cyclistTeam"`
}

// @Event(aggregate="Tour")
type EtappeCreated struct {
	Year                 int       `json:"year"`
	EtappeUid            string    `json:"etappeUid"`
	EtappeDate           time.Time `json:"etappeDate"`
	EtappeStartLocation  string    `json:"etappeStartLocation"`
	EtappeFinishLocation string    `json:"etappeFinishLocation"`
	EtappeLength         int       `json:"etappeLength"`
	EtappeKind           int       `json:"etappeKind"`
}

// @Event(aggregate="Tour")
type EtappeResultsCreated struct {
	Year                     int      `json:"year"`
	EtappeUid                string   `json:"EtappeUid"`
	BestDayCyclistIds        []string `json:"bestDayCyclistIds"`
	BestAllrounderCyclistIds []string `json:"bestAllrounderCyclistIds"`
	BestSprinterCyclistIds   []string `json:"bestSprinterCyclistIds"`
	BestClimberCyclistIds    []string `json:"bestClimberCyclistIds"`
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
	Year              int       `json:"year"`
	Timestamp         time.Time `json:"timestamp"`
	Message           string    `json:"message"`
	Sender            string    `json:"sender"`
	RelatedCyclistUid string    `json:"relatedCyclistUid"`
	RelatedEtappeUid  string    `json:"relatedEtappeUid"`
}
