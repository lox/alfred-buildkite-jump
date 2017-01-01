package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/buildkite/go-buildkite/buildkite"
)

const (
	apiTokenEnvKey = "BUILDKITE_API_TOKEN"
)

func configureUpdateCommand(app *kingpin.Application) {
	var apiToken string

	cmd := app.Command("update", "Update Pipelines from Buildkite API")

	cmd.Flag("token", "Buildkite API Token").
		OverrideDefaultFromEnvar(apiTokenEnvKey).
		StringVar(&apiToken)

	cmd.Action(func(c *kingpin.ParseContext) error {
		if apiToken == "" {
			fmt.Printf("Error: No %s set, configure in environment settings", apiTokenEnvKey)
			os.Exit(1)
			return nil
		}

		n, err := updatePipelines(apiToken)
		if err != nil {
			fmt.Println("Error", err)
			os.Exit(1)
		}

		fmt.Printf("Updated %d pipelines from Buildkite", n)
		return nil
	})
}

func listPipelines(apiToken string) ([]Pipeline, error) {
	config, err := buildkite.NewTokenConfig(apiToken, true)
	if err != nil {
		return nil, err
	}

	client := buildkite.NewClient(config.Client())

	orgs, _, err := client.Organizations.List(nil)
	if err != nil {
		return nil, err
	}

	results := []Pipeline{}

	for _, org := range orgs {
		log.Printf("Listing builds for org %s", *org.Slug)
		ps, _, err := client.Pipelines.List(*org.Slug, nil)
		if err != nil {
			return nil, err
		}
		for _, p := range ps {
			results = append(results, Pipeline{
				Org:  *org.Slug,
				Name: *p.Slug,
				URL:  *p.WebURL,
			})
		}
	}

	return results, nil
}

func nilableString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func buildkiteTime(t *buildkite.Timestamp) *time.Time {
	if t == nil {
		return nil
	}
	return &t.Time
}

func updatePipelines(apiToken string) (int64, error) {
	pipelines, err := listPipelines(apiToken)
	if err != nil {
		return 0, err
	}

	db, err := OpenDB()
	if err != nil {
		return 0, err
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	found := map[string]struct{}{}
	counter := int64(0)

	for _, pipeline := range pipelines {
		id := pipeline.FullName()
		log.Printf("Updating %s", id)
		res, err := db.Exec(
			`INSERT OR REPLACE INTO pipeline (
					id,
					url,
					description,
					name, org,
					updated_at,
					created_at
				) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			id,
			pipeline.URL,
			pipeline.Description,
			pipeline.Name,
			pipeline.Org,
			pipeline.LastUpdated,
			pipeline.Created,
		)
		if err != nil {
			return counter, err
		}
		found[id] = struct{}{}
		rows, _ := res.RowsAffected()
		counter += rows
	}

	existing, err := listPipelinesFromDB()
	if err != nil {
		return 0, err
	}

	// purge pipelines that don't exit any more
	for _, pipeline := range existing {
		if _, exists := found[pipeline.FullName()]; !exists {
			log.Printf("Pipeline %s doesn't exist, deleting", pipeline.FullName())

			_, err := db.Exec(
				`DELETE FROM pipeline WHERE id=?`,
				pipeline.FullName(),
			)
			if err != nil {
				return 0, err
			}

		}
	}

	return counter, tx.Commit()
}
