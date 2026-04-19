package slack

import (
	"fmt"

	slackgo "github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"k8s.io/client-go/kubernetes"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned" // NEW IMPORT
)

// Start initializes Socket Mode and begins listening for events
// NOTICE: Added metricsClient to the arguments
func Start(appToken, botToken string, k8sClient *kubernetes.Clientset, metricsClient *metricsv.Clientset) {
	api := slackgo.New(
		botToken,
		slackgo.OptionDebug(false),
		slackgo.OptionAppLevelToken(appToken),
	)
	client := socketmode.New(api)

	go func() {
		for evt := range client.Events {
			switch evt.Type {

			case socketmode.EventTypeEventsAPI:
				eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
				if !ok {
					continue
				}

				client.Ack(*evt.Request)

				if eventsAPIEvent.Type == slackevents.CallbackEvent {
					innerEvent := eventsAPIEvent.InnerEvent
					switch ev := innerEvent.Data.(type) {
					case *slackevents.AppMentionEvent:
						// Pass BOTH clients to the router
						handleMention(api, k8sClient, metricsClient, ev)
					}
				}

			case socketmode.EventTypeInteractive:
				callback, ok := evt.Data.(slackgo.InteractionCallback)
				if !ok {
					continue
				}

				client.Ack(*evt.Request)

				// Pass BOTH clients to the router
				handleInteraction(api, k8sClient, metricsClient, callback)
			}
		}
	}()

	fmt.Println("🤖 Kubebot is actively listening to Slack via Socket Mode...")
	client.Run()
}