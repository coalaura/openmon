package main

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/coalaura/openingrouter"
)

type ModelPricing struct {
	Input  float64
	Output float64
}

type Model struct {
	Slug        string
	Name        string
	Description string
	Context     int64
	Modality    string
	CreatedAt   int64
	Pricing     ModelPricing
}

func FetchModels(cfg *Config) ([]Model, error) {
	list, err := openingrouter.ListFrontendModels(context.Background())
	if err != nil {
		return nil, err
	}

	models := make([]Model, 0, len(list))

	for _, model := range list {
		if model.Endpoint == nil {
			continue
		}

		if len(cfg.Providers.Include) > 0 && slices.Contains(cfg.Providers.Include, model.Author) {
			continue
		}

		if slices.Contains(cfg.Providers.Exclude, model.Author) {
			continue
		}

		models = append(models, Model{
			Slug:        model.Slug,
			Name:        model.Name,
			Description: model.Description,
			Context:     int64(model.ContextLength),
			Modality:    Modalities(model.InputModalities, model.OutputModalities),
			CreatedAt:   model.CreatedAt.Unix(),
			Pricing: ModelPricing{
				Input:  model.Endpoint.Pricing.Prompt.Float64() * 1000000,
				Output: model.Endpoint.Pricing.Completion.Float64() * 1000000,
			},
		})
	}

	sort.Slice(models, func(i, j int) bool {
		return models[i].CreatedAt > models[j].CreatedAt
	})

	return models, nil
}

func Modalities(in, out []string) string {
	return fmt.Sprintf(
		"%s 🡒 %s",
		strings.Join(in, "+"),
		strings.Join(out, "+"),
	)
}

func GetNewModels(prev, next []Model) []Model {
	if len(prev) == 0 {
		return next
	}

	newer := make([]Model, 0, len(next))

	latest := prev[0].CreatedAt

	for _, model := range next {
		if latest >= model.CreatedAt {
			break
		}

		newer = append(newer, model)
	}

	return newer
}
