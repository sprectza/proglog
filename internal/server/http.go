package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func NewHTTPServer(addr string) *http.Server {
	httpserv := httpHanlderFunc()
	r := mux.NewRouter()

	r.HandleFunc("/", httpserv.handleProduce).Methods("POST")
	r.HandleFunc("/", httpserv.handleConsume).Methods("GET")

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

/*
	Standard way to handle the request coming to the server

a) Unmarshal the request's JSON into a struct
b) Run the endpoint's logic with request to get a result
c) Marshal and write that result back to response
*/
type LogHTTPHandler struct {
	Log *Log
}

func httpHanlderFunc() *LogHTTPHandler {
	return &LogHTTPHandler{
		Log: NewLog(),
	}
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ProduceRequest struct {
	Record Record `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"offset"`
}

type ConsumerRequest struct {
	Offset uint64 `json:"offset"`
}

func (s *LogHTTPHandler) handleProduce(w http.ResponseWriter, r *http.Request) {
	// First step: Decode the incoming JSON
	var request ProduceRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Second step: Call the bizz logic
	offset, err := s.Log.Append(request.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Third step: Encode the response and write it back
	var response ProduceResponse = ProduceResponse{Offset: offset}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *LogHTTPHandler) handleConsume(w http.ResponseWriter, r *http.Request) {
	// First step: Decode the incoming JSON payload
	var request ConsumerRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Second step: Call the bizz logic
	record, err := s.Log.Read(request.Offset)
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Third step: Encode the response and give it back
	var response ConsumeResponse = ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
