package repositoryimpl

import "github.com/opensourceways/xihe-server/common/infrastructure/pgsql"

type pgsqlClient interface {
	Create(result interface{}) error
	Updates(filter, result interface{}) error
	Count(filter interface{}) (int, error)
	Filter(filter, result interface{}) error
	First(filter, result interface{}) error
	GetRecords(filter, result interface{}, p pgsql.Pagination, sort []pgsql.SortByColumn) error
	GetRecord(filter, result interface{}) error
	UpdateRecord(filter, update interface{}) error

	IsRowNotFound(err error) bool
	IsRowExists(err error) bool
}
