package service

import (
	"github.com/opensourceways/xihe-server/cloud/domain"
	"github.com/opensourceways/xihe-server/cloud/domain/message"
	"github.com/opensourceways/xihe-server/cloud/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
)

type CloudService struct {
	podRepo repository.Pod
	sender  message.CloudMessageProducer
}

func NewCloudService(
	pod repository.Pod,
	sender message.CloudMessageProducer,
) CloudService {
	return CloudService{
		pod,
		sender,
	}
}

func (r *CloudService) caculateRemain(
	c *domain.Cloud, p *repository.PodInfoList,
) (err error) {
	remain := c.CloudConf.Limited.CloudLimited() - len(p.PodInfos)
	if remain < 0 {
		remain = 0
	}

	if c.Remain, err = domain.NewCloudRemain(remain); err != nil {
		return
	}

	return
}

func (r *CloudService) ToCloud(c *domain.Cloud) (err error) {
	plist, err := r.podRepo.GetRunningPod(c.CloudConf.Id)
	if err != nil {
		return
	}

	r.caculateRemain(c, &plist)

	return
}

func (r *CloudService) SubscribeCloud(
	c *domain.CloudConf, u types.Account,
) (err error) {
	// save into repo
	p := new(domain.PodInfo)
	if err := p.SetStartingPodInfo(c.Id, u); err != nil {
		return err
	}

	var pid string
	if pid, err = r.podRepo.AddStartingPod(p); err != nil {
		return
	}

	// send msg to call pod instance api
	msg := new(message.MsgCloudConf)
	msg.ToMsgCloudConf(c, u, pid)

	return r.sender.SubscribeCloud(msg)
}

func (r *CloudService) ReleasePod(p *domain.Pod) error {
	msg := new(message.MsgPod)
	msg.ToMsgPod(p)

	return r.sender.ReleasePod(msg)
}
