package whatsapp

import (
	"encoding/json"
	"fmt"
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

	endpoint := fmt.Sprintf("/%s/media/%s", phoneId, mediaId)

	params := map[string]interface{}{}
	params["access_token"] = api.Token

	res, status, err := api.request(endpoint, "GET", params, nil)
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

type MediaResponse struct {
	MessagingProduct string `json:"messaging_product"`
	Url              string `json:"url"`
	MimeType         string `json:"mime_type"`
	Sha256           string `json:"sha256"`
	FileSize         string `json:"file_size"`
	Id               string `json:"id"`
}
