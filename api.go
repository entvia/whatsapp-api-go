package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var (
	DefaultURI     = "https://graph.facebook.com"
	DefaultVersion = "v17.0"
)

type API struct {
	Token               string
	WebHookVerification string
	URI                 string `default:"https://graph.facebook.com"`
	Version             string `default:"v17.0"`
	client              *http.Client
}

type ErrorResponse struct {
	Data struct {
		Message      string `json:"message"`
		Type         string `json:"type"`
		Code         int32  `json:"code"`
		ErrorSubCode int32  `json:"error_subcode"`
		FbTraceId    string `json:"fbtrace_id"`
		ErrorData    struct {
			MessagingProduct string `json:"messaging_product"`
			Details          string `json:"details"`
		} `json:"error_data"`
	} `json:"error"`
}

func (e *ErrorResponse) Error() string {
	return e.Data.Message
}

func (api *API) request(endpoint string, method string, params map[string]interface{}, body map[string]interface{}) (result []byte, status int, err error) {
	if api.client == nil {
		api.client = &http.Client{}
	}

	uri := fmt.Sprintf("%s/%s%s", api.URI, api.Version, endpoint)

	if params != nil {
		uri = fmt.Sprintf("%s?", uri)
		query := url.Values{}

		for k, v := range params {
			query.Add(k, v.(string))
		}
		uri = uri + query.Encode()
	}
	var reqBody []byte
	if body != nil {
		reqBody, err = json.Marshal(body)
	}
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(reqBody))
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api.Token))

	resp, err := api.client.Do(req)
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return
	}
	result, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	status = resp.StatusCode
	return
}

// downloadRequest is a helper function to download a media file or a document
// it takes url of the file and uses token to authenticate the request
func (api *API) downloadRequest(url string) (io.ReadCloser, int, error) {
	if api.client == nil {
		api.client = &http.Client{}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api.Token))

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	return resp.Body, resp.StatusCode, nil
}
