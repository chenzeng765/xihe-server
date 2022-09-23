package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type TagsMapper interface {
	List(string) ([]DomainTagsDo, error)
}

type DomainTagsDo = domain.DomainTags
type TagsDo = domain.Tags

func NewTagsRepository(mapper TagsMapper) repository.Tags {
	return tags{mapper}
}

type tags struct {
	mapper TagsMapper
}

func (impl tags) List(resourceType domain.ResourceType) ([]domain.DomainTags, error) {
	v, err := impl.mapper.List(resourceType.ResourceType())
	if err != nil {
		return nil, convertError(err)
	}

	return v, nil
}