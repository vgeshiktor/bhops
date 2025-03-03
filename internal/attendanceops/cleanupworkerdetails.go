package attendanceops

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	// Define the path to the JSON file
	filePath := "config/id2worker.json"

	// create archive directory if needed
	err := os.MkdirAll("archive", os.ModePerm)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			panic(err)
		}
	}

	// copy file to archive directory
	timestamp := time.Now().Format("2006-01-02T15:04:05")
	archivePath := fmt.Sprintf(
		"archive/id2worker%s.json", timestamp)
	err = os.Rename(filePath, archivePath)
	if err != nil {
		log.Error().Msgf("failed to rename file: %s to %s, error: %v", filePath, archivePath, err)
		return
	}

	// open worker details file
	file, err := os.Open(filePath)
	if err != nil {
		log.Error().Msgf(
			"failed to open worker details file: %s, error: %v", filePath, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error().Msgf("failed to close worker details file: %s, error: %v",
				filePath, err)
		}
	}()

	// read worker details file
	content, err := io.ReadAll(file)
	if err != nil {
		log.Error().Msgf("Error reading file: %v", err)
		return
	}

	// Unmarshal the JSON data into a map
	var workers map[string]WorkerDetails
	err = json.Unmarshal( content, &workers)
	if err != nil {
		log.Error().Msgf(
			"Error unmarshalling JSON: %v", err)
		return
	}

	// Update some values in the map
	for id, details := range workers {
		details.Holidays = 0
		details.HolidayPresent = 0
		details.HoursAdjustment = 0
		details.Hours125Adjustment = 0
		workers[id] = details
	}

	// Marshal the map back to JSON
	updatedData, err := 
		json.MarshalIndent(workers, "", "  ")
	if err != nil {
		log.Error().Msgf(
			"Error marshalling JSON: %v", err)
		return
	}
	_ = updatedData

	// Write the updated JSON back to the file
	// err = io.WriteAll(filePath, updatedData, 0644)
	// if err != nil {
	// 	fmt.Println("Error writing file:", err)
	// 	return
	// }

	fmt.Println("File updated successfully")
}