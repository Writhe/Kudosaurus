package server

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/slack-go/slack/slackevents"
)

func (s *CommandServer) handleHomeEvent(ev *slackevents.AppHomeOpenedEvent, teamID string) {
	s.publishView(
		s.getHomeView(ev.User, teamID),
	)
}

func (s *CommandServer) handleEvent(res http.ResponseWriter, req *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	body := buf.String()
	eventsAPIEvent, e := slackevents.ParseEvent(
		json.RawMessage(body),
		slackevents.OptionVerifyToken(
			&slackevents.TokenComparator{VerificationToken: s.config.VerificationToken},
		),
	)

	if e != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
		}
		res.Header().Set("Content-Type", "text")
		res.Write([]byte(r.Challenge))
		return
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent

		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppHomeOpenedEvent:
			s.handleHomeEvent(ev, eventsAPIEvent.TeamID)

		default:
			s.logger.Printf("Unhandled event - %s", innerEvent.Type)
		}
	}
}
