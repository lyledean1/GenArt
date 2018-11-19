package main

import (
	"GenArt/generative"
	"GenArt/scrape"
	"flag"
	"github.com/azer/go-flickr"
	"github.com/jasonlvhit/gocron"
	"log"
	"net/http"
	"os"
)

func main() {

	client := flickr.Client{
		Key: os.Getenv("FLICKR"),
	}

	go cron(client)

	port := flag.String("p", "8100", "port to serve on")
	directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))

}

func cron(client flickr.Client) {
	gocron.Every(1).Minute().Do(scrape.ScrapeFlickr, client, "trees")
	gocron.Every(2).Minutes().Do(generative.GenerateImage)
	<-gocron.Start()
}
