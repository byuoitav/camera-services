package cameraservices

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Service struct {
	DisableAuth bool
	opaAddress  string
	opaToken    string
	log         *zap.Logger
}

type opaResponse struct {
	DecisionID string    `json:"decision_id"`
	Result     opaResult `json:"result"`
}

type opaResult struct {
	Allow bool `json:"allow"`
}

type opaRequest struct {
	Input requestData `json:"input"`
}

type requestData struct {
	User   string `json:"user"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

func New(d bool, logger *zap.Logger, addr, token string) Service {
	return Service{
		DisableAuth: d,
		opaAddress:  addr,
		opaToken:    token,
		log:         logger,
	}
}

func (s *Service) Authorize(c *gin.Context) {
	oReq, err := json.Marshal(
		opaRequest{
			Input: requestData{
				User:   "",
				Path:   c.FullPath(),
				Method: c.Request.Method,
			},
		},
	)
	if err != nil {
		s.log.Error("Error trying to create request to OPA: %s\n", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error while contacting authorization server")
		return
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/v1/data/cameras", s.opaAddress),
		bytes.NewReader(oReq),
	)
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", s.opaToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		s.log.Error("Error while making request to OPA: %s", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error while contacting authorization server")
		return
	}
	if res.StatusCode != http.StatusOK {
		s.log.Error("Unable to read body from OPA: %s", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error while contacting authorization server")
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		s.log.Error("Unable to read body from OPA: %s", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error while contacting authorization server")
		return
	}

	oRes := opaResponse{}
	err = json.Unmarshal(body, &oRes)
	if err != nil {
		s.log.Error("Unable to parse body from OPA: %s", zap.Error(err))
		c.String(http.StatusInternalServerError, "Error while contacting authorization server")
		return
	}

	if oRes.Result.Allow {
		c.String(http.StatusOK, "Authorized")
		return
	} else {
		c.String(http.StatusForbidden, "Unauthorized")
	}
}
