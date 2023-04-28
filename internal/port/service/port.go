package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/bcicen/jstream"
	"github.com/boriszhilko/ports-service/internal/port/domain"
)

var (
	ErrRepositoryRequired = errors.New("nil repository")
	ErrInputRequired      = errors.New("nil input")
	batchSize             = 100
)

type Service interface {
	CreateOrUpdatePorts(ctx context.Context) error
}

type PortService struct {
	input      domain.PortInput
	repository domain.PortRepository
}

func NewPortService(input domain.PortInput, pr domain.PortRepository) (Service, error) {
	if pr == nil {
		return nil, ErrRepositoryRequired
	}
	if input == nil {
		return nil, ErrInputRequired
	}
	return &PortService{
		repository: pr,
		input:      input,
	}, nil
}

func (p *PortService) CreateOrUpdatePorts(ctx context.Context) error {
	batch := make([]domain.Port, 0, batchSize)
	data := p.input.GetData()

	decoder := jstream.NewDecoder(data, 1).EmitKV()
	for val := range decoder.Stream() {
		kv := val.Value.(jstream.KV)
		value := kv.Value.(map[string]interface{})
		key := kv.Key
		port := p.toModel(key, value)
		batch = append(batch, port)

		if len(batch) == batchSize {
			err := p.mergeAndWrite(ctx, batch)
			if err != nil {
				return err
			}
			batch = batch[:0]
		}
	}

	if len(batch) > 0 {
		if err := p.mergeAndWrite(ctx, batch); err != nil {
			return err
		}
	}

	return nil
}

func (p *PortService) mergeAndWrite(ctx context.Context, batch []domain.Port) error {
	mergedBatch, err := p.mergeWithExisting(ctx, batch)
	if err != nil {
		return err
	}
	if err := p.repository.CreateOrUpdate(ctx, mergedBatch); err != nil {
		return err
	}
	return nil
}

func (p *PortService) mergeWithExisting(ctx context.Context, ports []domain.Port) ([]domain.Port, error) {
	mergedPorts := make([]domain.Port, 0, len(ports))

	for _, record := range ports {
		existingRecord, err := p.repository.Get(ctx, record.ID)
		if err != nil {
			if err == domain.ErrPortNotFound {
				mergedPorts = append(mergedPorts, record)
				continue
			}
			return mergedPorts, err
		}
		mergedPorts = append(mergedPorts, updateFields(existingRecord, record))
	}
	return mergedPorts, nil
}

func (p *PortService) toModel(key string, portData map[string]interface{}) domain.Port {
	return domain.Port{
		ID:          key,
		Name:        p.toString(portData["name"]),
		City:        p.toString(portData["city"]),
		Country:     p.toString(portData["country"]),
		Alias:       p.toStringSlice(portData["alias"]),
		Regions:     p.toStringSlice(portData["regions"]),
		Coordinates: p.toFloatSlice(portData["coordinates"]),
		Province:    p.toString(portData["province"]),
		Timezone:    p.toString(portData["timezone"]),
		Unlocs:      p.toStringSlice(portData["unlocs"]),
		Code:        p.toString(portData["code"]),
	}
}

func (p *PortService) toString(s interface{}) string {
	return fmt.Sprintf("%v", s)
}

func (p *PortService) toFloatSlice(val interface{}) []float64 {
	if val == nil {
		return nil
	}
	var result []float64
	for _, v := range val.([]interface{}) {
		result = append(result, v.(float64))
	}
	return result
}

func (p *PortService) toStringSlice(val interface{}) []string {
	if val == nil {
		return nil
	}
	var result []string
	for _, v := range val.([]interface{}) {
		result = append(result, v.(string))
	}
	return result
}

func updateFields(existing domain.Port, updates domain.Port) domain.Port {
	if updates.Name != "" {
		existing.Name = updates.Name
	}
	if updates.City != "" {
		existing.City = updates.City
	}
	if updates.Country != "" {
		existing.Country = updates.Country
	}
	if len(updates.Alias) > 0 {
		existing.Alias = updates.Alias
	}
	if len(updates.Regions) > 0 {
		existing.Regions = updates.Regions
	}
	if len(updates.Coordinates) > 0 {
		existing.Coordinates = updates.Coordinates
	}
	if updates.Province != "" {
		existing.Province = updates.Province
	}
	if updates.Timezone != "" {
		existing.Timezone = updates.Timezone
	}
	if len(updates.Unlocs) > 0 {
		existing.Unlocs = updates.Unlocs
	}
	if updates.Code != "" {
		existing.Code = updates.Code
	}
	return existing
}
