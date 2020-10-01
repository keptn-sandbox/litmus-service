package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"
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
	-> "sh.keptn.event.action.triggered"
*/

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

	// capture start time for tests
	startTime := time.Now()

	// run tests
	// ToDo: Implement your tests here

	log.Printf("looking for Litmus chaos experiment in Keptn git repo...")

	resourceHandler := keptnapi.NewResourceHandler("configuration-service:8080")

	keptnResourceContent, err := resourceHandler.GetServiceResource(data.Project, data.Stage, data.Service, LitmusExperimentFileName)
	var fileContent []byte
	if err != nil {
		logMessage := fmt.Sprintf("No %s file found for service %s in stage %s in project %s", LitmusExperimentFileName, data.Service, data.Stage, data.Project)
		log.Printf(logMessage)
		return errors.New(logMessage)
	}
	fileContent = []byte(keptnResourceContent.ResourceContent)

	_ = os.Mkdir("litmus", 0644)
	err = ioutil.WriteFile(LitmusExperimentFileName, fileContent, 0644)
	if err != nil {
		log.Printf("could not store experiment file locally: %s", err.Error())
	}

	log.Printf("executing Litmus chaos experiment...")
	output, err := ExecuteCommand("kubectl", []string{"apply", "-f", LitmusExperimentFileName})
	if err != nil {
		log.Printf("Error execute kubectl apply command: %s", err.Error())
	}
	log.Printf("Execute command finished with: %s", output)

	// Allow the chaos-operator to patch the engine with the initial status 
	time.Sleep(2 * time.Second)

	var chaosStatus string
	for chaosStatus != "completed" {
		log.Printf("Waiting for completion of chaos experiment..")
		chaosStatus, err = ExecuteCommand("kubectl", []string{"get", "chaosengine", "carts-chaos", "-o", "custom-columns=:status.engineStatus", "-n", "litmus-chaos", "--no-headers"})
		if err != nil {
                	log.Printf("Error while retrieving chaos status: %s", err.Error())
			break 
        	}
		// interval before we check the chaosengine status again
		time.Sleep(2 * time.Second)
	}

	log.Printf("Chaos experiment is completed")
	// Send Test Finished Event
	return myKeptn.SendTestsFinishedEvent(&incomingEvent, "", "", startTime, "pass", nil, "litmus-service")
	//return nil
}

//
// Handles TestsFinishedEventType = "sh.keptn.events.tests-finished"
// TODO: add in your handler code
//
func HandleTestsFinishedEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.TestsFinishedEventData) error {
	log.Printf("Handling Tests Finished Event: %s", incomingEvent.Context.GetID())

	// delete chaos experiment
	log.Printf("Deleting chaos experiment resources")
	_, err := ExecuteCommand("kubectl", []string{"delete", "-f", LitmusExperimentFileName})
	if err != nil {
                log.Printf("Error execute kubectl delete command: %s", err.Error())
        }

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

	// Deprecated since Keptn 0.7.0 - use the HandleActionTriggeredEvent instead

	return nil
}

//
// Handles ActionTriggeredEventType = "sh.keptn.event.action.triggered"
// TODO: add in your handler code
//
func HandleActionTriggeredEvent(myKeptn *keptn.Keptn, incomingEvent cloudevents.Event, data *keptn.ActionTriggeredEventData) error {
	log.Printf("Handling Action Triggered Event: %s", incomingEvent.Context.GetID())

	// check if action is supported
	if data.Action.Action == "action-xyz" {
		//myKeptn.SendActionStartedEvent() TODO: implement the SendActionStartedEvent in keptn/go-utils/pkg/lib/events.go

		// Implement your remediation action here

		//myKeptn.SendActionFinishedEvent() TODO: implement the SendActionFinishedEvent in keptn/go-utils/pkg/lib/events.go
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
