package main

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
)

// initializing a data structure to keep the scraped data
type Product struct {
	url, image, name, price string
}

func main() {
	// initializing the slice of structs to store the data to scrape
	var products []Product

	// creating a new Colly instance
	c := colly.NewCollector()

	// scraping logic
	c.OnHTML("li.product", func(e *colly.HTMLElement) {
		product := Product{}

		product.url = e.ChildAttr("a", "href")
		product.image = e.ChildAttr("img", "src")
		product.name = e.ChildText("h2")
		product.price = e.ChildText(".price")

		products = append(products, product)
	})

	c.OnScraped(func(r *colly.Response) {

		// opening the CSV file
		file, err := os.Create("products.csv")
		if err != nil {
			log.Fatalln("Failed to create output CSV file", err)
		}
		defer file.Close()

		// initializing a file writer
		writer := csv.NewWriter(file)

		// writing the CSV headers
		headers := []string{
			"url",
			"image",
			"name",
			"price",
		}
		writer.Write(headers)

		// writing each product as a CSV row
		for _, product := range products {
			// converting a Product to an array of strings
			record := []string{
				product.url,
				product.image,
				product.name,
				product.price,
			}

			// adding a CSV record to the output file
			writer.Write(record)
		}
		defer writer.Flush()
	})

	// downloading the target HTML page
	c.Visit("https://www.scrapingcourse.com/ecommerce/")
}
