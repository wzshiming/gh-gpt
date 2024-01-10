package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/wzshiming/gh-gpt/pkg/api"
	"github.com/wzshiming/gh-gpt/pkg/auth"
)

func setRespHeaders(h http.Header) {
	h.Set("Transfer-Encoding", "chunked")
	h.Set("X-Accel-Buffering", "no")
	h.Set("Content-Type", "text/event-stream; charset=utf-8")
	h.Set("Cache-Control", "no-cache")
	h.Set("Connection", "keep-alive")
}

type responseError struct {
	Error string `json:"error"`
}

func respJSONError(w http.ResponseWriter, code int, err string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(responseError{
		Error: err,
	})
}

type Option func(*Server)

func WithAuth(auth auth.Auth) func(*Server) {
	return func(s *Server) {
		s.auth = auth
	}
}

func WithClient(client *api.Client) func(*Server) {
	return func(s *Server) {
		s.client = client
	}
}

type Server struct {
	auth   auth.Auth
	client *api.Client
}

func NewServer(opts ...Option) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) ChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	oauthToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if oauthToken == "" {
		var err error
		oauthToken, err = s.auth.GetToken()
		if err != nil {
			respJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		respJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	req := api.ChatRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		respJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ctx := r.Context()
	token, err := s.client.TokenWishCache(ctx, oauthToken)
	if err != nil {
		respJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	setRespHeaders(w.Header())

	encodeer := json.NewEncoder(w)

	fn := func(chatResponse api.ChatResponse) error {
		if req.Stream {
			_, err := w.Write([]byte("data: "))
			if err != nil {
				return err
			}
		}

		err := encodeer.Encode(chatResponse)
		if err != nil {
			return err
		}

		_, err = w.Write([]byte("\n"))
		if err != nil {
			return err
		}
		return nil
	}

	err = s.client.ChatCompletions(ctx, token, &req, fn)
	if err != nil {
		respJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if req.Stream {
		_, err := w.Write([]byte("data: [DONE]\n"))
		if err != nil {
			respJSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
}
