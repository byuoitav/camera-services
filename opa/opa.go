package opa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	cameraservices "github.com/byuoitav/camera-services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Client struct {
	Address  string
	Endpoint string
	Token    string
	Disable  bool

	Logger *zap.Logger
}

type response struct {
	DecisionID string          `json:"decision_id"`
	Result     map[string]bool `json:"result"`
}

type request struct {
	Input requestData `json:"input"`
}

type requestData struct {
	User   string `json:"user"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

func (client *Client) AuthorizeFor(keys ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !client.IsAuthorizedFor(c.Request.Context(), keys...) {
			c.String(http.StatusForbidden, "Unauthorized")
			c.Abort()
			return
		}

		c.Next()
	}
}

func (client *Client) IsAuthorizedFor(ctx context.Context, keys ...string) bool {
	if client.Disable {
		return true
	}

	auth := cameraservices.CtxAuth(ctx)
	for _, key := range keys {
		if !auth[key] {
			return false
		}
	}

	return true
}

func (client *Client) FillAuth(c *gin.Context) {
	if client.Disable {
		c.Next()
		return
	}

	var user string
	if v, ok := c.Request.Context().Value("user").(string); ok {
		user = v
	}

	oReq, err := json.Marshal(
		request{
			Input: requestData{
				User:   user,
				Path:   c.Request.URL.Path,
				Method: c.Request.Method,
			},
		},
	)

	if err != nil {
		client.Logger.Error("error trying to marshal request to OPA", zap.Error(err))
		c.String(http.StatusInternalServerError, "error while contacting authorization server")
		c.Abort()
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	client.Logger.Debug("Sending authentication request", zap.ByteString("body", oReq))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", client.Address, client.Endpoint), bytes.NewReader(oReq))
	if err != nil {
		client.Logger.Error("error trying to create request to OPA", zap.Error(err))
		c.String(http.StatusInternalServerError, "error while contacting authorization server")
		c.Abort()
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.Token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		client.Logger.Error("error while making request to OPA", zap.Error(err))
		c.String(http.StatusInternalServerError, "error while contacting authorization server")
		c.Abort()
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		client.Logger.Error("unable to read body from OPA", zap.Error(err))
		c.String(http.StatusInternalServerError, "error while contacting authorization server")
		c.Abort()
		return
	}

	client.Logger.Debug("Authentication response", zap.Int("statusCode", res.StatusCode), zap.ByteString("body", body))

	if res.StatusCode != http.StatusOK {
		client.Logger.Error("bad response from opa", zap.Int("statusCode", res.StatusCode))
		c.String(http.StatusInternalServerError, "error while contacting authorization server")
		c.Abort()
		return
	}

	var oRes response
	if err := json.Unmarshal(body, &oRes); err != nil {
		client.Logger.Error("unable to parse body from OPA: %s", zap.Error(err))
		c.String(http.StatusInternalServerError, "error while contacting authorization server")
		c.Abort()
		return
	}

	c.Request = c.Request.WithContext(cameraservices.WithAuth(c.Request.Context(), oRes.Result))
	c.Next()
}
