package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Notify(cfg *Config, list []Model) error {
	var embeds []map[string]any

	for _, model := range list {
		created := time.Unix(model.CreatedAt, 0)

		ctx := "N/A"

		if model.Context != nil {
			ctx = strconv.FormatInt(*model.Context, 10)
		}

		description := fmt.Sprintf(
			"```\nModality: %s\nContext:  %s tokens\nPricing:  $%s 🡒 $%s\n```\n\n*%s*",
			model.Modality,
			ctx,
			strconv.FormatFloat(model.Pricing.Input, 'f', -1, 64),
			strconv.FormatFloat(model.Pricing.Output, 'f', -1, 64),
			strings.TrimSpace(model.Description),
		)

		embeds = append(embeds, map[string]any{
			"author": map[string]any{
				"name": "New OpenRouter Model",
				"url":  fmt.Sprintf("https://openrouter.ai/%s", model.Slug),
			},
			"title":       model.Name,
			"description": description,
			"color":       12041720,
			"footer": map[string]any{
				"text": "OpenMon",
			},
			"timestamp": created.Format(time.RFC3339Nano),
		})
	}

	body, err := json.Marshal(map[string]any{
		"embeds": embeds,
	})

	if err != nil {
		return err
	}

	resp, err := http.Post(cfg.Webhook, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}

	resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusNoContent {
		return errors.New(resp.Status)
	}

	return nil
}
