package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type httpServer struct {
	Log *Log
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

type ProduceRequest struct {
	Record Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}

func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request){
	var req ProduceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w,err.Error(), http.StatusBadRequest)
		return
	}

	offset,err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := ProduceResponse{Offset: offset}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
        return
	}
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record,err := s.Log.Read(req.Offset)
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func NewHTTPServer(addr string) *http.Server {
	httpSrv := newHTTPServer()
	r := mux.NewRouter()
	r.HandleFunc("/", httpSrv.handleProduce).Methods("POST")
	r.HandleFunc("/", httpSrv.handleConsume).Methods("GET")
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
