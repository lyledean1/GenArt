package main

import (
	"GenArt/generative"
	"github.com/sirupsen/logrus"
	"log"
)

func main() {
	img, err := generative.OpenImage("images/jpeg.jpg")
	if err != nil {
		log.Fatal(err)
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
