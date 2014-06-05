package LeagueApi

import (
	"net/http"
	"testing"
)

func GetMock(string url) (r *http.Response, e error) {
	return nil, nil
}
