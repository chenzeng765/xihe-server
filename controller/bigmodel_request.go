package controller

import (
	"errors"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

type describePictureResp struct {
	Desc string `json:"desc"`
}

type pictureGenerateRequest struct {
	Desc string `json:"desc"`
}

func (req *pictureGenerateRequest) validate() error {
	if req.Desc == "" {
		return errors.New("missing desc")
	}

	return nil
}

type pictureGenerateResp struct {
	Picture string `json:"picture"`
}

type multiplePicturesGenerateResp struct {
	Pictures []string `json:"pictures"`
}

type questionAskRequest struct {
	Question string `json:"question"`
	Picture  string `json:"picture"`
}

func (req *questionAskRequest) toCmd() (
	q domain.Question, p string, err error,
) {
	if q, err = domain.NewQuestion(req.Question); err != nil {
		return
	}

	if req.Picture == "" {
		err = errors.New("missing picture")

		// TODO doesnot allow chinese name for picture

		return
	}

	p = req.Picture

	return
}

type questionAskResp struct {
	Answer string `json:"answer"`
}

type pictureUploadResp struct {
	Path string `json:"path"`
}

type panguRequest struct {
	Question string `json:"question"`
}

type panguResp struct {
	Answer string `json:"answer"`
}

type luojiaResp struct {
	Answer string `json:"answer"`
}

type CodeGeexRequest struct {
	Lang    string `json:"lang"`
	Content string `json:"content"`
}

func (req *CodeGeexRequest) toCmd() (
	cmd app.CodeGeexCmd, err error,
) {
	cmd.Lang = req.Lang
	cmd.Content = req.Content

	err = cmd.Validate()

	return
}

type wukongRequest struct {
	Desc  string `json:"desc"`
	Style string `json:"style"`
}

func (req *wukongRequest) toCmd() (cmd app.WuKongCmd, err error) {
	cmd.Style = req.Style

	if cmd.Desc, err = domain.NewWuKongPictureDesc(req.Desc); err != nil {
		return
	}

	err = cmd.Validate()

	return
}

type wukongPicturesGenerateResp struct {
	Pictures map[string]string `json:"pictures"`
}

type wukongAddLikeFromTempRequest struct {
	OBSPath string `json:"obspath"`
}

func (req *wukongAddLikeFromTempRequest) toCmd(user domain.Account) app.WuKongAddLikeFromTempCmd {
	return app.WuKongAddLikeFromTempCmd{
		User:    user,
		OBSPath: req.OBSPath,
	}
}

type wukongAddLikeFromPublicRequest struct {
	Id string `json:"id"`
}

func (req *wukongAddLikeFromPublicRequest) toCmd(user domain.Account) app.WuKongAddLikeFromPublicCmd {
	return app.WuKongAddLikeFromPublicCmd{
		User: user,
		Id:   req.Id,
	}
}

type wukongAddLikeResp struct {
	Id string `json:"id"`
}

type wukongPictureLink struct {
	Link string `json:"link"`
}
