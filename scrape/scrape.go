package scrape

import (
	"encoding/json"
	"github.com/azer/go-flickr"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

const StoreImage = "images/jpeg.jpg"

// Get Image ids from specified tag
func GetImagesIds(client flickr.Client, tag string) ([]string, error) {

	var ids []string

	args := []string{"tags", tag}

	tagResp, err := flickrRequest(client, args, "photos.search")
	if err != nil {
		logrus.Error("flickr request not successful ", err.Error())
		return ids, err
	}

	photos := tagResp["photos"].(map[string]interface{})

	photo := photos["photo"].([]interface{})

	for _, image := range photo {

		id := image.(map[string]interface{})["id"].(string)

		ids = append(ids, id)
	}

	return ids, nil
}

// Get image url from photo id
func GetImageUrl(client flickr.Client, id string) (string, error) {

	args := []string{"photo_id", id}
	url, err := flickrRequest(client, args, "photos.getSizes")
	if err != nil {
		logrus.Error("flickr request not successful ", err.Error())
		return "", err
	}

	urlResp := url["sizes"].(map[string]interface{})

	urls := urlResp["size"].([]interface{})

	for _, image := range urls {

		label := image.(map[string]interface{})["label"].(string)

		if label == "Large" {
			return image.(map[string]interface{})["source"].(string), nil
		}

	}

	return "", nil

}

func flickrRequest(client flickr.Client, args []string, method string) (map[string]interface{}, error) {

	generic := make(map[string]interface{})

	resp, err := client.Request(method, flickr.Params{
		args[0]: args[1],
	})
	if err != nil {
		logrus.Error("flickr request not successful ", err.Error())
		return generic, err
	}

	err = json.Unmarshal(resp, &generic)
	if err != nil {
		logrus.Error("unable to unmarshal flickr response", err.Error())
		return generic, err
	}
	return generic, nil
}

func SaveImage(url string) {

	response, err := http.Get(url)
	if err != nil {
		logrus.Error("unable to GET image from URL", err.Error())
	}
	defer response.Body.Close()

	file, err := os.Create(StoreImage)
	if err != nil {
		logrus.Error("unable to save image from URL", err.Error())
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		logrus.Error("io.Copy unable to dump the response body to the file", err.Error())
	}

}
