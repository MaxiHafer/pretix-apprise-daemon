package apprise

import (
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

type ClientTestSuite struct {
	suite.Suite
}

func (s *ClientTestSuite) TestClient_SendNotification() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Require().Equal("/notify/test", r.URL.Path)
		s.Require().Equal("POST", r.Method)
		s.Require().Equal("application/json", r.Header.Get("Content-Type"))

		var notification Notification
		s.Require().NoError(json.NewDecoder(r.Body).Decode(&notification))

		s.Require().Equal("test", notification.Body)

		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := NewClient(&ClientConfig{
		Host: srv.URL,
		Key:  "test",
	})

	err := client.SendNotification(&Notification{
		Body: "test",
	})
	s.Require().NoError(err)
}
