# Temporal Rails and Go example

This is an experimental example app that uses Ruby on Rails to talk to a Temporal Workflow running in Go.

This application lets you manage email drip campaigns and let users subscribe to them. The campaign and email management happens through the Rails application and its database. Subscription creation and cancellation start on the Rails side, but the subscription itself is a Temporal Workflow.

The Workflow calls back to the Rails application's REST API to fetch the emails in the campaign to send.


Prerequisites:

* Ruby 3.1.4 (Ruby 3.2 is not supported yet.)
* Go 19+
* The Temporal CLI application

## Setup

Clone the repository.

Install the Temporal CLI :

```bash
curl -sSf https://temporal.download/cli.sh | sh
```

Start the Temporal server:

```bash
temporal server start-dev
```

Install the Rails app's dependencies:

```bash
bundle install
```

Start the Rails app:

```
rails s
```

Visit http://localhost:3000/campaigns and add a campaign. Then add a couple emails.

Now run the Temporal Worker:

```bash
cd temporal
go mod tidy
go run worker/main.go
```

Now visit http://localhost:3000/subscriptions/new and subscribe to the campaign.

To cancel the subscription, delete the Subscription record which cancels the Workflow.


You can also use the Temporal CLI to start a subscription:

```
temporal workflow start --task-queue=email_drips --type=UserSubscriptionWorkflow --input='{"EmailAddress": "brian@example.com", "CampaignID": "1"}'
```

## TODO

* Properly handle cancellation instead of swalling errors
* Switch subscription ID to use a UUID rather than the Active Record ID.
* Move subscription list to restrict by email and move overall subscription management to an admin controller.
* lock campaigns behind auth tokens.

