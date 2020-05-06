package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/keptn-sandbox/sdk-go/pkg/keptn"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
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

// httpMonitor Struct
type httpMonitor struct {
	EntityID                  string                `json:"entityId"`
	Name                      string                `json:"name"`
	FrequencyMin              int                   `json:"frequencyMin"`
	Enabled                   bool                  `json:"enabled"`
	Type                      string                `json:"type"`
	CreatedFrom               string                `json:"createdFrom"`
	Script                    Script                `json:"script"`
	Locations                 []string              `json:"locations"`
	AnomalyDetection          AnomalyDetection      `json:"anomalyDetection"`
	Tags                      []Tags                `json:"tags"`
	ManagementZones           []ManagementZones     `json:"managementZones"`
	AutomaticallyAssignedApps []string              `json:"automaticallyAssignedApps,omitempty"`
	ManuallyAssignedApps      []string              `json:"manuallyAssignedApps,omitempty"`
	Requests                  []httpMonitorRequests `json:"requests"`
}

// Configuration httpMonitor
type Configuration struct {
	AcceptAnyCertificate bool `json:"acceptAnyCertificate"`
	FollowRedirects      bool `json:"followRedirects"`
}

// ScriptRequests httpMonitor.Script
type ScriptRequests struct {
	Description          string        `json:"description"`
	URL                  string        `json:"url"`
	Method               string        `json:"method"`
	RequestBody          string        `json:"requestBody"`
	Configuration        Configuration `json:"configuration"`
	Validation           Validation    `json:"validation"`
	PreProcessingScript  string        `json:"preProcessingScript"`
	PostProcessingScript string        `json:"postProcessingScript"`
}

// Rules Validation
type Rules struct {
	Value       string `json:"value"`
	PassIfFound bool   `json:"passIfFound"`
	Type        string `json:"type"`
}

// Tags httpMonitor
type Tags struct {
	Source  string `json:"source"`
	Context string `json:"context"`
	Key     string `json:"key"`
}

// Validation ScriptRequests
type Validation struct {
	Rules         []Rules `json:"rules"`
	RulesChaining string  `json:"rulesChaining"`
}

// Script httpMonitor
type Script struct {
	Version  string           `json:"version"`
	Requests []ScriptRequests `json:"requests"`
}

// LocalOutagePolicy httpMonitor
type LocalOutagePolicy struct {
	AffectedLocations int `json:"affectedLocations"`
	ConsecutiveRuns   int `json:"consecutiveRuns"`
}

// OutageHandling httpMonitor
type OutageHandling struct {
	GlobalOutage      bool              `json:"globalOutage"`
	LocalOutage       bool              `json:"localOutage"`
	LocalOutagePolicy LocalOutagePolicy `json:"localOutagePolicy"`
}

// Thresholds httpMonitor
type Thresholds struct {
	Type    string `json:"type"`
	ValueMs int    `json:"valueMs"`
}

// LoadingTimeThresholds httpMonitor
type LoadingTimeThresholds struct {
	Enabled    bool         `json:"enabled"`
	Thresholds []Thresholds `json:"thresholds"`
}

// AnomalyDetection httpMonitor
type AnomalyDetection struct {
	OutageHandling        OutageHandling        `json:"outageHandling"`
	LoadingTimeThresholds LoadingTimeThresholds `json:"loadingTimeThresholds"`
}

