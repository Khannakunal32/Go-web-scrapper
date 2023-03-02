package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

type item struct {
	Name         string `json:"name"`
	Image_url    string `json:"image_url"`
	Description  string `json:"description"`
	Origin       string `json:"origin"`
	OfferPrice   int64  `json:"offer_price"`
	Price        string `json:"price"`
	Review_count int64  `json:"review_count"`
	Rating       int64  `json:"rating"`
	Id           int64  `json:"id"`
	Weight       int64  `json:"weight"`
	Item_link    string `json:"item_link"`
}

func main() {

	c := colly.NewCollector(
	// colly.AllowedDomains("https://cococart.in/"),
	)
	c.SetRequestTimeout(120 * time.Second)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	// var items []item

	c.OnHTML("a.grid-item__link", func(h *colly.HTMLElement) {
		fmt.Println(h.ChildText("div.grid-product__title"))                                                         // Name
		fmt.Println(h.ChildText("div.grid-product__price span.grid-product__price--original span.visually-hidden")) // Price
		fmt.Println(h.ChildText("div.grid-product__price span.grid-product__price--current span.visually-hidden"))  // OfferPrice
		fmt.Println(h.ChildAttr("div.jdgm-widget div.jdgm-prev-badge", "data-average-rating"))                      // rating
		fmt.Println(h.ChildAttr("div.jdgm-widget div.jdgm-prev-badge", "data-number-of-reviews"))                   // reveiw_count
	})

	c.OnHTML("[title=Next]", func(h *colly.HTMLElement) {
		c.Visit(h.Request.AbsoluteURL(h.Attr("href")))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
	})

	// c.OnHTML("div.col-sm-4 div[itemprop=itemListElement]", func(h *colly.HTMLElement) {
	// 	item := item{
	// 		Name:   h.ChildText("h2.product-title"),
	// 		Price:  h.ChildText("div.sale-price"),
	// 		ImgUrl: h.ChildAttr("img", "src"),
	// 	}

	// 	items = append(items, item)
	// })

	// c.OnHTML("[title=Next]", func(h *colly.HTMLElement) {
	// 	c.Visit(h.Request.AbsoluteURL(h.Attr("href")))
	// })

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println(r.URL.String())
	// })

	// c.Visit("https://j2store.net/demo/index.php/shop")
	// // fmt.Println(items)

	// content, err := json.MarshalIndent(items, "", "\t")

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// os.WriteFile("products.json", content, 0644)
	// fmt.Println("output", string(content))
	c.Visit("https://cococart.in/collections/shop-all?sort_by=title-ascending")
}

// type item struct {
// 	Name   string `json:"name"`
// 	Price  string `json:"price"`
// 	ImgUrl string `json:"imgUrl"`
// }

// func main() {

// 	c := colly.NewCollector(
// 		colly.AllowedDomains("j2store.net"),
// 	)

// 	var items []item

// 	c.OnHTML("div.col-sm-4 div[itemprop=itemListElement]", func(h *colly.HTMLElement) {
// 		item := item{
// 			Name:   h.ChildText("h2.product-title"),
// 			Price:  h.ChildText("div.sale-price"),
// 			ImgUrl: h.ChildAttr("img", "src"),
// 		}

// 		items = append(items, item)
// 	})

// 	c.OnHTML("[title=Next]", func(h *colly.HTMLElement) {
// 		c.Visit(h.Request.AbsoluteURL(h.Attr("href")))
// 	})

// 	c.OnRequest(func(r *colly.Request) {
// 		fmt.Println(r.URL.String())
// 	})

// 	c.Visit("https://j2store.net/demo/index.php/shop")
// 	// fmt.Println(items)

// 	content, err := json.MarshalIndent(items, "", "\t")

// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	os.WriteFile("products.json", content, 0644)
// 	fmt.Println("output", string(content))
// }
