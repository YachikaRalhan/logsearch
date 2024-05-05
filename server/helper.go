package server

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"logpuller/pkg/model"
)

func grep(filename, searchTerm string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var results []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, searchTerm) {
			results = append(results, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func writeErrorResponse(w http.ResponseWriter, msg string, err error, statusCode int) {
	log.Println("got error: ", err)

	resp := model.ErrorResponse{
		Message: msg,
		Detail:  err.Error(),
		Status:  strconv.Itoa(http.StatusBadRequest),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}

func writeSuccessResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
