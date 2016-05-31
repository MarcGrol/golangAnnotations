package example

import (
	"fmt"
	"time"
)

// +event -> aggregate: tour
type TourCreated struct {
	Year int `json:"year"`
}

func (e TourCreated) GetUid() string {
	return fmt.Sprintf("%d", e.Year)
}

// +event -> aggregate: tour
type CyclistCreated struct {
	Year        int    `json:"year"`
	CyclistId   int    `json:"cyclistId"`
	CyclistName string `json:"cyclistName"`
	CyclistTeam string `json:"cyclistTeam"`
}

func (e CyclistCreated) GetUid() string {
	return fmt.Sprintf("%d", e.Year)
}

// +event -> aggregate: tour
type EtappeCreated struct {
	Year                 int       `json:"year"`
	EtappeId             int       `json:"etappeId"`
	EtappeDate           time.Time `json:"etappeDate"`
	EtappeStartLocation  string    `json:"etappeStartLocation"`
	EtappeFinishLocation string    `json:"etappeFinishLocation"`
	EtappeLength         int       `json:"etappeLength"`
	EtappeKind           int       `json:"etappeKind"`
}

func (e EtappeCreated) GetUid() string {
	return fmt.Sprintf("%d", e.Year)
}

// +event -> aggregate: tour
type EtappeResultsCreated struct {
	Year                     int   `json:"year"`
	LastEtappeId             int   `json:"lastEtappeId"`
	BestDayCyclistIds        []int `json:"bestDayCyclistIds"`
	BestAllrounderCyclistIds []int `json:"bestAllrounderCyclistIds"`
	BestSprinterCyclistIds   []int `json:"bestSprinterCyclistIds"`
	BestClimberCyclistIds    []int `json:"bestClimberCyclistIds"`
}

func (e EtappeResultsCreated) GetUid() string {
	return fmt.Sprintf("%d", e.Year)
}

// +event -> aggregate: gambler
type GamblerCreated struct {
	GamblerUid       string `json:"gamblerUid"`
	GamblerName      string `json:"gamblerName"`
	GamblerEmail     string `json:"gamblerEmail"`
	GamblerImageIUrl string `json:"gamblerImageIUrl"`
}

func (e GamblerCreated) GetUid() string {
	return fmt.Sprintf("%d", e.GamblerUid)
}

// +event -> aggregate: gambler
type GamblerTeamCreated struct {
	GamblerUid      string `json:"gamblerUid"`
	Year            int    `json:"year"`
	GamblerCyclists []int  `json:"gamblerCyclists"`
}

func (e GamblerTeamCreated) GetUid() string {
	return fmt.Sprintf("%d", e.GamblerUid)
}

// +event -> aggregate: news
type NewsItemCreated struct {
	Year      int       `json:"year"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Sender    string    `json:"sender"`
}

func (e NewsItemCreated) GetUid() string {
	return fmt.Sprintf("%d", e.Year)
}
