package repository

import (
	"encoding/json"
	"fmt"
	"hot-coffee/internal/customErrors"
	"io"
	"log/slog"
	"os"
	"regexp"
)

func readJSON(filePath string) ([]byte, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY|os.O_CREATE, 0o644)
	if err != nil {
		slog.Error("repo error: opening JSON", "filePath", filePath, "error", err)
		return nil, fmt.Errorf("%w: %s", customErrors.ErrJsonOpen, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		slog.Error("repo error: reading JSON file", "filePath", filePath, "error", err)
		return nil, fmt.Errorf("%w: %s", customErrors.ErrJsonRead, err)
	}

	if len(data) <= 0 {
		slog.Warn("repo warning: JSON file is empty", "filePath", filePath)
		data = []byte("[]")
		// return nil, fmt.Errorf("the JSON file is empty")
	}

	return data, nil
}

func saveJSONToFile(filePath string, data any) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		slog.Error("repo error: encoding JSON", "error", err, "filePath", filePath)
		return fmt.Errorf("%w: %s", customErrors.ErrJsonMarshal, err)
	}

	if err := os.WriteFile(filePath, jsonData, 0o644); err != nil {
		slog.Error("repo error: writing JSON to file", "error", err, "filePath", filePath)
		return fmt.Errorf("%w: %s", customErrors.ErrJsonWrite, err)
	}
	return nil
}

func isValidID(id string) bool {
	ValidBucketPath := regexp.MustCompile(`^order[1-9][0-9]*$`)
	return ValidBucketPath.MatchString(id)
}

func maxID(nums []int) int {
	maxNum := nums[0]
	for _, num := range nums[1:] {
		if num > maxNum {
			maxNum = num
		}
	}
	return maxNum
}
