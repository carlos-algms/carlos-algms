package cache

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"scraper/pkg/github"
)

func ComputeHash(issues, prs []github.Issue) string {
	data := map[string]interface{}{
		"issues": issues,
		"prs":    prs,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal cache data: %v", err)
	}
	hash := sha256.Sum256(jsonData)
	return fmt.Sprintf("%x", hash)
}

func HasChanged(newHash, filePath string) bool {
	oldHash, err := os.ReadFile(filePath)
	if err != nil {
		return true
	}
	return string(oldHash) != newHash
}

func Save(hash, filePath string) error {
	return os.WriteFile(filePath, []byte(hash), 0644)
}