package pretix

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

func (s *ClientTestSuite) TestClient_GetOrder() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Require().Equal("/api/v1/organizers/test/events/test/orders/test/", r.URL.Path)
		s.Require().Equal("GET", r.Method)
		s.Require().Equal("Token test", r.Header.Get("Authorization"))

		order := &Order{
			Code:            "test",
			Event:           "test",
			Status:          OrderStatusPending,
			TestMode:        false,
			PaymentProvider: "bank_transfer",
			Total:           "1337€",
			OrderURL:        "https://1337.1337.1337",
		}

		body, err := json.Marshal(order)
		s.Require().NoError(err)

		_, err = w.Write(body)
		s.Require().NoError(err)
	}))
	defer srv.Close()

	client := NewClient(&ClientConfig{
		Token: "test",
		Host:  srv.URL,
	})

	order, err := client.GetOrder("test", "test", "test")
	s.Require().NoError(err)

	s.Require().Equal("test", order.Code)
	s.Require().Equal("test", order.Event)
	s.Require().Equal(OrderStatusPending, order.Status)
	s.Require().Equal(false, order.TestMode)
	s.Require().Equal("bank_transfer", order.PaymentProvider)
	s.Require().Equal("1337€", order.Total)
	s.Require().Equal("https://1337.1337.1337", order.OrderURL)
}
