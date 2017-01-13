// +build !ci

package event

import (
	"fmt"
	"time"
)

//go:generate golangAnnotations -input-dir .

// @JsonEnum()
type Color int

const (
	Red Color = iota
	Green
	Blue
)

// @JsonStruct()
// @Event( aggregate = "Tour")
type TourCreated struct {
	Year      int       `json:"year"`
	Tags      []string  `json:"tags"`
	Timestamp time.Time `json:"-"`
}

func (t TourCreated) GetUID() string {
	return fmt.Sprintf("%d", t.Year)
}

/*
// @Event(aggregate="Tour")
type CyclistCreated struct {
	Year        int       `json:"year"`
	CyclistUid  string    `json:"cyclistUid"`
	CyclistName string    `json:"cyclistName"`
	CyclistTeam string    `json:"cyclistTeam"`
	Timestamp   time.Time `json:"-"`
}

func (t CyclistCreated) GetUID() string {
	return fmt.Sprintf("%d", t.Year)
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
	Timestamp            time.Time `json:"-"`
}

func (t EtappeCreated) GetUID() string {
	return fmt.Sprintf("%d", t.Year)
}

// @Event(aggregate="Tour")
type EtappeResultsCreated struct {
	Year                     int       `json:"year"`
	EtappeUid                string    `json:"EtappeUid"`
	BestDayCyclistIds        []string  `json:"bestDayCyclistIds"`
	BestAllrounderCyclistIds []string  `json:"bestAllrounderCyclistIds"`
	BestSprinterCyclistIds   []string  `json:"bestSprinterCyclistIds"`
	BestClimberCyclistIds    []string  `json:"bestClimberCyclistIds"`
	Timestamp                time.Time `json:"-"`
}

func (t EtappeResultsCreated) GetUID() string {
	return fmt.Sprintf("%d", t.Year)
}

// @Event(aggregate="Gambler")
type GamblerCreated struct {
	GamblerUid       string    `json:"gamblerUid"`
	GamblerName      string    `json:"gamblerName"`
	GamblerEmail     string    `json:"gamblerEmail"`
	GamblerImageIUrl string    `json:"gamblerImageIUrl"`
	Timestamp        time.Time `json:"-"`
}

func (t GamblerCreated) GetUID() string {
	return t.GamblerUid
}

// @Event(aggregate="Gambler")
type GamblerTeamCreated struct {
	GamblerUid      string    `json:"gamblerUid"`
	Year            int       `json:"year"`
	GamblerCyclists []string  `json:"gamblerCyclists"`
	Timestamp       time.Time `json:"-"`
}

func (t GamblerTeamCreated) GetUID() string {
	return t.GamblerUid
}

// @Event(aggregate="News")
type NewsItemCreated struct {
	Year              int       `json:"year"`
	Message           string    `json:"message"`
	Sender            string    `json:"sender"`
	RelatedCyclistUid string    `json:"relatedCyclistUid"`
	RelatedEtappeUid  string    `json:"relatedEtappeUid"`
	Timestamp         time.Time `json:"-"`
}

func (t NewsItemCreated) GetUID() string {
	return fmt.Sprintf("%d", t.Year)
}
*/
