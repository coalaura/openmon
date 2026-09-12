package main

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/coalaura/openingrouter"
)

func FetchModels(cfg *Config) ([]openingrouter.FrontendModel, error) {
	list, err := openingrouter.ListFrontendModels(context.Background())
	if err != nil {
		return nil, err
	}

	models := make([]openingrouter.FrontendModel, 0, len(list))
	bySlug := make(map[string]int, len(list))

	for index := range list {
		model := &list[index]

		if model.Endpoint == nil {
			continue
		}

		if len(cfg.Providers.Include) > 0 && !slices.Contains(cfg.Providers.Include, model.Author) {
			continue
		}

		if slices.Contains(cfg.Providers.Exclude, model.Author) {
			continue
		}

		existingIndex, exists := bySlug[model.Slug]
		if exists {
			existing := &models[existingIndex]

			if isBatchVariant(existing) && !isBatchVariant(model) {
				models[existingIndex] = *model
			}

			continue
		}

		bySlug[model.Slug] = len(models)

		models = append(models, *model)
	}

	sort.Slice(models, func(firstIndex, secondIndex int) bool {
		return models[firstIndex].CreatedAt.After(models[secondIndex].CreatedAt.Time)
	})

	return models, nil
}

func GetNewModels(seen map[string]struct{}, list []openingrouter.FrontendModel) []openingrouter.FrontendModel {
	newer := make([]openingrouter.FrontendModel, 0, len(list))

	batch := make(map[string]struct{}, len(list))

	for index := range list {
		model := &list[index]

		if _, exists := seen[model.Slug]; exists {
			continue
		}

		if model.Permaslug != "" {
			if _, exists := seen[model.Permaslug]; exists {
				continue
			}
		}

		if _, exists := batch[model.Slug]; exists {
			continue
		}

		batch[model.Slug] = struct{}{}

		if model.Permaslug != "" {
			batch[model.Permaslug] = struct{}{}
		}

		newer = append(newer, *model)
	}

	return newer
}

func Modalities(inputModalities, outputModalities []string) string {
	return fmt.Sprintf(
		"%s 🡒 %s",
		strings.Join(inputModalities, "+"),
		strings.Join(outputModalities, "+"),
	)
}

func isBatchVariant(model *openingrouter.FrontendModel) bool {
	if model.Endpoint == nil {
		return false
	}

	if model.Endpoint.Variant == "batch" {
		return true
	}

	return strings.HasSuffix(model.Endpoint.ModelVariantSlug, ":batch")
}
