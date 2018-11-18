package main

import (
	"GenArt/generative"
	"GenArt/scrape"
	"github.com/azer/go-flickr"
	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
)

func main() {

	client := flickr.Client{
		Key: os.Getenv("FLICKR"),
	}

	// Do jobs with params
	gocron.Every(1).Second().Do(scrapeFlickr, client)
	gocron.Every(1).Second().Do(generateImage)

	// function Start start all the pending jobs
	<-gocron.Start()

	//port := flag.String("p", "8100", "port to serve on")
	//directory := flag.String("d", ".", "the directory of static file to host")
	//flag.Parse()
	//
	//http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	//
	//log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	//log.Fatal(http.ListenAndServe(":"+*port, nil))

}

func scrapeFlickr(client flickr.Client) {
	ids, err := scrape.GetImagesIds(client, "trees")
	if err != nil {
		logrus.Error("unable to get image IDs ", err.Error())
	}
	url, err := scrape.GetImageUrl(client, ids[rand.Intn(len(ids))])
	if err != nil {
		logrus.Error("unable to get image url ", err.Error())
	}
	if url != "" {
		scrape.SaveImage(url)
	}
}

func generateImage() {

	img, err := generative.OpenImage(scrape.StoreImage)
	if err != nil {
		logrus.Error("unable to open image")
	}

	// Crop attempts to find the best crop of img based on the given width and height values.
	img, err = generative.Crop(img, 1000, 1000)
	if err != nil {
		logrus.Error("unable to crop image " + err.Error())
	}
	err = generative.SaveImage(img, ".", "images/cropped.jpg")
	if err != nil {
		logrus.Error("unable to save mage " + err.Error())
	}
	img, err = generative.OpenImage("images/cropped.jpg")
	if err != nil {
		logrus.Error("unable to open cropped image " + err.Error())
	}
	sat := generative.Saturate(img)
	err = generative.SaveImage(sat, ".", "images/saturated.jpg")
	if err != nil {
		logrus.Error("unable to save saturated image " + err.Error())
	}
	mult := generative.Multiply(img)
	err = generative.SaveImage(mult, ".", "images/multiplied.jpg")
	if err != nil {
		logrus.Error("unable to save multiplied image " + err.Error())
	}
	shrp := generative.Sharpen(sat)
	err = generative.SaveImage(shrp, ".", "images/sharpened.jpg")
	if err != nil {
		logrus.Error("unable to save sharpened image " + err.Error())
	}
	pri := generative.PrimitivePicture(sat)
	err = generative.SaveImage(pri, ".", "images/primitive.jpg")
	if err != nil {
		logrus.Error("unable to save primitive image " + err.Error())
	}
}
