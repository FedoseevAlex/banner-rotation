package server

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/FedoseevAlex/banner-rotation/internal/common"
	"github.com/FedoseevAlex/banner-rotation/internal/config"
	"github.com/FedoseevAlex/banner-rotation/internal/types"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	app        types.Application
	httpServer *http.Server
	timeout    time.Duration
}

func NewServer(application types.Application, cfg config.Server) (*Server, error) {
	mux := httprouter.New()
	httpServer := http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: mux,
	}

	timeout, err := time.ParseDuration(cfg.Timeout)
	if err != nil {
		return nil, err
	}

	server := &Server{
		app:        application,
		httpServer: &httpServer,
		timeout:    timeout,
	}
	requestLogger := server.app.GetLogger("request info")

	mux.Handle("POST", "/banners", loggingMiddleware(
		server.addBannerHandler,
		requestLogger,
	))
	mux.Handle("DELETE", "/banners/:banner_id", loggingMiddleware(
		server.deleteBannerHandler,
		requestLogger,
	))
	mux.Handle("GET", "/banners/:banner_id", loggingMiddleware(
		server.getBannerHandler,
		requestLogger,
	))
	mux.Handle("POST", "/group/:group_id/slots/:slot_id/banners/:banner_id/click", loggingMiddleware(
		server.registerClickHandler,
		requestLogger,
	))
	mux.Handle("GET", "/group/:group_id/slots/:slot_id/banners/:banner_id/stats", loggingMiddleware(
		server.getStatsHandler,
		requestLogger,
	))
	mux.Handle("GET", "/group/:group_id/slots/:slot_id/banner", loggingMiddleware(
		server.chooseBannerHandler,
		requestLogger,
	))
	mux.Handle("GET", "/version", server.versionHandler)

	return server, nil
}

func jsonResponse(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bodyBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) versionHandler(w http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-type", "application/json")
	if err := common.PrintVersion(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) addBannerHandler(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	body := DescriptionBody{}
	err := json.NewDecoder(request.Body).Decode(&body)
	if err != nil {
		jsonResponse(
			w,
			http.StatusBadRequest,
			BadRequestResponse{
				Error: err.Error(),
				Msg:   "failed to decode request body",
			},
		)
		return
	}

	ctx, cancel := context.WithTimeout(request.Context(), s.timeout)
	defer cancel()

	banner, err := s.app.AddBanner(ctx, body.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, banner)
}

func (s *Server) deleteBannerHandler(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	bannerID, err := uuid.Parse(params.ByName("banner_id"))
	if err != nil {
		jsonResponse(
			w,
			http.StatusBadRequest,
			BadRequestResponse{
				Error: err.Error(),
				Msg:   "failed to parse banner uuid",
			},
		)
		return
	}

	ctx, cancel := context.WithTimeout(request.Context(), s.timeout)
	defer cancel()

	err = s.app.DeleteBanner(ctx, bannerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusNoContent, nil)
}

func (s *Server) getBannerHandler(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	bannerID, err := uuid.Parse(params.ByName("banner_id"))
	if err != nil {
		jsonResponse(
			w,
			http.StatusBadRequest,
			BadRequestResponse{
				Error: err.Error(),
				Msg:   "failed to parse banner uuid",
			},
		)
		return
	}

	ctx, cancel := context.WithTimeout(request.Context(), s.timeout)
	defer cancel()

	banner, err := s.app.GetBanner(ctx, bannerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, banner)
}

func (s *Server) registerClickHandler(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	bannerID, err := uuid.Parse(params.ByName("banner_id"))
	if err != nil {
		jsonResponse(
			w,
			http.StatusBadRequest,
			BadRequestResponse{
				Error: err.Error(),
				Msg:   "failed to parse banner uuid",
			},
		)
		return
	}

	slotID, err := uuid.Parse(params.ByName("slot_id"))
	if err != nil {
		jsonResponse(
			w,
			http.StatusBadRequest,
			BadRequestResponse{
				Error: err.Error(),
				Msg:   "failed to parse slot uuid",
			},
		)
		return
	}

	groupID, err := uuid.Parse(params.ByName("group_id"))
	if err != nil {
		jsonResponse(
			w,
			http.StatusBadRequest,
			BadRequestResponse{
				Error: err.Error(),
				Msg:   "failed to parse group uuid",
			},
		)
		return
	}

	ctx, cancel := context.WithTimeout(request.Context(), s.timeout)
	defer cancel()

	err = s.app.RegisterClick(ctx, bannerID, slotID, groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusNoContent, nil)
}

func (s *Server) getStatsHandler(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	bannerID, err := uuid.Parse(params.ByName("banner_id"))
	if err != nil {
		jsonResponse(
			w,
			http.StatusBadRequest,
			BadRequestResponse{
				Error: err.Error(),
				Msg:   "failed to parse banner uuid",
			},
		)
		return
	}

	slotID, err := uuid.Parse(params.ByName("slot_id"))
	if err != nil {
		jsonResponse(
			w,
			http.StatusBadRequest,
			BadRequestResponse{
				Error: err.Error(),
				Msg:   "failed to parse slot uuid",
			},
		)
		return
	}

	groupID, err := uuid.Parse(params.ByName("group_id"))
	if err != nil {
		jsonResponse(
			w,
			http.StatusBadRequest,
			BadRequestResponse{
				Error: err.Error(),
				Msg:   "failed to parse group uuid",
			},
		)
		return
	}

	ctx, cancel := context.WithTimeout(request.Context(), s.timeout)
	defer cancel()

	stats, err := s.app.GetStats(ctx, bannerID, slotID, groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusNoContent, stats)
}

func (s *Server) chooseBannerHandler(w http.ResponseWriter, request *http.Request, params httprouter.Params) {
	slotID, err := uuid.Parse(params.ByName("slot_id"))
	if err != nil {
		jsonResponse(
			w,
			http.StatusBadRequest,
			BadRequestResponse{
				Error: err.Error(),
				Msg:   "failed to parse slot uuid",
			},
		)
		return
	}

	groupID, err := uuid.Parse(params.ByName("group_id"))
	if err != nil {
		jsonResponse(
			w,
			http.StatusBadRequest,
			BadRequestResponse{
				Error: err.Error(),
				Msg:   "failed to parse group uuid",
			},
		)
		return
	}

	ctx, cancel := context.WithTimeout(request.Context(), s.timeout)
	defer cancel()

	rotation, err := s.app.ChooseBanner(ctx, slotID, groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, rotation)
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.httpServer.Shutdown(context.Background())
}
