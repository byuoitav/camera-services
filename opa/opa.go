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
	"github.com/byuoitav/camera-services/cmd/control/middleware"
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
	APIKey string `json:"api_key"`
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

// New OPA stuff
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

func (client *Client) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Initial data
		opaData := opaRequest{
			Input: requestData{
				Path:   c.Request.URL.Path,
				Method: c.Request.Method,
			},
		}
		fmt.Printf("Context Output: %v\n", c.Request.Context().Value("user"))
		fmt.Printf("Context Output: %v\n", c.Request.Context().Value("userBYUID"))

		// use either the user netid for the authorization request or an
		// API key if one was used instead
		if user, ok := c.Request.Context().Value("user").(string); ok {
			opaData.Input.User = user
			fmt.Printf("User Found\n")
		} else if apiKey, ok := middleware.GetAVAPIKey(c.Request.Context()); ok {
			opaData.Input.APIKey = apiKey
		}

		// Prep the request
		oReq, err := json.Marshal(opaData)
		if err != nil {
			fmt.Printf("Error trying to create request to OPA: %s\n", err)
			c.String(http.StatusInternalServerError, "Error while contacting authorization server")
			c.Abort()
			return
		}

		req, err := http.NewRequest(
			"POST",
			fmt.Sprintf("%s/v1/data/viaalert", client.Address),
			bytes.NewReader(oReq),
		)

		req.Header.Set("authorization", fmt.Sprintf("Bearer %s", client.Token))

		fmt.Printf("Data: %s\n", opaData)
		fmt.Printf("URL: %s\n", client.Address)

		// Make the request
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("Error while making request to OPA: %s", err)
			c.String(http.StatusInternalServerError, "Error while contacting authorization server")
			c.Abort()
			return
		}
		if res.StatusCode != http.StatusOK {
			fmt.Printf("Got back non 200 status from OPA: %d", res.StatusCode)
			c.String(http.StatusInternalServerError, "Error while contacting authorization server")
			c.Abort()
			return
		}

		// Read the body
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Printf("Unable to read body from OPA: %s", err)
			c.String(http.StatusInternalServerError, "Error while contacting authorization server")
			c.Abort()
			return
		}

		// Unmarshal the body
		oRes := opaResponse{}
		err = json.Unmarshal(body, &oRes)
		if err != nil {
			fmt.Printf("Unable to parse body from OPA: %s", err)
			c.String(http.StatusInternalServerError, "Error while contacting authorization server")
			c.Abort()
			return
		}
		fmt.Printf("Results: %v\n", oRes.Result)
		// If OPA approved then allow the request, else reject with a 403
		if oRes.Result.Allow {
			c.Next()
		} else {
			fmt.Printf("Unauthorized\n")
			c.String(http.StatusForbidden, "Unauthorized")
			c.Abort()
		}
	}
}
