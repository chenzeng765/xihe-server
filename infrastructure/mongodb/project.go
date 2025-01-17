package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func NewProjectMapper(name string) repositories.ProjectMapper {
	return project{name}
}

type project struct {
	collectionName string
}

func (col project) newDoc(owner string) error {
	docFilter := resourceOwnerFilter(owner)

	doc := bson.M{
		fieldOwner: owner,
		fieldItems: bson.A{},
	}

	f := func(ctx context.Context) error {
		_, err := cli.newDocIfNotExist(
			ctx, col.collectionName, docFilter, doc,
		)

		return err
	}

	if err := withContext(f); err != nil && isDBError(err) {
		return err
	}

	return nil
}

func (col project) Insert(do repositories.ProjectDO) (identity string, err error) {
	if identity, err = col.insert(do); err == nil || !isDocNotExists(err) {
		return
	}

	// doc is not exist or duplicate insert

	if err = col.newDoc(do.Owner); err == nil {
		if identity, err = col.insert(do); err != nil && isDocNotExists(err) {
			err = repositories.NewErrorDuplicateCreating(err)
		}
	}

	return
}

func (col project) insert(do repositories.ProjectDO) (identity string, err error) {
	identity = newId()

	do.Id = identity
	doc, err := col.toProjectDoc(&do)
	if err != nil {
		return
	}
	doc[fieldVersion] = 0
	doc[fieldLikeCount] = 0
	doc[fieldForkCount] = 0
	doc[fieldDownloadCount] = 0
	doc[fieldModels] = bson.A{}
	doc[fieldDatasets] = bson.A{}

	err = insertResource(col.collectionName, do.Owner, do.Name, doc)

	return
}

func (col project) Delete(do *repositories.ResourceIndexDO) error {
	return deleteResource(col.collectionName, do)
}

func (col project) UpdateProperty(do *repositories.ProjectPropertyDO) error {
	p := &ProjectPropertyItem{
		Level:    do.Level,
		Name:     do.Name,
		FL:       do.FL,
		Desc:     do.Desc,
		CoverId:  do.CoverId,
		RepoType: do.RepoType,
		Tags:     do.Tags,
		TagKinds: do.TagKinds,
	}

	p.setDefault()

	return updateResourceProperty(col.collectionName, &do.ResourceToUpdateDO, p)
}

func (col project) Get(owner, identity string) (do repositories.ProjectDO, err error) {
	var v []dProject

	if err = getResourceById(col.collectionName, owner, identity, &v); err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	col.toProjectDO(owner, &v[0].Items[0], &do)

	return
}

func (col project) GetByName(owner, name string) (do repositories.ProjectDO, err error) {
	var v []dProject

	if err = getResourceByName(col.collectionName, owner, name, &v); err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	col.toProjectDO(owner, &v[0].Items[0], &do)

	return
}

func (col project) GetSummary(owner string, projectId string) (
	do repositories.ResourceSummaryDO, err error,
) {
	var v []dProject

	err = getResourceSummary(col.collectionName, owner, projectId, &v)
	if err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	item := &v[0].Items[0]
	do.Id = projectId
	do.Name = item.Name
	do.Owner = owner
	do.RepoId = item.RepoId
	do.RepoType = item.RepoType

	return
}

func (col project) GetSummaryByName(owner, name string) (
	do repositories.ResourceSummaryDO, err error,
) {
	var v []dProject

	err = getResourceSummaryByName(col.collectionName, owner, name, &v)
	if err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	item := &v[0].Items[0]
	do.Id = item.Id
	do.Name = name
	do.Owner = owner
	do.RepoId = item.RepoId
	do.RepoType = item.RepoType

	return
}

func (col project) ListUsersProjects(opts map[string][]string) (
	r []repositories.ProjectSummaryDO, err error,
) {
	var v []dProject

	err = listUsersResources(col.collectionName, opts, &v)
	if err != nil || len(v) == 0 {
		return
	}

	r = make([]repositories.ProjectSummaryDO, 0, len(v))

	for i := range v {
		owner := v[i].Owner
		items := v[i].Items

		dos := make([]repositories.ProjectSummaryDO, len(items))
		for j := range items {
			col.toProjectSummaryDO(owner, &items[j], &dos[j])
		}

		r = append(r, dos...)
	}

	return
}

func (col project) toProjectDoc(do *repositories.ProjectDO) (bson.M, error) {
	docObj := projectItem{
		Id:        do.Id,
		Type:      do.Type,
		Protocol:  do.Protocol,
		Training:  do.Training,
		RepoId:    do.RepoId,
		CreatedAt: do.CreatedAt,
		UpdatedAt: do.UpdatedAt,
		ProjectPropertyItem: ProjectPropertyItem{
			FL:       do.FL,
			Name:     do.Name,
			Desc:     do.Desc,
			CoverId:  do.CoverId,
			RepoType: do.RepoType,
			Tags:     do.Tags,
			TagKinds: do.TagKinds,
		},
	}

	docObj.ProjectPropertyItem.setDefault()

	return genDoc(docObj)
}

func (col project) toProjectDO(owner string, item *projectItem, do *repositories.ProjectDO) {
	*do = repositories.ProjectDO{
		Id:            item.Id,
		Owner:         owner,
		Name:          item.Name,
		Desc:          item.Desc,
		Type:          item.Type,
		Level:         item.Level,
		CoverId:       item.CoverId,
		Protocol:      item.Protocol,
		Training:      item.Training,
		RepoType:      item.RepoType,
		RepoId:        item.RepoId,
		Tags:          item.Tags,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
		Version:       item.Version,
		LikeCount:     item.LikeCount,
		ForkCount:     item.ForkCount,
		DownloadCount: item.DownloadCount,

		RelatedModels:   toResourceIndexDO(item.RelatedModels),
		RelatedDatasets: toResourceIndexDO(item.RelatedDatasets),
	}
}
