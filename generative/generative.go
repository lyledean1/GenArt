package generative

import (
	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/transform"
	"github.com/artyom/smartcrop"
	"github.com/fogleman/primitive/primitive"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"image"
	"image/jpeg"
	"math/rand"
	"os"
	"path"
	"runtime"
	"time"
)

const StoreImage = "images/jpeg.jpg"
const Cropped = "images/cropped.jpg"
const Saturated = "images/saturated.jpg"
const Multiplied = "images/multiplied.jpg"
const Sharpended = "images/sharpened.jpg"
const Primitive = "images/primitive.jpg"

//openImage imports an image from a given path.
func OpenImage(path string) (image.Image, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot open "+path)
	}

	// Decode from JPG into image.Image format.
	img, err := jpeg.Decode(imgFile)
	if err != nil {
		return nil, errors.Wrap(err, "Decoding the image failed.")
	}

	return img, nil
}

// saveImage saves the image to `pname/fname.jpg`.
func SaveImage(img image.Image, pname, fname string) error {
	fpath := path.Join(pname, fname)

	f, err := os.Create(fpath)
	if err != nil {
		return errors.Wrap(err, "Cannot create file: "+fpath)
	}
	err = jpeg.Encode(f, img, &jpeg.Options{Quality: 85})
	if err != nil {
		return errors.Wrap(err, "Failed to encode the image as JPEG")
	}
	return nil
}

// The SubImager interface exposes the SubImage method to facilitate the type conversion from `Image` to the appropriate color type.
type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

// `crop` auto-crops the image in-place.
func Crop(img image.Image, width, height int) (image.Image, error) {

	rect, err := smartcrop.Crop(img, width, height)
	if err != nil {
		return nil, errors.Wrap(err, "Smartcrop failed")
	}

	si, ok := (img).(SubImager)
	if !ok {
		return nil, errors.New("crop(): img does not support SubImage()")
	}
	subImg := si.SubImage(rect)

	return subImg, nil
}

// Apply 50% saturation
func Saturate(img image.Image) image.Image {
	return adjust.Saturation(img, 0.5)
}

// Multiply the image with itself
func Multiply(img image.Image) image.Image {
	return blend.Multiply(img, img)
}

// Sharpen the image using unsharp masking.
func Sharpen(img image.Image) image.Image {
	return effect.UnsharpMask(img, 0.6, 1.2)
}

//Making art.
func PrimitivePicture(img image.Image) image.Image {

	// Resize the image to 256x256 to save processing time.
	// `transform` is a `bild` package.

	img = transform.Resize(img, 256, 256, transform.Linear)

	// Seed random number generator.
	rand.Seed(time.Now().UTC().UnixNano())

	// Set the background color.
	bg := primitive.MakeColor(primitive.AverageImageColor(img))

	// NewModel(image, background color, output size, # of workers)
	model := primitive.NewModel(img, bg, 1024, runtime.NumCPU())

	for i := 0; i < 100; i++ {
		// 5 = rotated rectangles,
		// 128 = default alpha,
		// 0 = default repeat
		model.Step(primitive.ShapeType(5), 128, 0)
	}

	return model.Context.Image()
}

func GenerateImage() {

	logrus.Info("generating new image")
	img, err := OpenImage(StoreImage)
	if err != nil {
		logrus.Error("unable to open image")
	}

	// Crop attempts to find the best crop of img based on the given width and height values.
	img, err = Crop(img, 1000, 1000)
	if err != nil {
		logrus.Error("unable to crop image " + err.Error())
	}
	err = SaveImage(img, ".", Cropped)
	if err != nil {
		logrus.Error("unable to save mage " + err.Error())
	}
	img, err = OpenImage(Cropped)
	if err != nil {
		logrus.Error("unable to open cropped image " + err.Error())
	}
	sat := Saturate(img)
	err = SaveImage(sat, ".", Saturated)
	if err != nil {
		logrus.Error("unable to save saturated image " + err.Error())
	}
	mult := Multiply(img)
	err = SaveImage(mult, ".", Multiplied)
	if err != nil {
		logrus.Error("unable to save multiplied image " + err.Error())
	}
	shrp := Sharpen(sat)
	err = SaveImage(shrp, ".", Sharpended)
	if err != nil {
		logrus.Error("unable to save sharpened image " + err.Error())
	}
	pri := PrimitivePicture(sat)
	err = SaveImage(pri, ".", Primitive)
	if err != nil {
		logrus.Error("unable to save primitive image " + err.Error())
	}
}
