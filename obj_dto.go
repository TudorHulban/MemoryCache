package cache

import "time"

// DomainModel as per DDD.
type DomainModel string

// DTO as data transfer object brings together the application layers and the cache.
type DTO struct {
	Key         int64
	DomainModel DomainModel // cannot be nil
	Data        []byte
}

// NewDTO constructor should be used to get a DTO from domain model data.
func NewDTO(domain DomainModel, data []byte) *DTO {
	return &DTO{
		Key:         time.Now().UnixNano(),
		DomainModel: domain,
		Data:        data,
	}
}
