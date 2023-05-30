package main

import (
	"context"
	"emaildrips"
	"log"

	"go.temporal.io/sdk/client"
)

func main() {

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "email_drip_1",
		TaskQueue: emaildrips.TaskQueueName,
	}

	input := emaildrips.Subscription{
		EmailAddress: "bphogan@example.com",
		CampaignID:   "1",
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, emaildrips.UserSubscriptionWorkflow, input)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	var result string

	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}

	log.Printf("Workflow result: %s\n", result)

}
