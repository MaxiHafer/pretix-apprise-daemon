package daemon

import (
	"bytes"
	"encoding/json"
	"github.com/maxihafer/pretix-apprise-daemon/pkg/pretix"
	"github.com/maxihafer/pretix-apprise-daemon/templates"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDaemonTestSuite(t *testing.T) {
	suite.Run(t, new(DaemonTestSuite))
}

type DaemonTestSuite struct {
	suite.Suite
}

var (
	order = &pretix.Order{
		Code:            "test",
		Event:           "test",
		Status:          pretix.OrderStatusPending,
		TestMode:        false,
		PaymentProvider: "bank_transfer",
		Total:           "1337€",
		OrderURL:        "https://1337.1337.1337",
	}

	webhook = &pretix.Webhook{
		Organizer: "test",
		Event:     "test",
		Code:      "test",
		Action:    pretix.ActionOrderPlaced,
	}
)

func (s *DaemonTestSuite) TestDaemon_buildNotification() {

	daemon := NewServer(&Config{
		BindAddress: "localhost",
		Port:        8080,
		AppriseHost: "test",
		AppriseKey:  "test",
		PretixToken: "test",
		PretixHost:  "test",
	})

	expectedBody := &bytes.Buffer{}
	s.Require().NoError(templates.ExecuteOrderPlacedDE(expectedBody, nil))

	notification, err := daemon.buildNotification(pretix.ActionOrderPlaced, order)
	s.Require().NoError(err)

	s.Require().NotNil(notification)
	s.Require().Equal(expectedBody.String(), notification.Body)
}

func (s *DaemonTestSuite) TestDaemon_handleWebhook() {
	appriseSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Require().Equal("/notify/test", r.URL.Path)
		s.Require().Equal("POST", r.Method)
		s.Require().Equal("application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
	}))
	defer appriseSrv.Close()

	pretixSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Require().Equal("/api/v1/organizers/test/events/test/orders/test/", r.URL.Path)
		s.Require().Equal("GET", r.Method)
		s.Require().Equal("Token test", r.Header.Get("Authorization"))

		s.Require().NoError(json.NewEncoder(w).Encode(order))
	}))
	defer pretixSrv.Close()

	daemon := NewServer(&Config{
		BindAddress: "localhost",
		Port:        8080,
		AppriseHost: appriseSrv.URL,
		AppriseKey:  "test",
		PretixToken: "test",
		PretixHost:  pretixSrv.URL,
	})

	webhookBody, err := json.Marshal(webhook)
	s.Require().NoError(err)

	req := httptest.NewRequest("POST", "/webhook", bytes.NewBuffer(webhookBody))

	w := httptest.NewRecorder()
	daemon.handleWebhook(w, req)

	s.Require().Equal(http.StatusOK, w.Code)
}
