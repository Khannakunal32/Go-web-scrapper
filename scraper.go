package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type item struct {
	Name   string `json:"name"`
	Price  string `json:"price"`
	ImgUrl string `json:"imgUrl"`
}

func main() {

	c := colly.NewCollector(
		colly.AllowedDomains("j2store.net"),
	)

	var items []item

	c.OnHTML("div.col-sm-4 div[itemprop=itemListElement]", func(h *colly.HTMLElement) {
		item := item{
			Name:   h.ChildText("h2.product-title"),
			Price:  h.ChildText("div.sale-price"),
			ImgUrl: h.ChildAttr("img", "src"),
		}

		items = append(items, item)
	})

	c.OnHTML("[title=Next]", func(h *colly.HTMLElement) {
		c.Visit(h.Request.AbsoluteURL(h.Attr("href")))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
	})

	c.Visit("https://j2store.net/demo/index.php/shop")
	// fmt.Println(items)

	content, err := json.MarshalIndent(items, "", "\t")

	if err != nil {
		fmt.Println(err.Error())
	}

	os.WriteFile("products.json", content, 0644)
	fmt.Println("output", string(content))
}
