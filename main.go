package main

import (
	"bytes"
	"embed"
	"image"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"

	"github.com/artyom/smartcrop"
	"github.com/edwvee/exiffix"
	"github.com/gofiber/fiber/v2"
)

//go:embed index.html

var content embed.FS

var imageDirectory = os.Getenv("IMAGE_DIRECTORY")

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		body, _ := content.ReadFile("index.html")
		c.Type("html")
		return c.Send(body)
	})

	app.Get("/image", func(c *fiber.Ctx) error {
		photos, _ := filepath.Glob(filepath.Join(imageDirectory, "*.jpg"))
		width, _ := strconv.Atoi(c.Query("width", "1280"))
		height, _ := strconv.Atoi(c.Query("height", "800"))

		img := cropImage(photos[rand.Intn(len(photos))], width, height)

		buf := new(bytes.Buffer)
		jpeg.Encode(buf, img, &jpeg.Options{Quality: 100})

		c.Type("jpg")
		return c.SendStream(buf, buf.Len())
	})

	app.Listen(":3000")
}

func cropImage(imgSrc string, width, height int) image.Image {
	f, err := os.Open(imgSrc)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	img, _, err := exiffix.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	topCrop, err := smartcrop.Crop(img, width, height)
	if err != nil {
		log.Fatal(err)
	}

	type subImager interface {
		SubImage(image.Rectangle) image.Image
	}
	if si, ok := img.(subImager); ok {
		return si.SubImage(topCrop)
	}
	return nil
}
