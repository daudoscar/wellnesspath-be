package seeder

import (
	"context"
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"wellnesspath/config"
	"wellnesspath/models"
)

func SeedExercisesFromAzure() error {
	if config.DB == nil {
		log.Fatal("❌ Database is not connected.")
	}

	blobName := "exercises.csv"
	container := config.ENV.AzureContainerName

	downloadResp, err := config.BlobClient.DownloadStream(context.TODO(), container, blobName, nil)
	if err != nil {
		log.Printf("❌ Failed to download blob: %v", err)
		return err
	}
	defer downloadResp.Body.Close()

	body, err := ioutil.ReadAll(downloadResp.Body)
	if err != nil {
		log.Fatalf("❌ Failed to read blob stream: %v", err)
	}

	log.Printf("📄 Blob size: %d bytes", len(body))
	log.Printf("🔍 First 300 characters:\n%s", string(body[:min(300, len(body))]))

	reader := csv.NewReader(strings.NewReader(string(body)))
	reader.Comma = ';'

	// ✅ Skip the header row
	if _, err := reader.Read(); err != nil {
		log.Fatalf("❌ Failed to read CSV header: %v", err)
	}

	successCount := 0
	failCount := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("⚠️ Error reading row: %v", err)
			failCount++
			continue
		}

		log.Printf("📥 Row: %v", record)

		ex := models.Exercise{
			Name:                   safeGet(record, 0),
			BodyPart:               safeGet(record, 5),
			Equipment:              safeGet(record, 7),
			Description:            safeGet(record, 9),
			StepByStepInstructions: safeGet(record, 9),
			Difficulty:             safeGet(record, 10),
			Category:               safeGet(record, 11),
			ExerciseType:           safeGet(record, 12),
			GoalTag:                safeGet(record, 13),
		}

		if ex.Name == "" || ex.BodyPart == "" {
			log.Printf("⚠️ Skipped row (missing Name/BodyPart): %v", record)
			failCount++
			continue
		}

		if err := config.DB.Create(&ex).Error; err != nil {
			log.Printf("❌ Insert failed for '%s': %v", ex.Name, err)
			failCount++
		} else {
			successCount++
		}
	}

	log.Printf("✅ Finished seeding: %d inserted, %d failed.\n", successCount, failCount)
	return nil
}

func safeGet(record []string, index int) string {
	if len(record) > index {
		return strings.TrimSpace(record[index])
	}
	return ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
