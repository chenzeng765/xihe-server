package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ProjectCreateCmd struct {
	Owner    domain.Account
	Name     domain.ProjName
	Desc     domain.ProjDesc
	Type     domain.ProjType
	CoverId  domain.CoverId
	RepoType domain.RepoType
	Protocol domain.ProtocolName
	Training domain.TrainingPlatform
}

func (cmd *ProjectCreateCmd) Validate() error {
	b := cmd.Owner != nil &&
		cmd.Name != nil &&
		cmd.Type != nil &&
		cmd.CoverId != nil &&
		cmd.RepoType != nil &&
		cmd.Protocol != nil &&
		cmd.Training != nil

	if !b {
		return errors.New("invalid cmd of creating project")
	}

	return nil
}

func (cmd *ProjectCreateCmd) toProject() domain.Project {
	return domain.Project{
		Owner:    cmd.Owner,
		Type:     cmd.Type,
		Protocol: cmd.Protocol,
		Training: cmd.Training,
		ProjectModifiableProperty: domain.ProjectModifiableProperty{
			Name:     cmd.Name,
			Desc:     cmd.Desc,
			CoverId:  cmd.CoverId,
			RepoType: cmd.RepoType,
		},
	}
}

type ProjectDTO struct {
	Id       string   `json:"id"`
	Owner    string   `json:"owner"`
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	Type     string   `json:"type"`
	CoverId  string   `json:"cover_id"`
	Protocol string   `json:"protocol"`
	Training string   `json:"training"`
	RepoType string   `json:"repo_type"`
	RepoId   string   `json:"repo_id"`
	Tags     []string `json:"tags"`
}

type ProjectService interface {
	Create(*ProjectCreateCmd) (ProjectDTO, error)
	GetByName(domain.Account, domain.ProjName) (ProjectDTO, error)
	List(domain.Account, *ResourceListCmd) ([]ProjectDTO, error)
	Update(*domain.Project, *ProjectUpdateCmd) (ProjectDTO, error)
	Fork(*ProjectForkCmd) (ProjectDTO, error)

	AddLike(domain.Account, string) error
	RemoveLike(domain.Account, string) error
}

func NewProjectService(repo repository.Project, pr platform.Repository) ProjectService {
	return projectService{repo: repo, pr: pr}
}

type projectService struct {
	repo repository.Project
	pr   platform.Repository
}

func (s projectService) Create(cmd *ProjectCreateCmd) (dto ProjectDTO, err error) {
	pid, err := s.pr.New(platform.RepoOption{
		Name:     cmd.Name,
		Desc:     cmd.Desc,
		RepoType: cmd.RepoType,
	})
	if err != nil {
		return
	}

	v := cmd.toProject()
	v.RepoId = pid

	p, err := s.repo.Save(&v)
	if err != nil {
		return
	}

	s.toProjectDTO(&p, &dto)

	return
}

func (s projectService) GetByName(
	owner domain.Account, name domain.ProjName,
) (dto ProjectDTO, err error) {
	v, err := s.repo.GetByName(owner, name)
	if err != nil {
		return
	}

	s.toProjectDTO(&v, &dto)

	return
}

type ResourceListCmd struct {
	Name     string
	RepoType domain.RepoType
}

func (cmd *ResourceListCmd) toResourceListOption() (
	option repository.ResourceListOption,
) {
	option.Name = cmd.Name
	option.RepoType = cmd.RepoType

	return
}

func (s projectService) List(owner domain.Account, cmd *ResourceListCmd) (
	dtos []ProjectDTO, err error,
) {
	v, err := s.repo.List(owner, cmd.toResourceListOption())
	if err != nil || len(v) == 0 {
		return
	}

	dtos = make([]ProjectDTO, len(v))
	for i := range v {
		s.toProjectDTO(&v[i], &dtos[i])
	}

	return
}

func (s projectService) toProjectDTO(p *domain.Project, dto *ProjectDTO) {
	*dto = ProjectDTO{
		Id:       p.Id,
		Owner:    p.Owner.Account(),
		Name:     p.Name.ProjName(),
		Type:     p.Type.ProjType(),
		CoverId:  p.CoverId.CoverId(),
		Protocol: p.Protocol.ProtocolName(),
		Training: p.Training.TrainingPlatform(),
		RepoType: p.RepoType.RepoType(),
		RepoId:   p.RepoId,
		Tags:     p.Tags,
	}

	if p.Desc != nil {
		dto.Desc = p.Desc.ProjDesc()
	}
}
