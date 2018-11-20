package scrape

import (
	"GenArt/generative"
	"encoding/json"
	"github.com/azer/go-flickr"
	"github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net/http"
	"os"
)

// Get Image ids from specified tag
func getImagesIds(client flickr.Client) ([]string, error) {

	var ids []string

	tag, err := getHotTag(client)
	if err != nil {
		logrus.Error("hot tag request not successful ", err.Error())
		return ids, err
	}

	args := []string{"tags", tag}

	tagResp, err := flickrRequest(client, args, "photos.search")
	if err != nil {
		logrus.Error("flickr request not successful ", err.Error())
		return ids, err
	}

	if !checkStatus(tagResp) {
		logrus.Debugf("status check failed - see message ->", tagResp["stat"])
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
func getImageUrl(client flickr.Client, id string) (string, error) {

	args := []string{"photo_id", id}
	url, err := flickrRequest(client, args, "photos.getSizes")
	if err != nil {
		logrus.Error("flickr request not successful ", err.Error())
		return "", err
	}

	if !checkStatus(url) {
		logrus.Debugf("status check failed - see message ->", url["stat"])
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

func getHotTag(client flickr.Client) (string, error) {
	args := []string{"count", "20"}
	hotList, err := flickrRequest(client, args, "tags.getHotList")
	if err != nil {
		logrus.Error("flickr request not successful ", err.Error())
		return "", err
	}

	if !checkStatus(hotList) {
		logrus.Debugf("status check failed - see message ->", hotList["stat"])
	}

	tags := hotList["hottags"].(map[string]interface{})
	tag := tags["tag"].([]interface{})
	hot := tag[rand.Intn(len(tag))].(map[string]interface{})["_content"].(string)

	return hot, nil

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

func saveImage(url string) {

	response, err := http.Get(url)
	if err != nil {
		logrus.Error("unable to GET image from URL", err.Error())
	}
	defer response.Body.Close()

	os.MkdirAll("images", os.ModePerm)

	file, err := os.Create(generative.StoreImage)
	if err != nil {
		logrus.Error("unable to save image from URL", err.Error())
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		logrus.Error("io.Copy unable to dump the response body to the file", err.Error())
	}

}

func ScrapeFlickr(client flickr.Client) {
	logrus.Info("scraping new image")
	ids, err := getImagesIds(client)
	if err != nil {
		logrus.Error("unable to get image IDs ", err.Error())
	}

	if len(ids) > 0 {

		url, err := getImageUrl(client, ids[rand.Intn(len(ids))])
		if err != nil {
			logrus.Error("unable to get image url ", err.Error())
		}
		if url != "" {
			saveImage(url)
		}

	}

}

func checkStatus(resp map[string]interface{}) bool {

	if resp["stat"] == "ok" {
		return true
	} else {
		return false
	}

}
