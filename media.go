package whatsapp

import (
	"encoding/json"
	"fmt"
	"io"
)

var (
	MediaAudio    = "audio"
	MediaDocument = "document"
	MediaImage    = "image"
	MediaSticker  = "sticker"
	MediaVideo    = "video"
)

type Media struct {
	File     string `json:"file"`
	IsItAnID bool   `json:"is_it_an_id"`
	Type     string `json:"type"`
	api      *API
}

type MediaId struct {
	Type string `json:"-"`
	Id   string `json:"id"`
	api  *API
}

type MediaLink struct {
	Type string `json:"-"`
	Link string `json:"link"`
	api  *API
}

func (obj *Media) ToId() *MediaId {
	return &MediaId{api: obj.api, Type: obj.Type, Id: obj.File}
}
func (obj *Media) ToLink() *MediaLink {
	return &MediaLink{api: obj.api, Type: obj.Type, Link: obj.File}
}

func (m *MediaLink) ToMedia() *Media {
	return &Media{api: m.api, Type: m.Type, File: m.Link, IsItAnID: false}
}

func (m *MediaId) ToMedia() *Media {
	return &Media{api: m.api, Type: m.Type, File: m.Id, IsItAnID: true}
}

func (api *API) GetMediaData(phoneId string, mediaId string) (*MediaResponse, error) {

	// endpoint := fmt.Sprintf("/%s/media/%s", phoneId, mediaId)
	endpoint := fmt.Sprintf("/%s", mediaId)

	res, status, err := api.request(endpoint, "GET", nil, nil)
	if err != nil {

		return nil, err
	}

	if status != 200 {
		e := ErrorResponse{}
		json.Unmarshal(res, &e)
		return nil, &e
	}

	var response MediaResponse
	err = json.Unmarshal(res, &response)
	return &response, err
}

// DownloadMediaByURL downloads media by URL
// Note: make sure you have `defer body.Close()` in the calling function
func (api *API) DownloadMediaByURL(url string) (io.ReadCloser, error) {

	body, status, err := api.downloadRequest(url)
	if err != nil {
		return nil, err
	}

	// Ensure that the HTTP request was successful
	if status != 200 {
		return nil, fmt.Errorf("failed to download media: HTTP %d, err: %v", status, err)
	}

	return body, nil
}

// DeleteMediaByID deletes media by its ID.
func (api *API) DeleteMediaByID(mediaID string) (bool, error) {
	// Prepare the endpoint with the media ID.
	endpoint := fmt.Sprintf("/%s", mediaID)

	// Call the existing request function.
	_, status, err := api.request(endpoint, "DELETE", nil, nil)
	if err != nil {
		return false, err
	}

	// Ensure that the HTTP request was successful.
	if status != 200 {
		return false, fmt.Errorf("failed to delete media: HTTP %d", status)
	}

	// If the request was successful, return true.
	return true, nil
}

type MediaResponse struct {
	MessagingProduct string `json:"messaging_product"`
	Url              string `json:"url"`
	MimeType         string `json:"mime_type"`
	Sha256           string `json:"sha256"`
	FileSize         int32  `json:"file_size"`
	Id               string `json:"id"`
}