// ManagementZones httpMonitor
type ManagementZones struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// httpMonitorRequests
type httpMonitorRequests struct {
	EntityID       string `json:"entityId"`
	Name           string `json:"name"`
	SequenceNumber int    `json:"sequenceNumber"`
}

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

	logger := keptnutils.NewLogger(incomingEvent.Context.GetID(), incomingEvent.ID(), "synthetic-service")

	// Create http Client
	client := &http.Client{}

	var locationsObject SyntheticLocations
	var synteticsMonitorsObject SyntheticMonitors

	// Get Dynatrace Secrets
	dtTenant := os.Getenv("DT_TENANT")
	dtAPItoken := os.Getenv("DT_API_TOKEN")

	if dtAPItoken == "" || dtTenant == "" {
		logger.Error("Dynatrace Credentials not found!")
		return nil
	} else if data.DeploymentURIPublic == "" {
		logger.Info("DeploymentURIPublic not found.")
		return nil
	}

	logger.Debug("DeploymentURIPublic: " + data.DeploymentURIPublic)
	logger.Debug("Using Tenant: " + dtTenant)

	var manuallyAssignedApps []string
	if v, found := data.Labels["SyntheticManuallyAssignedApp"]; found {
		logger.Info("ManuallyAssignedApps found: " + v)
		manuallyAssignedApps = strings.Split(v, ",")
		//manuallyAssignedApps = "\"" + strings.Join(tApps, "\", \"") + "\""
	}

	var SyntheticFrequency = 5
	if v, found := data.Labels["SyntheticFrequency"]; found {
		logger.Info("SyntheticFrequency found: " + v)
		sfTemp, err := strconv.Atoi(v)
		if err == nil {
			SyntheticFrequency = sfTemp
		}
	}

	// Get Private Synthetic Check Locations
	dtAPIUrl := "https://" + dtTenant + "/api/v1/synthetic/locations?type=PRIVATE"
	req, err := http.NewRequest("GET", dtAPIUrl, nil)
	req.Header.Set("Authorization", "Api-Token "+dtAPItoken)

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("The HTTP request failed with error")
		return nil
	} else {
		data, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(data, &locationsObject)
	}
	defer resp.Body.Close()
	var CheckLocations []string
	for _, location := range locationsObject.Locations {
		logger.Info("Synthetic Loctioan: " + location.EntityID)
		CheckLocations = append(CheckLocations, location.EntityID)
	}

	// Check if Synthetic is existing
	// /monitors?tag=keptn_service:${APPLICATION_SHORT_NAME}&tag=keptn_stage:${CI_ENVIRONMENT_SLUG}
	dtAPIUrl = "https://" + dtTenant + "/api/v1/synthetic/monitors?tag=keptn_service:" + data.Service + "&tag=keptn_stage:" + data.Stage + "&tag=keptn_project:" + data.Project
	req, err = http.NewRequest("GET", dtAPIUrl, nil)
	req.Header.Set("Authorization", "Api-Token "+dtAPItoken)

	resp, err = client.Do(req)
	if err != nil {
		logger.Error("The HTTP request failed with error")
		return nil
	} else {
		data, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal(data, &synteticsMonitorsObject)
	}
	defer resp.Body.Close()
	var monitorEntityID = ""
	var monitorMethod = "POST"
	for _, monitor := range synteticsMonitorsObject.Monitors {
		logger.Info("Found existing Monitor " + monitor.EntityID)
		monitorEntityID = monitor.EntityID
		monitorMethod = "PUT"
	}

	Monitor := httpMonitor{
		Name:         data.Project + `.` + data.Service + `.` + data.Stage,
		Enabled:      true,
		Type:         "HTTP",
		FrequencyMin: SyntheticFrequency,
		Locations:    CheckLocations,
		AnomalyDetection: AnomalyDetection{
			LoadingTimeThresholds: LoadingTimeThresholds{
				Enabled: false,
				Thresholds: []Thresholds{
					Thresholds{
						Type:    "TOTAL",
						ValueMs: 10000,
					},
				},
			},
			OutageHandling: OutageHandling{
				GlobalOutage: true,
				LocalOutage:  false,
				LocalOutagePolicy: LocalOutagePolicy{
					AffectedLocations: 1,
					ConsecutiveRuns:   3,
				},
			},
		},
		Script: Script{
			Version: "1.0",
			Requests: []ScriptRequests{
				ScriptRequests{
					Description: data.Project + `.` + data.Service + `.` + data.Stage,
					Method:      "GET",
					URL:         data.DeploymentURIPublic,
					RequestBody: "",
					Validation: Validation{
						RulesChaining: "or",
						Rules: []Rules{
							Rules{
								Value:       ">=400",
								PassIfFound: false,
								Type:        "httpStatusesList",
							},
						},
					},
					Configuration: Configuration{
						AcceptAnyCertificate: true,
						FollowRedirects:      true,
					},
				},
			},
		},
		Tags: []Tags{
			Tags{
				Source:  "USER",
				Context: "CONTEXTLESS",
				Key:     "keptn_project:" + data.Project,
			},
			Tags{
				Source:  "USER",
				Context: "CONTEXTLESS",
				Key:     "keptn_service:" + data.Service,
			},
			Tags{
				Source:  "USER",
				Context: "CONTEXTLESS",
				Key:     "keptn_stage:" + data.Stage,
			},
		},
		ManuallyAssignedApps: manuallyAssignedApps,
	}

	jsonStr, err := json.Marshal(Monitor)
	if err != nil {
		logger.Debug("Error formating json")
		return nil
	}
	logger.Debug(string(jsonStr))

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
		logger.Error("The HTTP request failed with error")
		return nil
	} else {
		data, _ := ioutil.ReadAll(resp.Body)
		logger.Debug(string(data))
	}
	defer resp.Body.Close()

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
