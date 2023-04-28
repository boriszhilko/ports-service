package domain

import (
	"context"
	"encoding/json"
	"errors"
	"io"
)

var (
	ErrPortNotFound = errors.New("port not found")
)

type Port struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float64 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code"`
}

type PortInput interface {
	GetData() io.Reader
}

type PortRepository interface {
	CreateOrUpdate(ctx context.Context, port []Port) error
	Get(ctx context.Context, id string) (Port, error)
}

func (p Port) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}
