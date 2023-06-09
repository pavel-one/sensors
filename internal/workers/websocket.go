package workers

import (
	"github.com/goccy/go-json"
	"github.com/olahol/melody"
	"github.com/pavel-one/SimpleServerMonitor/internal/events"
)

type message struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

// WebsocketWorker worker from work websocket's
func WebsocketWorker(server *melody.Melody, ch events.Chan) error {
	for event := range ch {
		msg := message{
			Event: event.GetEvent(),
			Data:  event.GetData(),
		}

		b, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		err = server.BroadcastFilter(b, func(session *melody.Session) bool {
			return !session.IsClosed()
		})
		if err != nil {
			return err
		}

	}

	return nil
}
