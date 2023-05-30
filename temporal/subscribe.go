package emaildrips

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

// TaskQueueName is the Temporal Task Queue for the worker and Client.
var TaskQueueName = "email_drips"

var baseURL = "http://localhost:3000/campaigns/"

// Subscription is the user email and the campaign they'll receive.
type Subscription struct {
	EmailAddress string
	CampaignID   string
	ID           string
}

// EmailInfo is the data that the SendContentEmail uses to send the message.
type EmailInfo struct {
	ToAddress string
	EmailPath string
}

// Email represents the email from the campaign
type Email struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Body string `json:"body"`
}

// UserSubscriptionWorkflow handles subscribing users. Accepts a Subsription.
func UserSubscriptionWorkflow(ctx workflow.Context, subscription Subscription) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Subscription created for " + subscription.EmailAddress)

	// How frequently to send the messages
	duration := time.Minute

	// errors
	var err error

	// duration := (24 * 7) * time.Hour

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
		WaitForCancellation: true,
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	// Handle any cleanup, including cancellations..
	defer func() {

		if !errors.Is(ctx.Err(), workflow.ErrCanceled) {
			return
		}

		//	Cancellation received, which will trigger an unsubscribe email.
		logger.Info("Sending unsubscribe email to " + subscription.EmailAddress)

		newCtx, cancel := workflow.NewDisconnectedContext(ctx)

		defer cancel()

		data := EmailInfo{
			ToAddress: subscription.EmailAddress,
			EmailPath: baseURL + subscription.CampaignID + "/unsubscribe_message.json",
		}

		err = workflow.ExecuteActivity(newCtx, SendContentEmail, data).Get(newCtx, nil)

		if err != nil {
			logger.Error("Unable to send unsubscribe message", "Error", err)
		}

	}()

	// welcome message
	logger.Info("Sending welcome email to " + subscription.EmailAddress)

	data := EmailInfo{
		ToAddress: subscription.EmailAddress,
		EmailPath: baseURL + subscription.CampaignID + "/welcome_message.json",
	}

	err = workflow.ExecuteActivity(ctx, SendContentEmail, data).Get(ctx, nil)

	if err != nil {
		logger.Error("Failed to send welcome email", "Error", err)
	}

	// Get campaign emails
	logger.Info("Getting emails for campaign" + subscription.CampaignID)

	var emails []string

	err = workflow.ExecuteActivity(ctx, GetCampaignEmails, subscription.CampaignID).Get(ctx, &emails)

	if err != nil {
		logger.Error("Failed to get campaign emails", "Error", err)
	}

	for _, mail := range emails {

		data := EmailInfo{
			ToAddress: subscription.EmailAddress,
			EmailPath: baseURL + subscription.CampaignID + "/emails/" + mail + ".json",
		}

		err = workflow.ExecuteActivity(ctx, SendContentEmail, data).Get(ctx, nil)

		if err != nil {
			logger.Error("Failed to send email "+mail, "Error", err)
		}

		logger.Info("sent content email " + mail + " to " + subscription.EmailAddress)

		workflow.Sleep(ctx, duration)
	}

	return nil
}

// GetCampaignEmails is the activity that fetches the email IDs from the API.
func GetCampaignEmails(ctx context.Context, campaignID string) ([]string, error) {

	logger := activity.GetLogger(ctx)
	logger.Info("getting email IDs from server")

	url := baseURL + campaignID + "/emails.json"

	logger.Info("fetching email IDs from" + url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var emails []Email
	err = json.Unmarshal(body, &emails)
	if err != nil {
		return nil, err
	}

	ids := make([]string, len(emails))
	for i, email := range emails {
		ids[i] = strconv.Itoa(email.ID)
	}

	return ids, nil
}

// SendContentEmail is the activity that sends the email to the customer.
func SendContentEmail(ctx context.Context, emailInfo EmailInfo) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Sending email " + emailInfo.EmailPath + " to " + emailInfo.ToAddress)

	message, err := getEmailBody(emailInfo.EmailPath)

	// call mailer api here.
	if err != nil {
		logger.Error("Failed getting email", err)
		return errors.New("unable to locate message to send")
	}

	logger.Info("message is " + message)
	return sendMail(message, emailInfo.ToAddress)
}

// functions

// getEmailFromFile gets the email from the specified text file.
func getEmailBody(url string) (string, error) {

	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var email Email
	err = json.Unmarshal(body, &email)
	if err != nil {
		return "", err
	}

	return email.Body, nil
}

// mocked mail.
func sendMail(message string, email string) error {
	return nil
}
