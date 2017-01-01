package main

import (
	"fmt"
	"log"
	"time"

	"github.com/alecthomas/kingpin"
	alfred "github.com/pascalw/go-alfred"
)

func configurePipelinesCommand(app *kingpin.Application) {
	var keywords []string

	cmd := app.Command("pipelines", "List Pipelines.")

	cmd.Arg("filter", "Keywords to apply as filters").
		StringsVar(&keywords)

	cmd.Action(func(c *kingpin.ParseContext) error {
		response := alfred.NewResponse()
		defer response.Print()

		alfred.InitTerms(keywords)

		pipelines, err := listPipelinesFromDB()
		if err != nil {
			response.AddItem(alfredError(err))
			return nil
		}

		for _, pipeline := range pipelines {
			log.Printf("Comparing %s with %s", keywords, pipeline.FullName())
			if alfred.MatchesTerms(keywords, pipeline.FullName()) {
				response.AddItem(&alfred.AlfredResponseItem{
					Valid:    true,
					Uid:      pipeline.URL,
					Title:    pipeline.FullName(),
					Subtitle: pipeline.Description,
					Arg:      pipeline.URL,
				})
			}
		}

		return nil
	})
}

type Pipeline struct {
	URL, Name, Org, Description string
	Created, LastUpdated        time.Time
}

func (r Pipeline) FullName() string {
	return fmt.Sprintf("%s/%s", r.Org, r.Name)
}

func listPipelinesFromDB() ([]Pipeline, error) {
	db, err := OpenDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query("SELECT id, url,description, name, org, created_at, updated_at FROM pipeline")
	if err != nil {
		return nil, err
	}

	pipelines := []Pipeline{}

	for rows.Next() {
		var id, url, descr, name, org string
		var created, updated time.Time
		err = rows.Scan(&id, &url, &descr, &name, &org, &created, &updated)
		if err != nil {
			return nil, err
		}

		pipelines = append(pipelines, Pipeline{
			URL:         url,
			Name:        name,
			Org:         org,
			Description: descr,
			Created:     created,
			LastUpdated: updated,
		})
	}

	return pipelines, nil
}
