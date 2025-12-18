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

	ticker := time.NewTicker(time.Minute)

	for range ticker.C {
		log.Println("Updating model list...")

		list, err := FetchModels(cfg)
		if err != nil {
			log.Warnf("Failed to fetch models: %v\n", err)

			continue
		}

		newer := GetNewModels(models, list)

		if len(newer) > 0 {
			log.Printf("%d new models\n", len(newer))

			err = Notify(cfg, newer)
			if err != nil {
				log.Warnf("Failed to notify: %v\n", err)

				continue
			}
		} else {
			log.Println("Nothing new")
		}

		models = list
	}
}
