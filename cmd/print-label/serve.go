package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"net/http"
	"strings"
	"sync"
	"time"

	printlabel "github.com/nateberkopec/print-label"
)

type printRequest struct {
	Labels    []string `json:"labels"`
	TapeWidth float64  `json:"tape_width"`
}

type printResponse struct {
	Labels int `json:"labels"`
	Bytes  int `json:"bytes"`
}

type printServer struct {
	cfg     printlabel.Config
	printer sync.Mutex
}

func serve(addr string, cfg printlabel.Config) error {
	s := &printServer{cfg: cfg}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("POST /print", s.handlePrint)
	fmt.Printf("Listening on %s\n", addr)
	return (&http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 10 * time.Second}).ListenAndServe()
}

func (s *printServer) handlePrint(w http.ResponseWriter, r *http.Request) {
	var req printRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	cfg := s.cfg
	if req.TapeWidth != 0 {
		cfg.TapeWidth = req.TapeWidth
	}
	if err := validate(req.Labels, cfg); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	images := make([]*image.Gray, len(req.Labels))
	for i, text := range req.Labels {
		var err error
		if images[i], err = printlabel.RenderText(text, cfg); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
	}
	job, err := printlabel.BuildLabelJob(images, cfg)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	s.printer.Lock()
	err = printlabel.SendToPrinter(job, cfg.PrinterIP, cfg.PrinterPort)
	s.printer.Unlock()
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	fmt.Printf("Sent %d labels (%d bytes) on %gmm tape\n", len(images), len(job), cfg.TapeWidth)
	writeJSON(w, http.StatusOK, printResponse{Labels: len(images), Bytes: len(job)})
}

func validate(labels []string, cfg printlabel.Config) error {
	if _, ok := printlabel.TapeInfos[cfg.TapeWidth]; !ok {
		return fmt.Errorf("unsupported tape width: %gmm", cfg.TapeWidth)
	}
	if len(labels) == 0 {
		return errors.New("at least one label is required")
	}
	for _, text := range labels {
		if strings.TrimSpace(text) == "" {
			return errors.New("labels must not be blank")
		}
	}
	return nil
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
