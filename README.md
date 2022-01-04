## What

This is an experiment to send audit events from GCP to slack.

## How

1. Create a pubsub topic
2. Create a log sink in Cloud Logging and filter cloudaudit log to be routed to pubsub
3. Create a trigger in pubsub to invoke Cloud Function
4. Cloud Function parse incoming message and send to slack
