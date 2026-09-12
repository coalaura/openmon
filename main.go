package main

import (
	"time"

	"github.com/coalaura/plain"
)

var log = plain.New(plain.WithDate(plain.RFC3339Local))

func main() {
	log.Println("Loading config...")

	cfg, err := LoadConfig()
	log.MustFail(err)

	log.Println("Fetching initial list...")

	models, err := FetchModels(cfg)
	log.MustFail(err)

	log.Printf("Loaded %d models\n", len(models))

	seen := make(map[string]struct{}, len(models))

	for index := range models {
		seen[models[index].Slug] = struct{}{}

		if models[index].Permaslug != "" {
			seen[models[index].Permaslug] = struct{}{}
		}
	}

	ticker := time.NewTicker(time.Minute)

	for range ticker.C {
		log.Println("Updating model list...")

		list, err := FetchModels(cfg)
		if err != nil {
			log.Warnf("Failed to fetch models: %v\n", err)

			continue
		}

		newer := GetNewModels(seen, list)

		if len(newer) > 0 {
			log.Printf("%d new models\n", len(newer))

			err = Notify(cfg, newer)
			if err != nil {
				log.Warnf("Failed to notify: %v\n", err)

				continue
			}

			for index := range newer {
				seen[newer[index].Slug] = struct{}{}

				if newer[index].Permaslug != "" {
					seen[newer[index].Permaslug] = struct{}{}
				}
			}
		} else {
			log.Println("Nothing new")
		}
	}
}
