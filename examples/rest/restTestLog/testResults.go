package restTestLog

// Generated automatically based on running of api-tests

import (
	"github.com/MarcGrol/golangAnnotations/generator/rest/testcase"
)

var TestResults = testcase.TestSuiteDescriptor{
	TestCases: []testcase.TestCaseDescriptor{
		{
			Name:        "TestGetTour",
			Description: "Should return 200",
			Operation:   "getTourOnUid",
			Request: testcase.RequestDescriptor{
				Method: "GET",
				Url:    "/api/tour/2016",
				Headers: []string{
					"Accept:application/json",
				},
			},
			Response: testcase.ResponseDescriptor{
				Status: 200,
				Headers: []string{
					"Content-Type:application/json",
				},
				Body: `{
	"year": 2016,
	"etappes": [],
	"cyclists": []
}
`,
			},
		},
		{
			Name:        "TestCreateCyclist",
			Description: "Should return 200",
			Operation:   "createCyclist",
			Request: testcase.RequestDescriptor{
				Method: "POST",
				Url:    "/api/tour/2016/cyclist",
				Headers: []string{
					"Accept:application/json",
				},
				Body: `{
	"uid": "1",
	"name": "Boogerd, Michael",
	"points": 42
}`,
			},
			Response: testcase.ResponseDescriptor{
				Status: 200,
				Headers: []string{
					"Content-Type:application/json",
				},
				Body: `{
	"uid": "42",
	"name": "Boogerd, Michael",
	"points": 180
}
`,
			},
		},
		{
			Name:        "TestCreateEtappe",
			Description: "Should return 200",
			Operation:   "createEtappe",
			Request: testcase.RequestDescriptor{
				Method: "POST",
				Url:    "/api/tour/2016/etappe",
				Headers: []string{
					"Accept:application/json",
				},
				Body: `{
	"uid": "14",
	"day": "2016-07-14T10:00:00Z",
	"startLocation": "Paris",
	"finishLocation": "Roubaix",
	"etappeResult": null
}`,
			},
			Response: testcase.ResponseDescriptor{
				Status: 200,
				Headers: []string{
					"Content-Type:application/json",
				},
				Body: `{
	"uid": "14",
	"day": "0001-01-01T00:00:00Z",
	"startLocation": "Paris",
	"finishLocation": "Roubaix",
	"etappeResult": null
}
`,
			},
		},
		{
			Name:        "TestCreateEtappeResult",
			Description: "Should return 200",
			Operation:   "addEtappeResults",
			Request: testcase.RequestDescriptor{
				Method:  "PUT",
				Url:     "/api/tour/2016/etappe/14",
				Headers: []string{},
				Body: `{
	"etappeUid": "14",
	"dayRankings": [
		"1",
		"2",
		"3"
	],
	"yellowRankings": [
		"11",
		"12",
		"13"
	],
	"climbRankings": [
		"21",
		"22",
		"23"
	],
	"sprintRankings": [
		"31",
		"32",
		"33"
	]
}`,
			},
			Response: testcase.ResponseDescriptor{
				Status:  204,
				Headers: []string{},
				Body:    ``,
			},
		},
		{
			Name:        "TestMarkCyclistAbandoned",
			Description: "Should return 200",
			Operation:   "markCyclistAbondoned",
			Request: testcase.RequestDescriptor{
				Method:  "DELETE",
				Url:     "/api/tour/2016/cyclist/42",
				Headers: []string{},
			},
			Response: testcase.ResponseDescriptor{
				Status:  204,
				Headers: []string{},
				Body:    ``,
			},
		},
	},
}
