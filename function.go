// Package p contains a Pub/Sub Cloud Function.
package p

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/slack-go/slack"
)

// PubSubMessage is the payload of a Pub/Sub event. Please refer to the docs for
// additional information regarding Pub/Sub events.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

type BindingDelta struct {
	Action string `json:"action"`
	Member string `json:"member"`
	Role   string `json:"role"`
}

type PolicyDelta struct {
	BindingDeltas []BindingDelta `json:"bindingDeltas"`
}

type ServiceData struct {
	PolicyDelta PolicyDelta `json:"policyDelta"`
}

type AuthenticationInfo struct {
	PrincipalEmail string `json:"principalEmail"`
}

type ProtoPayload struct {
	AuthenticationInfo AuthenticationInfo `json:"authenticationInfo"`
	ServiceData        ServiceData        `json:"serviceData"`
}

type Labels struct {
	ProjectID string `json:"project_id"`
}

type Resource struct {
	Labels Labels `json:"labels"`
}

type DecodedMessage struct {
	ProtoPayload ProtoPayload `json:"protoPayload"`
	Resource     Resource     `json:"resource"`
}

// ConsumePubSub consumes a Pub/Sub message.
func ConsumePubSub(ctx context.Context, m PubSubMessage) error {
	token := os.Getenv("SLACK_TOKEN")
	channel := os.Getenv("SLACK_CHANNEL")

	if token == "" && channel == "" {
		log.Println("Required SLACK_TOKEN & SLACK_CHANNEL environment variable")
		return errors.New("SLACK_TOKEN/SLACK_CHANNEL not found")
	}

	var msg DecodedMessage
	if err := json.Unmarshal(m.Data, &msg); err != nil {
		log.Println(err)
		return err
	}

	// user:xxx@17.media
	member := msg.ProtoPayload.ServiceData.PolicyDelta.BindingDeltas[0].Member
	email := strings.Split(member, ":")[1]
	admin := msg.ProtoPayload.AuthenticationInfo.PrincipalEmail
	project := msg.Resource.Labels.ProjectID
	action := msg.ProtoPayload.ServiceData.PolicyDelta.BindingDeltas[0].Action

	title := slack.NewTextBlockObject("mrkdown", "GCP", false, false)
	titleSection := slack.NewSectionBlock(title, nil, nil)

	body := slack.NewTextBlockObject("mrkdown", fmt.Sprintf("[GCP][%s] %s `%s` by `%s`", project, action, email, admin), false, false)
	bodySection := slack.NewSectionBlock(body, nil, nil)

	msgOptions := slack.MsgOptionBlocks(titleSection, bodySection)

	api := slack.New(token)
	api.PostMessage(channel, msgOptions)

	return nil
}
