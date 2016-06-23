package web

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTour(t *testing.T) {
	respCode, tour, err := getTourOnUidTestHelper("/api/tour/2016")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respCode)
	assert.Equal(t, 2016, tour.Year)
	assert.Empty(t, tour.Etappes)
	assert.Empty(t, tour.Cyclists)
}

func TestCreatCyclist(t *testing.T) {
	cyclist := Cyclist{
		UID:    "1",
		Name:   "Boogerd, Michael",
		Points: 42,
	}

	respCode, result, err := createCyclistTestHelper("/api/tour/2016/cyclist", cyclist)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respCode)
	assert.Equal(t, "42", result.UID)
	assert.Equal(t, "Boogerd, Michael", result.Name)
	assert.Equal(t, 180, result.Points)
}

func TestCreatEtappe(t *testing.T) {
	date, _ := time.Parse(time.RFC3339, "2016-07-14T10:00:00Z")
	etappe := Etappe{
		UID:            "14",
		StartLocation:  "Paris",
		FinishLocation: "Roubaix",
		Day:            date,
	}

	respCode, result, err := createEtappeTestHelper("/api/tour/2016/etappe", etappe)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, respCode)
	assert.Equal(t, "14", result.UID)
	assert.Equal(t, "Paris", result.StartLocation)
	assert.Equal(t, "Roubaix", result.FinishLocation)
}

func TestCreatEtappeResult(t *testing.T) {
	results := EtappeResult{
		EtappeUID:      "14",
		DayRankings:    []string{"1", "2", "3"},
		YellowRankings: []string{"11", "12", "13"},
		ClimbRankings:  []string{"21", "22", "23"},
		SprintRankings: []string{"31", "32", "33"},
	}

	respCode, err := addEtappeResultsTestHelper("/api/tour/2016/etappe/14", results)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, respCode)
}

func TestMarkCyclistAbandoned(t *testing.T) {
	respCode, err := markCyclistAbondonedTestHelper("/api/tour/2016/cyclist/42")

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, respCode)
}
