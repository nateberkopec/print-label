package main

import (
	"bytes"
	"encoding/base64"
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
	Images    []string `json:"images"`
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
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 8<<20)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	cfg := s.cfg
	if req.TapeWidth != 0 {
		cfg.TapeWidth = req.TapeWidth
	}
	if err := validate(req, cfg); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	images, status, err := buildImages(req, cfg)
	if err != nil {
		writeError(w, status, err)
		return
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

func validate(req printRequest, cfg printlabel.Config) error {
	if _, ok := printlabel.TapeInfos[cfg.TapeWidth]; !ok {
		return fmt.Errorf("unsupported tape width: %gmm", cfg.TapeWidth)
	}
	if req.Labels != nil && req.Images != nil {
		return errors.New("send either labels or images, not both")
	}
	if len(req.Labels) == 0 && len(req.Images) == 0 {
		return errors.New("at least one label or image is required")
	}
	for _, text := range req.Labels {
		if strings.TrimSpace(text) == "" {
			return errors.New("labels must not be blank")
		}
	}
	return nil
}

func buildImages(req printRequest, cfg printlabel.Config) ([]*image.Gray, int, error) {
	if len(req.Images) > 0 {
		images := make([]*image.Gray, len(req.Images))
		for i, encoded := range req.Images {
			data, err := base64.StdEncoding.DecodeString(encoded)
			if err != nil {
				return nil, http.StatusBadRequest, fmt.Errorf("image %d: %w", i+1, err)
			}
			if images[i], err = printlabel.LoadImage(bytes.NewReader(data), cfg.TapeWidth); err != nil {
				return nil, http.StatusBadRequest, fmt.Errorf("image %d: %w", i+1, err)
			}
		}
		return images, 0, nil
	}
	images := make([]*image.Gray, len(req.Labels))
	for i, text := range req.Labels {
		var err error
		if images[i], err = printlabel.RenderText(text, cfg); err != nil {
			return nil, http.StatusInternalServerError, err
		}
	}
	return images, 0, nil
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
