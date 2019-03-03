package votebot

import (
	"encoding/json"
	"net/http"

	"github.com/emoji/internal/mongodb"
	"github.com/go-chi/chi"
)

// Server configures router
type Server struct {
	client *mongodb.Client
}

// UseClient applies the mongodb client
func UseClient(c *mongodb.Client) func(*Server) error {
	return func(s *Server) error {
		s.client = c
		return nil
	}
}

// NewServer creats a new Server
func NewServer(opts ...func(*Server) error) (*Server, error) {
	s := &Server{}

	for _, f := range opts {
		if err := f(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Route retruns the routing rules
func (s *Server) Route() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/healthcheck", s.healthCheckHandler)
	r.Get("/vote/{emojiID}", s.voteHandler)
	r.Get("/emoji/{emojiID}", s.emojiHandler)
	return r
}

// healthCheckHandler handles healthcheck request
func (s *Server) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (s *Server) emojiHandler(w http.ResponseWriter, r *http.Request) {
	emojiID := chi.URLParam(r, "emojiID")

	emoji, err := s.client.GetEmoji(emojiID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(emoji)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// voteHandler handles vote request
func (s *Server) voteHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var emoji *mongodb.Emoji

	emojiID := chi.URLParam(r, "emojiID")

	if emoji, err = s.client.GetEmoji(emojiID); err != nil {
		emoji = &mongodb.Emoji{
			EmojiID: emojiID,
			Vote:    0,
		}
		s.client.AddEmoji(emoji)
	} else {
		s.client.UpdateEmoji(emoji)
	}

	js, err := json.Marshal(emoji)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
