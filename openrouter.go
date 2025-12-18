package main

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/revrost/go-openrouter"
)

type ModelPricing struct {
	Input  float64
	Output float64
}

type Model struct {
	Slug        string
	Name        string
	Description string
	Context     *int64
	Modality    string
	CreatedAt   int64
	Pricing     ModelPricing
}

func FetchModels(cfg *Config) ([]Model, error) {
	clientCfg := openrouter.DefaultConfig(cfg.ApiKey)

	clientCfg.XTitle = "OpenMon"
	clientCfg.HttpReferer = "https://github.com/coalaura/openmon"

	client := openrouter.NewClientWithConfig(*clientCfg)

	list, err := client.ListModels(context.Background())
	if err != nil {
		return nil, err
	}

	models := make([]Model, len(list))

	for i, model := range list {
		input, err := strconv.ParseFloat(model.Pricing.Prompt, 64)
		if err != nil {
			return nil, err
		}

		output, err := strconv.ParseFloat(model.Pricing.Completion, 64)
		if err != nil {
			return nil, err
		}

		models[i] = Model{
			Slug:        model.ID,
			Name:        model.Name,
			Description: model.Description,
			Context:     model.ContextLength,
			Modality:    Modalities(model.Architecture.InputModalities, model.Architecture.OutputModalities),
			CreatedAt:   model.Created,
			Pricing: ModelPricing{
				Input:  input * 1000000,
				Output: output * 1000000,
			},
		}
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
