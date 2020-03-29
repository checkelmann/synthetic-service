package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/keptn-sandbox/sdk-go/pkg/keptn"
)

/**
* Here are all the handler functions for the individual event
  See https://github.com/keptn/spec/blob/0.1.3/cloudevents.md for details on the payload

  -> "sh.keptn.event.configuration.change"
  -> "sh.keptn.events.deployment-finished"
  -> "sh.keptn.events.tests-finished"
  -> "sh.keptn.event.start-evaluation"
  -> "sh.keptn.events.evaluation-done"
  -> "sh.keptn.event.problem.open"
  -> "sh.keptn.events.problem"
*/

type SyntheticLocations struct {
	Locations []struct {
		Name     string `json:"name"`
		EntityID string `json:"entityId"`
		Type     string `json:"type"`
		Status   string `json:"status"`
	} `json:"locations"`
}

type SyntheticMonitors struct {
	Monitors []struct {
		Name     string `json:"name"`
		EntityID string `json:"entityId"`
		Type     string `json:"type"`
		Enabled  bool   `json:"enabled"`
	} `json:"monitors"`
}

//
// Handles ConfigurationChangeEventType = "sh.keptn.event.configuration.change"
// TODO: add in your handler code
//
func HandleConfigurationChangeEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ConfigurationChangeEventData) error {
	log.Printf("Handling Configuration Changed Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles DeploymentFinishedEventType = "sh.keptn.events.deployment-finished"
// TODO: add in your handler code
//
func HandleDeploymentFinishedEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.DeploymentFinishedEventData) error {
	log.Printf("Handling Deployment Finished Event: %s", incomingEvent.Context.GetID())

	// Create http Client
	client := &http.Client{}

	var locationsObject SyntheticLocations
	var synteticsMonitorsObject SyntheticMonitors

	// Get Dynatrace Secrets
	dtTenant := os.Getenv("DT_TENANT")
	dtAPItoken := os.Getenv("DT_API_TOKEN")

	if dtAPItoken == "" || dtTenant == "" {
		log.Println("No Dynatrace Credentials found!")
		return nil
	} else if data.DeploymentURIPublic == "" {
		log.Println("No DeploymentURIPublic found!")
		return nil
	}

	log.Println(data.DeploymentURIPublic)
	log.Println(data.Labels)
	log.Println(dtTenant)

	// Get Private Synthetic Check Locations
	dtAPIUrl := "https://" + dtTenant + "/api/v1/synthetic/locations?type=PRIVATE"
	req, err := http.NewRequest("GET", dtAPIUrl, nil)
	req.Header.Set("Authorization", "Api-Token "+dtAPItoken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(data, &locationsObject)
	}
	defer resp.Body.Close()
	var CheckLocations []string
	log.Println(locationsObject)
	for _, location := range locationsObject.Locations {
		log.Println(location.EntityID)
		CheckLocations = append(CheckLocations, location.EntityID)
	}

	// Check if Synthetic is existing
	// /monitors?tag=keptn_service:${APPLICATION_SHORT_NAME}&tag=keptn_stage:${CI_ENVIRONMENT_SLUG}
	dtAPIUrl = "https://" + dtTenant + "/api/v1/synthetic/monitors?tag=keptn_service:" + data.Service + "&tag=keptn_stage:" + data.Stage + "&tag=keptn_project:" + data.Project
	log.Println(dtAPIUrl)
	req, err = http.NewRequest("GET", dtAPIUrl, nil)
	req.Header.Set("Authorization", "Api-Token "+dtAPItoken)

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(data, &synteticsMonitorsObject)
	}
	defer resp.Body.Close()
	log.Println("Existing Monitors")
	log.Println(synteticsMonitorsObject)
	var monitorEntityID = ""
	var monitorMethod = "POST"
	for _, monitor := range synteticsMonitorsObject.Monitors {
		log.Println("Found existing Monitor " + monitor.EntityID)
		monitorEntityID = monitor.EntityID
		monitorMethod = "PUT"
	}

	var jsonPayload = `
		{
		"name": "` + data.Project + ` - ` + data.Service + ` - ` + data.Stage + `",
		"frequencyMin": 1,
		"enabled": true,
		"type": "HTTP",
		"script": {
		  "version": "1.0",
		  "requests": [
			{
			  "description": "` + data.Project + ` - ` + data.Service + ` - ` + data.Stage + `",
			  "url": "` + data.DeploymentURIPublic + `",
			  "method": "GET",
			  "requestBody": "",
			  "validation": {
				"rules": [
				  {
					"value": ">=400",
					"passIfFound": false,
					"type": "httpStatusesList"
				  }
				],
				"rulesChaining": "or"
			  },      
			  "configuration": {
				"acceptAnyCertificate": true,
				"followRedirects": true
			  },
			  "preProcessingScript": "",
			  "postProcessingScript": ""
			}
		  ]
		},
		"locations": [
			"` + strings.Join(CheckLocations, "\", \"") + `"
		],
		"anomalyDetection": {
		  "outageHandling": {
			"globalOutage": true,
			"localOutage": false,
			"localOutagePolicy": {
			  "affectedLocations": 1,
			  "consecutiveRuns": 3
			}
		  },
		  "loadingTimeThresholds": {
			"enabled": false,
			"thresholds": [
			  {
				"type": "TOTAL",
				"valueMs": 10000
			  }
			]
		  }
		},
		"tags": [
			"keptn_stage:` + data.Stage + `",
			"keptn_project:` + data.Project + `",
			"keptn_service:` + data.Service + `"	
		],
		"manuallyAssignedApps": [
		  "APPLICATION-EA7C4B59F27D43EB"
		]
		}
	`

	log.Println(jsonPayload)

	var jsonStr = []byte(jsonPayload)

	if monitorEntityID != "" {
		dtAPIUrl = "https://" + dtTenant + "/api/v1/synthetic/monitors/" + monitorEntityID
		req, err = http.NewRequest(monitorMethod, dtAPIUrl, bytes.NewBuffer(jsonStr))

	} else {
		dtAPIUrl = "https://" + dtTenant + "/api/v1/synthetic/monitors"
		req, err = http.NewRequest(monitorMethod, dtAPIUrl, bytes.NewBuffer(jsonStr))
	}

	req.Header.Set("Authorization", "Api-Token "+dtAPItoken)
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(data))
	}
	defer resp.Body.Close()
	// Create Monitor Object
	/*var Monitor SyntheticMonitor

	Monitor.Enabled = true
	Monitor.FrequencyMin = 5
	Monitor.Name = data.Project + " - " + data.Service + " - " + data.Stage
	Monitor.Type = "HTTP"
	Monitor.AnomalyDetection.OutageHandling.GlobalOutage = true
	Monitor.AnomalyDetection.OutageHandling.LocalOutage = false
	Monitor.AnomalyDetection.OutageHandling.LocalOutagePolicy.AffectedLocations = 1
	Monitor.AnomalyDetection.OutageHandling.LocalOutagePolicy.ConsecutiveRuns = 3
	Monitor.Locations = CheckLocations
	Monitor.Script.Version = "1.0"
	Monitor.Script.Requests = []SyntheticMonitor.Script.Requests{
		SyntheticMonitor.Script.Requests {
			Description: "Test"
		}
	}*/

	/*
			Requests []struct {
				Description   string `json:"description"`
				URL           string `json:"url"`
				Method        string `json:"method"`
				RequestBody   string `json:"requestBody"`
				Configuration struct {
					AcceptAnyCertificate bool `json:"acceptAnyCertificate"`
					FollowRedirects      bool `json:"followRedirects"`
				} `json:"configuration"`
				PreProcessingScript  string `json:"preProcessingScript"`
				PostProcessingScript string `json:"postProcessingScript"`
			} `json:"requests"`

		jsonValue, err := json.Marshal(Monitor)
		if err != nil {
			log.Println(err)
			return nil
		}

		log.Println(string(jsonValue))
	*/
	return nil
}

//
// Handles TestsFinishedEventType = "sh.keptn.events.tests-finished"
// TODO: add in your handler code
//
func HandleTestsFinishedEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.TestsFinishedEventData) error {
	log.Printf("Handling Tests Finished Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles EvaluationDoneEventType = "sh.keptn.events.evaluation-done"
// TODO: add in your handler code
//
func HandleStartEvaluationEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.StartEvaluationEventData) error {
	log.Printf("Handling Start Evaluation Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles DeploymentFinishedEventType = "sh.keptn.events.deployment-finished"
// TODO: add in your handler code
//
func HandleEvaluationDoneEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.EvaluationDoneEventData) error {
	log.Printf("Handling Evaluation Done Event: %s", incomingEvent.Context.GetID())

	return nil
}

//
// Handles ProblemOpenEventType = "sh.keptn.event.problem.open"
// Handles ProblemEventType = "sh.keptn.events.problem"
// TODO: add in your handler code
//
func HandleProblemEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ProblemEventData) error {
	log.Printf("Handling Problem Event: %s", incomingEvent.Context.GetID())

	return nil
}
