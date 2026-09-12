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

	"github.com/coalaura/openingrouter"
)

func Notify(cfg *Config, list []openingrouter.FrontendModel) error {
	embeds := make([]map[string]any, 0, len(list))

	for index := range list {
		model := &list[index]

		created := model.CreatedAt.Time

		contextLength := strconv.Itoa(model.ContextLength)

		var (
			promptPrice     float64
			completionPrice float64
		)

		if model.Endpoint != nil && model.Endpoint.Pricing != nil {
			promptPrice = model.Endpoint.Pricing.Prompt.Float64() * 1000000
			completionPrice = model.Endpoint.Pricing.Completion.Float64() * 1000000
		}

		description := fmt.Sprintf(
			"```\nModality: %s\nContext:  %s tokens\nPricing:  $%s 🡒 $%s\n```\n\n*%s*",
			Modalities(model.InputModalities, model.OutputModalities),
			contextLength,
			strconv.FormatFloat(promptPrice, 'f', -1, 64),
			strconv.FormatFloat(completionPrice, 'f', -1, 64),
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

	payload := map[string]any{
		"embeds": embeds,
	}

	body, err := json.Marshal(payload)
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
