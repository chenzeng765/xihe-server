package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/cloud/domain"
	"github.com/opensourceways/xihe-server/cloud/domain/message"
	"github.com/opensourceways/xihe-server/cloud/domain/repository"
	"github.com/opensourceways/xihe-server/cloud/domain/service"
)

type CloudService interface {
	// cloud
	ListCloud() ([]CloudDTO, error)
	SubscribeCloud(*SubscribeCloudCmd) (code string, err error)

	// pod
	Get(*PodInfoCmd) (PodInfoDTO, error)
	ReleasePod(*RelasePodCmd) (code string, err error)
}

var _ CloudService = (*cloudService)(nil)

func NewCloudService(
	cloudRepo repository.Cloud,
	podRepo repository.Pod,
	producer message.CloudMessageProducer,
) *cloudService {
	return &cloudService{
		cloudRepo:    cloudRepo,
		podRepo:      podRepo,
		producer:     producer,
		cloudService: service.NewCloudService(podRepo, producer),
	}
}

type cloudService struct {
	cloudRepo    repository.Cloud
	podRepo      repository.Pod
	producer     message.CloudMessageProducer
	cloudService service.CloudService
}

func (s *cloudService) ListCloud() (dto []CloudDTO, err error) {
	// list cloud conf
	confs, err := s.cloudRepo.ListCloudConf()
	if err != nil {
		return
	}

	// to cloud
	c := make([]domain.Cloud, len(confs))
	for i := range confs {
		c[i].CloudConf = confs[i]
		if err = s.cloudService.ToCloud(&c[i]); err != nil {
			return
		}
	}

	// to dto
	dto = make([]CloudDTO, len(c))
	for i := range c {
		dto[i].toCloudDTO(&c[i], c[i].HasFree())
	}

	return
}

func (s *cloudService) SubscribeCloud(cmd *SubscribeCloudCmd) (code string, err error) {
	// get cloud conf
	c := new(domain.Cloud)
	c.CloudConf, err = s.cloudRepo.GetCloudConf(cmd.CloudId)
	if err != nil {
		return
	}

	// check
	if err = s.cloudService.ToCloud(c); err != nil {
		return
	}

	if !c.HasFree() {
		code = errorResourceBusy
		err = errors.New("no free resource remain")

		return
	}

	// subscribe
	err = s.cloudService.SubscribeCloud(&c.CloudConf, cmd.User)

	return
}
