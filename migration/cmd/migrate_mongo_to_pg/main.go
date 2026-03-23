package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/leanote/leanote/migration"
	m2p "github.com/leanote/leanote/migration/mongodb_to_postgres"
)

func main() {
	mongoURL := flag.String("mongo-url", "mongodb://localhost:27017/leanote", "MongoDB connection URL")
	postgresURL := flag.String("postgres-url", "host=localhost port=5432 user=leanote password= dbname=leanote sslmode=disable", "PostgreSQL connection URL")
	batchSize := flag.Int("batch-size", 1000, "Batch size for migration")
	stopOnError := flag.Bool("stop-on-error", false, "Stop on error")
	validateAfter := flag.Bool("validate", true, "Validate after migration")
	verbose := flag.Bool("verbose", true, "Verbose output")
	flag.Parse()

	config := migration.MigrationConfig{
		BatchSize:     *batchSize,
		StopOnError:   *stopOnError,
		ValidateAfter: *validateAfter,
		Verbose:       *verbose,
	}

	migrator := m2p.NewMongoToPostgresMigrator(*mongoURL, *postgresURL, config)

	log.Println("Starting MongoDB to PostgreSQL migration...")
	log.Printf("MongoDB URL: %s", *mongoURL)
	log.Printf("PostgreSQL URL: %s", *postgresURL)

	err := migrator.Connect()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer migrator.Close()

	log.Println("Connected to both databases successfully")

	err = migrator.PreMigrate()
	if err != nil {
		log.Fatalf("Pre-migration failed: %v", err)
	}

	tables := []string{
		"users",
		"notebooks",
		"notes",
		"note_contents",
		"tags",
		"attachs",
		"files",
		"albums",
	}

	totalSuccess := 0
	totalFailed := 0

	for _, table := range tables {
		result, err := migrator.MigrateTable(table)
		if err != nil {
			log.Printf("Failed to migrate table %s: %v", table, err)
			if *stopOnError {
				os.Exit(1)
			}
			continue
		}

		totalSuccess += result.SuccessRecords
		totalFailed += result.FailedRecords

		log.Printf("Table %s: %d/%d records migrated in %v",
			table, result.SuccessRecords, result.TotalRecords, result.Duration)

		if len(result.Errors) > 0 {
			log.Printf("Errors in table %s:", table)
			for _, errMsg := range result.Errors {
				log.Printf("  - %s", errMsg)
			}
		}
	}

	err = migrator.PostMigrate()
	if err != nil {
		log.Printf("Post-migration warnings: %v", err)
	}

	if *validateAfter {
		log.Println("\n=== Validating Migration ===")
		result, err := migrator.Validate()
		if err != nil {
			log.Printf("Validation failed: %v", err)
		} else if result.Passed {
			log.Println("✓ All validations passed!")
		} else {
			log.Println("✗ Some validations failed!")
			for _, issue := range result.DataIntegrity {
				log.Printf("  - %s", issue)
			}
		}
	}

	log.Println("\n=== Migration Summary ===")
	log.Printf("Total Success: %d records", totalSuccess)
	log.Printf("Total Failed: %d records", totalFailed)

	if totalFailed > 0 {
		log.Println("⚠ Migration completed with errors")
		os.Exit(1)
	} else {
		log.Println("✓ Migration completed successfully!")
	}
}
