package bot

import (
	"context"
	"errors"
	"fmt"
	"github.com/her0elt/florida-man-bot/dto"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"
	"math/rand"
	"os"
)

func Run() {
	OAUTH_TOKEN := os.Getenv("OAUTH_TOKEN")
	CHANNEL_ID := os.Getenv("CHANNEL_ID")
	APP_TOKEN := os.Getenv("APP_TOKEN")

	client := slack.New(OAUTH_TOKEN, slack.OptionDebug(true), slack.OptionAppLevelToken(APP_TOKEN))
	socket := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context, client *slack.Client, socket *socketmode.Client, channel_id string) {
		for {
			select {
			case <-ctx.Done():
				log.Println("Shutting down socketmode listener")
				return
			case event := <-socket.Events:
				switch event.Type {

				case socketmode.EventTypeEventsAPI:
					eventsAPI, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
						continue
					}
					socket.Ack(*event.Request)
					err := HandleEventMessage(eventsAPI, client, channel_id)
					if err != nil {
						log.Fatal(err)
						continue
					}
				}
			}
		}
	}(ctx, client, socket, CHANNEL_ID)

	socket.Run()
}

func HandleEventMessage(event slackevents.EventsAPIEvent, client *slack.Client, slack_channel string) error {
	switch event.Type {
	case slackevents.CallbackEvent:

		innerEvent := event.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			err := HandleAppMentionEventToBot(ev, client, slack_channel)
			if err != nil {
				return err
			}
		}
	default:
		return errors.New("unsupported event type")
	}
	return nil
}

func HandleAppMentionEventToBot(event *slackevents.AppMentionEvent, client *slack.Client, slack_channel string) error {

	user, err := client.GetUserInfo(event.User)
	if err != nil {
		return err
	}

	GOOGLE_API_KEY := os.Getenv("GOOGLE_API_KEY")
	GOOGLE_SEARCH_CONTEXT := os.Getenv("GOOGLE_SEARCH_CONTEXT")
	gResp := dto.MakeRequest(GOOGLE_API_KEY, GOOGLE_SEARCH_CONTEXT)
	randomIndex := rand.Intn(len(gResp.Items))
	attachment := slack.Attachment{
		Pretext: gResp.Items[randomIndex].Title,
		Text:    fmt.Sprintf("<%s|Florida man news today>", gResp.Items[randomIndex].Link),
	}

	_, _, err = client.PostMessage(
		slack_channel,
		slack.MsgOptionText(fmt.Sprintf("Your florida man news for today <@%s>", user.Name), false),
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true),
	)
	if err != nil {
		return fmt.Errorf("failed to post message: %w", err)
	}
	return nil
}
