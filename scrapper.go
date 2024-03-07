package main

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/gocolly/colly/v2"
)

type PokemonProduct struct{
	url , image,name , price string 
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str{
			return false
		}
	}
	return false
}

var pokemonProducts []PokemonProduct


func main() {
	
	// Instantiate default collector
	c := colly.NewCollector()
	i := 1
	limit := 5

	var pagesToScrape []string
	
	pageToScrape := "https://scrapeme.live/shop/page/1"

	// initializing the list of pages discovered with a pageToScrape 
	pagesDiscovered := []string{ pageToScrape } 

	c.OnHTML("a.page-numbers" , func( e *colly.HTMLElement){
		newPaginationLink := e.Attr("href")

		if !contains(pagesToScrape,newPaginationLink){
			pagesToScrape = append(pagesToScrape,newPaginationLink)
		}

		pagesDiscovered = append(pagesDiscovered, newPaginationLink)

	})
 	

	c.OnHTML( "li.product" , func(e *colly.HTMLElement){
		// Extract data from the HTML elements
		url := e.ChildAttr("a.woocommerce-LoopProduct-link", "href")
		image := e.ChildAttr("img", "src")
		name := e.ChildText("h2.woocommerce-loop-product__title")
		price := e.ChildText("span.price")

		// Append the extracted data to the slice
		pokemonProducts = append(pokemonProducts, PokemonProduct{
			url:   url,
			image: image,
			name:  name,
			price: price,
		})
	})

	c.OnScraped( func(r *colly.Response){
		if len(pagesToScrape) != 0 && i < limit {
			
			pageToScrape = pagesToScrape[0]
			pagesToScrape = pagesToScrape[1:]

			i++

			c.Visit(pageToScrape)
		}
	})

	c.OnRequest( func( r *colly.Request){
		fmt.Println("Visiting : ",r.URL.String())
	})

	// setting a valid User-Agent header 
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
	
	// downloading the target HTML page 
	// visiting the first page 
	c.Visit(pageToScrape) 

	for _, pokemonProduct := range pokemonProducts { 
		// converting a PokemonProduct to an array of strings 
		 
			fmt.Printf("%v , %v , %v , %v \n",  
			pokemonProduct.url, 
			pokemonProduct.image, 
			pokemonProduct.name, 
			pokemonProduct.price)
		
	 
	} 

	//openeing a CSV file
	file , err := os.Create("pokemon.csv")
	if err != nil{
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers:= []string{
		"url","image","name","price",
	}

	err = writer.Write(headers)
	if err != nil{
		panic(err)
	}

	for _, pokemonProduct := range pokemonProducts { 
		// converting a PokemonProduct to an array of strings 
		 
			record := []string{
				pokemonProduct.url, 
				pokemonProduct.image, 
				pokemonProduct.name, 
				pokemonProduct.price,
			} 

		
			err = writer.Write(record)
			if err != nil{
				panic(err)
			}
			
	} 

	defer writer.Flush()
	
}