package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
	"unicode"

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

func countWhitespace(s string) int {
	count := 0
	for _, char := range s {
		if unicode.IsSpace(char) {
			count++
		}
	}
	return count
}

var(
	 pokemonProducts []PokemonProduct
	 limit int
	//  pagesToScrape  []string

	 lastRequestTime map[string]time.Time
)

var pagesToScrape = []string{ 
	"https://scrapeme.live/shop/page/1/", 
	"https://scrapeme.live/shop/page/2/", 
	"https://scrapeme.live/shop/page/3/", 
	"https://scrapeme.live/shop/page/4/", 
	"https://scrapeme.live/shop/page/5/", 
	"https://scrapeme.live/shop/page/6/", 
	"https://scrapeme.live/shop/page/7/", 
	"https://scrapeme.live/shop/page/8/", 
	"https://scrapeme.live/shop/page/9/", 
	"https://scrapeme.live/shop/page/10/", 
 
}


func main() {
	
	// Instantiate default collector
	c := colly.NewCollector( 
		colly.Async(true),
	)
	// i := 1

	var startUrl string
	var customSelector string
	var autopaginate bool
	var linkSelector string
	var paginationPattern string

	//define CLI flags
	flag.StringVar(&startUrl,"url","","starting web crawler URL")
	flag.StringVar(&customSelector,"selector","","Use to select custom CSS element")
	flag.IntVar(&limit,"limit",0,"maximum pages to crawl")
	
	//auto pagginate , link-selector
	flag.StringVar(&paginationPattern, "pagination", "", "Pagination URL pattern (e.g., https://example.com/page/%d)")
	flag.BoolVar(&autopaginate,"autopaginate",false,"use regexp to autopaginate")
	flag.StringVar(&linkSelector,"linkSelector","","use to add links to pages to Scrape")


	//PArse CLI flags
	flag.Parse()

	// initializing the list of pages discovered with a pageToScrape 
	// pagesDiscovered  := []string{ startUrl } 
	
	// // Print defined flags and their values
	// fmt.Println("Defined flags and their values:")
	// fmt.Printf("startUrl: %v\n", startUrl)
	// fmt.Printf("customSelector: %v\n", customSelector)
	// fmt.Printf("limit: %v\n", limit)
	// fmt.Printf("paginationPattern: %v\n", paginationPattern)
	// fmt.Printf("autopaginate: %v\n", autopaginate)
	// fmt.Printf("linkSelector: %v\n", linkSelector)


	// registering all pages to scrape 
	for _, pageToScrape := range pagesToScrape { 
		c.Visit(pageToScrape) 
	} 

	// if autopaginate && paginationPattern != ""{


	// 	lastRequestTime = make( map[string]time.Time )

	// 	// fmt.Println("Link-Selector : ",linkSelector)
	// 	c.OnHTML(linkSelector , func( e *colly.HTMLElement){
	// 		newPaginationLink := e.Attr("href")

	// 		//MATCHING STRINGS HERE 
	// 		// fmt.Println("Matching Stirngs here baby")
	// 		match , err := regexp.MatchString(paginationPattern,newPaginationLink)
	// 		if err != nil{
	// 			log.Fatal("Error in matching strings")
	// 		}
			
	// 		newPaginationLink = strings.TrimSpace(newPaginationLink)
	// 		paginationPattern = strings.TrimSpace(paginationPattern)
			
	// 		// fmt.Printf("newPaginationLink: %v\n", newPaginationLink)
	// 		// fmt.Printf("Pagination Pattern: %v\n", paginationPattern)
	// 		// fmt.Println("Match : ",match)
		
	// 		if match{
	// 			// fmt.Println("INSIDE PAGES")
	// 			if !contains(pagesToScrape,newPaginationLink){
	// 				// fmt.Println("NEW PAGE FOUND")
	// 				pagesToScrape = append(pagesToScrape,newPaginationLink)

	// 				// for _ , page := range pagesToScrape{
	// 				// 	fmt.Printf("url : %v\n",page)
	// 				// }
	// 			}

	// 			pagesDiscovered = append(pagesDiscovered, newPaginationLink)
	// 		}

	// 	})
	// }

	// c.OnHTML("a.page-numbers" , func( e *colly.HTMLElement){
	// 	newPaginationLink := e.Attr("href")
	// 	// fmt.Println("INSIDE PAGES")
	// 	if !contains(pagesToScrape,newPaginationLink){
	// 		// fmt.Println("NEW PAGE FOUND")
	// 		pagesToScrape = append(pagesToScrape,newPaginationLink)
	// 		// for _ , page := range pagesToScrape{
	// 		// 	fmt.Printf("url : %v\n",page)
	// 		// }
	// 	}
	// 	pagesDiscovered = append(pagesDiscovered, newPaginationLink)
	// })
 	

	c.OnHTML("li.product" , func(e *colly.HTMLElement){
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

	// c.OnScraped( func(r *colly.Response){
	// 	if len(pagesToScrape) != 0 && i < limit {
			
	// 		pageToScrape := pagesToScrape[0]
	// 		pagesToScrape = pagesToScrape[1:]
	// 		// fmt.Println("Page to scrape is : ",pageToScrape)
	// 		i++

	// 		c.Visit(pageToScrape)
	// 	}
	// })

	c.OnRequest( func( r *colly.Request){
		fmt.Println("Visiting : ",r.URL.String())
	})


	
	// downloading the target HTML page 
	// visiting the first page 
	c.Visit(startUrl) 
	
	c.Limit(&colly.LimitRule{
		Parallelism: 4,
		DomainGlob: "*",
		RandomDelay: 5*time.Second,
	})

	// wait for tColly to visit all pages 
	c.Wait() 

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