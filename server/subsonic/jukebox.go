package subsonic

import (
	"fmt"
	"net/http"

	"github.com/navidrome/navidrome/core/playback"
	"github.com/navidrome/navidrome/log"
	"github.com/navidrome/navidrome/server/subsonic/responses"
)

func (api *Router) JukeboxControl(r *http.Request) (*responses.Subsonic, error) {
	user, err := requiredParamString(r, "u")
	if err != nil {
		return nil, err
	}

	actionString, err := requiredParamString(r, "action")
	if err != nil {
		return nil, err
	}

	actionType := parseAction(actionString)
	if actionType == ActionUnknown {
		return nil, newError(responses.ErrorMissingParameter, "Unknown action: %s", actionString)
	}

	action, err := parseActionParameter(actionType, r)
	if err != nil {
		return nil, err
	}

	action.actionType = actionType
	action.user = user
	return handleJukeboxAction(action)
}

func handleJukeboxAction(action Action) (*responses.Subsonic, error) {
	log.Debug(fmt.Sprintf("Handle action: %s for user: %s, parameter: %v", action.actionType, action.user, action))
	playback := playback.GetInstance()

	switch action.actionType {
	case ActionGet:
		playback.Get(action.user)
	case ActionStatus:
		playback.Status(action.user)
	case ActionSet:
		playback.Set(action.user, action.Id)
	case ActionStart:
		playback.Start(action.user)
	case ActionStop:
		playback.Stop(action.user)
	case ActionSkip:
		playback.Skip(action.user, action.Index, action.Offset)
	case ActionAdd:
		playback.Add(action.user, action.Id)
	case ActionClear:
		playback.Clear(action.user)
	case ActionRemove:
		playback.Remove(action.user, action.Index)
	case ActionShuffle:
		playback.Shuffle(action.user)
	case ActionSetGain:
		playback.SetGain(action.user, action.Gain)
	}

	response := createJukeboxStatus(0, false, 0, 0)
	return response, nil
}

func createJukeboxStatus(currentIndex int64, playing bool, gain float64, position int64) *responses.Subsonic {
	response := newResponse()
	response.JukeboxStatus = &responses.JukeboxStatus{
		CurrentIndex: currentIndex,
		Playing:      playing,
		Gain:         gain,
		Position:     position,
	}
	return response
}