package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2" // make sure to use v2 cloudevents here

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

/**
* Here are all the handler functions for the individual event
* See https://github.com/keptn/spec/blob/0.8.0-alpha/cloudevents.md for details on the payload
**/

// GenericLogKeptnCloudEventHandler is a generic handler for Keptn Cloud Events that logs the CloudEvent
func GenericLogKeptnCloudEventHandler(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data interface{}) error {
	log.Printf("Handling %s Event: %s", incomingEvent.Type(), incomingEvent.Context.GetID())
	log.Printf("CloudEvent %T: %v", data, data)

	return nil
}

// HandleTestsTriggered hanldes test.triggered (used to be sh.keptn.events.deployment-finished)
func HandleTestsTriggered(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.TestTriggeredEventData) error {
	log.Printf("Handling Tests Triggered Event: %s", incomingEvent.Context.GetID())

	myKeptn.SendTaskStartedEvent(data, ServiceName)

	// run tests
	log.Printf("Looking for Litmus chaos experiment in Keptn git repo...")

	keptnResourceContent, err := myKeptn.GetKeptnResource(LitmusExperimentFileName)

	if err != nil {
		logMessage := fmt.Sprintf("No %s file found for service %s in stage %s in project %s", LitmusExperimentFileName, data.Service, data.Stage, data.Project)
		log.Printf(logMessage)
		_, err = myKeptn.SendTaskFinishedEvent(&keptnv2.EventData{
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
			Message: logMessage,
		}, ServiceName)

		return err
	}

	_ = os.Mkdir("litmus", 0644)
	err = ioutil.WriteFile(LitmusExperimentFileName, []byte(keptnResourceContent), 0644)
	if err != nil {
		log.Printf("Could not store experiment file locally: %s", err.Error())
	}

	// construct the namespace of the chaos resources
	projectAndNamespace := data.Project + "-" + data.Stage

	// Obtain the ChaosEngine name from the Keptn experiment manifest
	chaosEngineName, err := ExecuteCommand("kubectl", []string{"apply", "-f", LitmusExperimentFileName, "--dry-run", "-o", "jsonpath='{.metadata.name}'"})
	if err != nil {
		log.Printf("Error while extracting chaosengine name from manifest: %s", err.Error())
	}
	chaosEngineName = strings.Trim(chaosEngineName, `'"`)
	log.Printf("Name of ChaosEngine: %s", chaosEngineName)

	log.Printf("Executing Litmus chaos experiment...")
	output, err := ExecuteCommand("kubectl", []string{"apply", "-f", LitmusExperimentFileName})
	if err != nil {
		log.Printf("Error execute kubectl apply command: %s", err.Error())
	}
	log.Printf("ChaosEngine create command finished with: %s", output)

	// Allow the chaos-operator to patch the engine with the initial status
	time.Sleep(2 * time.Second)

	// Extract the chaosUID for use in result extraction
	uid, err := ExecuteCommand("kubectl", []string{"get", "chaosengine", chaosEngineName, "-o", "jsonpath='{.metadata.uid}'", "-n", projectAndNamespace})
	if err != nil {
		log.Printf("Error while retrieving chaosengine UID: %s", err.Error())
	}
	chaosUID := strings.Trim(string(uid), `'"`)
	log.Printf("UID of ChaosEngine %s: %s", chaosEngineName, chaosUID)

	var chaosStatus string
	for chaosStatus != "completed" {
		log.Printf("Waiting for completion of chaos experiment..")
		chaosStatus, err = ExecuteCommand("kubectl", []string{"get", "chaosengine", chaosEngineName, "-o", "jsonpath='{.status.engineStatus}'", "-n", projectAndNamespace})
		if err != nil {
			log.Printf("Error while retrieving chaos status: %s", err.Error())
			break
		}
		chaosStatus = strings.Trim(chaosStatus, `'"`)
		// interval before we check the chaosengine status again
		time.Sleep(2 * time.Second)
	}

	log.Printf("Chaos experiment is completed")

	// Construct the jsonpath filter to extract verdict of the chaosresult
	jsonPathFilterForResult := fmt.Sprintf("jsonpath='{.items[?(@.metadata.labels.chaosUID==\"%s\")].status.experimentstatus.verdict}'", chaosUID)

	// Getting ChaosResult Data
	verdict, err := ExecuteCommand("kubectl", []string{"get", "chaosresult", "-o", jsonPathFilterForResult, "-n", projectAndNamespace})
	if err != nil {
		log.Printf("Error while retrieving chaos result: %s", err.Error())
	}
	verdict = strings.Trim(verdict, `'"`)
	log.Println("ChaosExperiment Verdict: " + verdict)

	testResult := keptnv2.ResultFailed

	// check verdict
	log.Println("Final Result: " + verdict)

	if verdict == "Pass" {
		testResult = keptnv2.ResultPass
	}

	// send tests.finished
	_, err = myKeptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Status:  keptnv2.StatusSucceeded,
		Result:  testResult,
		Message: "Chaos tests finished",
	}, ServiceName)

	return err
}

// HandleTestFinished handles test.finished event
func HandleTestFinished(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *keptnv2.TestFinishedEventData) error {
	log.Printf("Handling Tests Finished Event: %s", incomingEvent.Context.GetID())

	if incomingEvent.Source() == ServiceName {
		// skip test.finished, it has been send out by litmus service
		return nil
	}

	// delete chaos experiment
	log.Printf("Deleting chaos experiment resources")
	_, err := ExecuteCommand("kubectl", []string{"delete", "-f", LitmusExperimentFileName})
	if err != nil {
		log.Printf("Error execute kubectl delete command: %s", err.Error())
	}

	return nil
}

// ExecuteCommand exectues the command using the args
func ExecuteCommand(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("Error executing command %s %s: %s\n%s", command, strings.Join(args, " "), err.Error(), string(out))
	}
	return string(out), nil
}
