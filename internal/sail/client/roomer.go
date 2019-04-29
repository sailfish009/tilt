package client

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"github.com/windmilleng/tilt/internal/model"
)

// For injecting room creation logic (because the real way involves an HTTP request)
type SailRoomer interface {
	NewRoom(ctx context.Context) (info model.SailRoomInfo, err error)
}

type httpRoomer struct {
	addr model.SailURL
}

func (r httpRoomer) NewRoom(ctx context.Context) (info model.SailRoomInfo, err error) {
	u := r.addr.Http()
	u.Path = "/room"
	resp, err := http.Post(u.String(), "text/plain", nil)
	if err != nil {
		return model.SailRoomInfo{}, errors.Wrapf(err, "GET %s", u.String())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return model.SailRoomInfo{}, errors.Wrap(err, "reading /new_room response")
	}

	var roomInfo model.SailRoomInfo
	err = json.Unmarshal(body, &roomInfo)
	if err != nil {
		return model.SailRoomInfo{}, errors.Wrapf(err, "unmarshaling json: %s", string(body))
	}

	return roomInfo, nil
}

func ProvideSailRoomer(addr model.SailURL) SailRoomer {
	return httpRoomer{addr}
}