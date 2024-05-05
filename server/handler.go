package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"logpuller/pkg/model"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// Function to retrieve log files within the specified time range
func retrieveLogFiles(from, to time.Time) ([]string, error) {
	var logFiles []string

	// Iterate over the hours within the specified time range
	for t := from; !t.After(to); t = t.Add(time.Hour) {
		// Construct the log file name based on the hour
		filename := t.Format("2006-01-02/15.txt") // Format: yyyy-mm-dd/HH.txt
		folderName := t.Format("2006-01-02")
		// Check if the folder already exists
		if _, err := os.Stat(folderName); os.IsNotExist(err) {
			// Create the folder
			err := os.MkdirAll(folderName, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
		// Retrieve the log file contents from remote storage
		err := downloadFromS3(filename)
		if err != nil {
			return nil, err
		}
		// Simulate retrieval of log file contents (replace with actual implementation)
		logFiles = append(logFiles, filename)
	}
	return logFiles, nil
}

func downloadFromS3(item string) error {
	// NOTE: you need to store your AWS credentials in ~/.aws/credentials
	// 1) Define your bucket and item names
	bucket := os.Getenv("AWS_S3_BUCKET_NAME")
	if bucket == "" {
		return errors.New("bucket name not found in environment")
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		return errors.New("region name not found in environment")
	}
	// 2) Create an AWS session
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	// 3) Create a new AWS S3 downloader
	downloader := s3manager.NewDownloader(sess)

	// 4) Download the item from the bucket. If an error occurs, log it and exit. Otherwise, notify the user that the download succeeded.
	file, err := os.Create(item)
	if err != nil {
		log.Printf("Unable to create file %q, %v", item, err)
		return nil
	}
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})

	return err
}

func logsearch(w http.ResponseWriter, r *http.Request) {
	logger := log.New(os.Stdout, "logpuller: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("Handling log search request...")

	// Context with timeout
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var requestData model.LogSearchRequest
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		writeErrorResponse(w, "failed to parse request body", err, http.StatusBadRequest)
		logger.Printf("Failed to parse request body: %v\n", err)
		return
	}

	logger.Printf("Received search request: %+v\n", requestData)

	// Retrieve log files within the specified time range
	logFiles, err := retrieveLogFiles(requestData.From, requestData.To)
	if err != nil {
		writeErrorResponse(w, "failed to retrieve log files", err, http.StatusBadRequest)
		logger.Printf("Failed to retrieve log files: %v\n", err)
		return
	}

	logger.Printf("Retrieved %d log files\n", len(logFiles))

	var results []string
	for _, log := range logFiles {
		result, _ := grep(log, requestData.SearchKeyword)
		results = append(results, result...)
		// Split the filename by "/"
		parts := strings.Split(log, "/")
		// Defer the deletion of the folder
		defer func() {
			if err := os.RemoveAll(parts[0]); err != nil {
				logger.Printf("Error deleting folder: %v\n", err)
				fmt.Println("Error deleting folder:", err) // Also print to console for immediate visibility
			}
		}()
	}

	writeSuccessResponse(w, model.ShortnerResponse{
		LogLines: results,
	})
}
