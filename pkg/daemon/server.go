package daemon

import (
	"bytes"
	"encoding/json"
	"github.com/maxihafer/pretix-apprise-daemon/pkg/apprise"
	"github.com/maxihafer/pretix-apprise-daemon/pkg/pretix"
	"github.com/maxihafer/pretix-apprise-daemon/templates"
	"log/slog"
	"net"
	"net/http"
	"strconv"
)

type Config struct {
	BindAddress string `mapstructure:"daemon-bind-address"`
	Port        int    `mapstructure:"daemon-port"`
	AppriseHost string `mapstructure:"apprise-host"`
	AppriseKey  string `mapstructure:"apprise-key"`
	PretixToken string `mapstructure:"pretix-token"`
	PretixHost  string `mapstructure:"pretix-host"`
}

func NewServer(config *Config) *Server {
	pretixClient := pretix.NewClient(&pretix.ClientConfig{
		Token: config.PretixToken,
		Host:  config.PretixHost,
	})

	appriseClient := apprise.NewClient(&apprise.ClientConfig{
		Host: config.AppriseHost,
		Key:  config.AppriseKey,
	})

	return &Server{
		config:  config,
		pretix:  pretixClient,
		apprise: appriseClient,
	}
}

type Server struct {
	config *Config

	pretix  *pretix.Client
	apprise *apprise.Client
}

func (s *Server) Run() error {
	srv := http.NewServeMux()

	srv.HandleFunc("/webhook", s.handleWebhook)

	return http.ListenAndServe(net.JoinHostPort(s.config.BindAddress, strconv.Itoa(s.config.Port)), srv)
}

func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	webhook := &pretix.Webhook{}

	if err := json.NewDecoder(r.Body).Decode(webhook); err != nil {

		slog.ErrorContext(r.Context(), "error decoding webhook payload: %v", err)
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	slog.InfoContext(r.Context(), "received webhook", "webhook", webhook)

	order, err := s.pretix.GetOrder(webhook.Organizer, webhook.Event, webhook.Code)
	if err != nil {
		slog.ErrorContext(r.Context(), "error fetching order", "error", err)
		http.Error(w, "error fetching order", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(r.Context(), "fetched order", "order", order)

	notification, err := s.buildNotification(webhook.Action, order)

	if err != nil {
		slog.ErrorContext(r.Context(), "error building notification", "error", err)
		http.Error(w, "error building notification", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(r.Context(), "built notification", "notification", notification)

	if err := s.apprise.SendNotification(notification); err != nil {
		slog.ErrorContext(r.Context(), "error sending notification", "error", err)
		http.Error(w, "error sending notification", http.StatusInternalServerError)
		return
	}

	slog.InfoContext(r.Context(), "sent notification")

	w.WriteHeader(http.StatusOK)
}

func (s *Server) buildNotification(action pretix.Action, order *pretix.Order) (*apprise.Notification, error) {
	var notification *apprise.Notification
	body := bytes.Buffer{}

	switch action {
	case pretix.ActionOrderPlaced:
		if err := templates.ExecuteOrderPlacedDE(&body, order); err != nil {
			return nil, err
		}

		notification = &apprise.Notification{
			Title: "Neue Bestellung eingegangen",
			Body:  body.String(),
			Type:  apprise.Success,
		}
	}
	/*case pretix.ActionOrderPaid:
		if err := templates.ExecuteOrderPaidDE(&body, order); err != nil {
			return nil, err
		}

		notification = apprise.Notification{
			Title: "Bestellung bezahlt",
			Body:  body.String(),
			Type:  apprise.Success,
		}
	case pretix.ActionOrderCanceled:
		if err := templates.ExecuteOrderCanceledDE(&body, order); err != nil {
			return nil, err
		}

		notification = apprise.Notification{
			Title: "Bestellung storniert",
			Body:  body.String(),
			Type:  apprise.Failure,
		}
	case pretix.ActionOrderExpired:
		if err := templates.ExecuteOrderExpiredDE(&body, order); err != nil {
			return nil, err
		}

		notification = apprise.Notification{
			Title: "Bestellung abgelaufen",
			Body:  body.String(),
			Type:  apprise.Warning,
		}*/

	return notification, nil
}
