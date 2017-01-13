// +build !ci

package rest

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetTour(t *testing.T) {
	testCase("TestGetTour", "Should return 200")

	respCode, tour, errResp, err := getTourOnUidTestHelper("/api/tour/2016")

	assert.NoError(t, err)
	assert.Nil(t, errResp)
	assert.Equal(t, http.StatusOK, respCode)
	assert.Equal(t, 2016, tour.Year)
	assert.Empty(t, tour.Etappes)
	assert.Empty(t, tour.Cyclists)
}

func TestCreateCyclist(t *testing.T) {

	testCase("TestCreateCyclist", "Should return 200")

	cyclist := Cyclist{
		UID:    "1",
		Name:   "Boogerd, Michael",
		Points: 42,
	}

	respCode, result, errResp, err := createCyclistTestHelper("/api/tour/2016/cyclist", cyclist)
	assert.NoError(t, err)
	assert.Nil(t, errResp)
	assert.Equal(t, http.StatusOK, respCode)
	assert.Equal(t, "42", result.UID)
	assert.Equal(t, "Boogerd, Michael", result.Name)
	assert.Equal(t, 180, result.Points)
}

func TestCreateEtappe(t *testing.T) {
	testCase("TestCreateEtappe", "Should return 200")

	date, _ := time.Parse(time.RFC3339, "2016-07-14T10:00:00Z")
	etappe := Etappe{
		UID:            "14",
		StartLocation:  "Paris",
		FinishLocation: "Roubaix",
		Day:            date,
	}

	respCode, result, errResp, err := createEtappeTestHelper("/api/tour/2016/etappe", etappe)
	assert.NoError(t, err)
	assert.Nil(t, errResp)
	assert.Equal(t, http.StatusOK, respCode)
	assert.Equal(t, "14", result.UID)
	assert.Equal(t, "Paris", result.StartLocation)
	assert.Equal(t, "Roubaix", result.FinishLocation)
}

func TestCreateEtappeResult(t *testing.T) {
	testCase("TestCreateEtappeResult", "Should return 200")

	results := EtappeResult{
		EtappeUID:      "14",
		DayRankings:    []string{"1", "2", "3"},
		YellowRankings: []string{"11", "12", "13"},
		ClimbRankings:  []string{"21", "22", "23"},
		SprintRankings: []string{"31", "32", "33"},
	}

	respCode, errResp, err := addEtappeResultsTestHelper("/api/tour/2016/etappe/14", results)
	assert.NoError(t, err)
	assert.Nil(t, errResp)
	assert.Equal(t, http.StatusNoContent, respCode)
}

func TestMarkCyclistAbandoned(t *testing.T) {
	testCase("TestMarkCyclistAbandoned", "Should return 200")

	respCode, errResp, err := markCyclistAbondonedTestHelper("/api/tour/2016/cyclist/42")

	assert.NoError(t, err)
	assert.Nil(t, errResp)
	assert.Equal(t, http.StatusNoContent, respCode)
}
