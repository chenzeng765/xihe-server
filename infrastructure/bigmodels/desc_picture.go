package bigmodels

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

type pictureDescInfo struct {
	endpoint string
}

func newPictureDescInfo(cfg *Config) pictureDescInfo {
	ce := &cfg.Endpoints

	es, _ := ce.parse(ce.DescPicture)

	return pictureDescInfo{
		endpoint: es[0],
	}
}

type descOfPicture struct {
	Result struct {
		Instances struct {
			Image []string `json:"image"`
		} `json:"instances"`
	} `json:"inference_result"`
}

func (s *service) DescribePicture(
	picture io.Reader, name string, length int64,
) (string, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	file, err := writer.CreateFormFile("file", name)
	if err != nil {
		return "", err
	}

	n, err := io.Copy(file, picture)
	if err != nil {
		return "", err
	}
	if n != length {
		return "", errors.New("copy file failed")
	}

	writer.Close()

	req, err := http.NewRequest(
		http.MethodPost, s.pictureDescInfo.endpoint, buf,
	)
	if err != nil {
		return "", err
	}

	t, err := s.token()
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Auth-Token", t)

	desc := new(descOfPicture)

	if _, err = s.hc.ForwardTo(req, desc); err != nil {
		return "", err
	}

	if v := desc.Result.Instances.Image; len(v) > 0 {
		return v[0], nil
	}

	return "", errors.New("no content")
}
