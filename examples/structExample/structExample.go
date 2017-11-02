// +build !ci

package structExample

import (
	"fmt"
	"time"
)

//go:generate golangAnnotations -input-dir .

type Metadata struct {
	UUID          string
	Timestamp     time.Time
	EventTypeName string
}

// @JsonStruct()
// @Event( aggregate = "Tour")
type TourCreated struct {
	Year     int      `json:"year"`
	Tags     []string `json:"tags"`
	Metadata Metadata `json:"-"`
}

func (t TourCreated) GetUID() string {
	return fmt.Sprintf("%d", t.Year)
}

// @JsonStruct()
// @Event(aggregate = "Tour")
type CyclistCreated struct {
	Year        int      `json:"year"`
	CyclistUID  string   `json:"cyclistUid"`
	CyclistName string   `json:"cyclistName"`
	CyclistTeam string   `json:"cyclistTeam"`
	Metadata    Metadata `json:"-"`
}

func (t CyclistCreated) GetUID() string {
	return fmt.Sprintf("%d", t.Year)
}

// @Event(aggregate = "Tour")
type EtappeCreated struct {
	Year                 int       `json:"year"`
	EtappeUID            string    `json:"etappeUid"`
	EtappeDate           time.Time `json:"etappeDate"`
	EtappeStartLocation  string    `json:"etappeStartLocation"`
	EtappeFinishLocation string    `json:"etappeFinishLocation"`
	EtappeLength         int       `json:"etappeLength"`
	EtappeKind           int       `json:"etappeKind"`
	Metadata             Metadata  `json:"-"`
}

func (t EtappeCreated) GetUID() string {
	return fmt.Sprintf("%d", t.Year)
}

// @Event(aggregate = "Tour")
type EtappeResultsCreated struct {
	Year                     int      `json:"year"`
	EtappeUID                string   `json:"EtappeUid"`
	BestDayCyclistIds        []string `json:"bestDayCyclistIds"`
	BestAllrounderCyclistIds []string `json:"bestAllrounderCyclistIds"`
	BestSprinterCyclistIds   []string `json:"bestSprinterCyclistIds"`
	BestClimberCyclistIds    []string `json:"bestClimberCyclistIds"`
	Metadata                 Metadata `json:"-"`
}

func (t EtappeResultsCreated) GetUID() string {
	return fmt.Sprintf("%d", t.Year)
}

// @Event(aggregate = "Gambler")
type GamblerCreated struct {
	GamblerUID       string   `json:"gamblerUid"`
	GamblerName      string   `json:"gamblerName"`
	GamblerEmail     string   `json:"gamblerEmail"`
	GamblerImageIUrl string   `json:"gamblerImageIUrl"`
	Metadata         Metadata `json:"-"`
}

func (t GamblerCreated) GetUID() string {
	return t.GamblerUID
}

// @Event(aggregate = "Gambler")
type GamblerTeamCreated struct {
	GamblerUID      string   `json:"gamblerUid"`
	Year            int      `json:"year"`
	GamblerCyclists []string `json:"gamblerCyclists"`
	Metadata        Metadata `json:"-"`
}

func (t GamblerTeamCreated) GetUID() string {
	return t.GamblerUID
}

// @Event(aggregate = "News")
type NewsItemCreated struct {
	Year              int      `json:"year"`
	Message           string   `json:"message"`
	Sender            string   `json:"sender"`
	RelatedCyclistUID string   `json:"relatedCyclistUid"`
	RelatedEtappeUID  string   `json:"relatedEtappeUid"`
	Metadata          Metadata `json:"-"`
}

func (t NewsItemCreated) GetUID() string {
	return fmt.Sprintf("%d", t.Year)
}
