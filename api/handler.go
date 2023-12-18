package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

var (
	ErrNotFound      = errors.New("NotFound")
	ErrInternalError = errors.New("InternalError")
)

// Post posts request built from server request struct.
func Post[ReqT any, ResT any](client *http.Client, url string) func(context.Context, ReqT) (ResT, error) {
	return func(ctx context.Context, r ReqT) (ResT, error) {
		var res ResT
		b, err := json.Marshal(r)
		if err != nil {
			return res, err
		}
		buf := bytes.NewBuffer(b)
		resp, err := client.Post(url, "application/json", buf)
		if err != nil {
			return res, err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return res, err
		}

		switch resp.StatusCode {
		case http.StatusOK:
			if err := json.Unmarshal(body, &res); err != nil {
				return res, err
			}
			return res, nil
		case http.StatusNotFound:
			return res, ErrNotFound
		default:
			return res, ErrInternalError
		}
	}
}

// API builds a http handler from server request handler.
func API[ReqT any, ResT any](f func(context.Context, ReqT) (ResT, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		logger := slog.With(slog.String("url", r.URL.String()))

		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("read body", slog.Any("error", err))
			return
		}

		logger.Info("read body", slog.Any("body", b))

		var req ReqT
		if err := json.Unmarshal(b, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "Invalid request")
			logger.Info("unmarshal failed", slog.Any("error", err))
			return
		}

		res, err := f(r.Context(), req)
		switch {
		case errors.Is(err, ErrRecordNotFound):
			w.WriteHeader(http.StatusNotFound)
			logger.Info("handle", slog.String("error", "not found"))
		case err != nil:
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("handle", slog.Any("error", err))
		default:
			rb, err := json.Marshal(res)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				logger.Error("write body", slog.Any("error", err))
				return
			}
			logger.Info("response", slog.Any("body", rb))
			if _, err := w.Write(rb); err != nil {
				logger.Error("write body", slog.Any("error", err))
			}
		}
	}
}

func RedirectHandler(redirector Redirector, pattern string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		name := strings.TrimPrefix(r.URL.Path, pattern)
		logger := slog.With(slog.String("url", r.URL.String()), slog.String("name", name))
		res, err := redirector.Redirect(r.Context(), &RedirectRequest{
			Name: name,
		})
		switch {
		case errors.Is(err, ErrRecordNotFound):
			w.WriteHeader(http.StatusNotFound)
			logger.Info("hanle", slog.String("error", "not found"))
		case err != nil:
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("handle", slog.Any("error", err))
		default:
			w.Header().Set("Location", res.To)
			w.WriteHeader(http.StatusMovedPermanently)
			logger.Info("handle", slog.String("to", res.To))
		}
	}
}

func StatusHandler(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "OK")
}
