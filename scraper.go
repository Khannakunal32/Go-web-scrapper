package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

type item struct {
	Sno          int    `json:"sno"`
	Name         string `json:"name"`
	Image_url    string `json:"image_url"`
	Description  string `json:"description"`
	Origin       string `json:"origin"`
	OfferPrice   string `json:"offer_price"`
	Price        string `json:"price"`
	Review_count string `json:"review_count"`
	Rating       string `json:"rating"`
	Weight       string `json:"weight"`
	Item_link    string `json:"item_link"`
}
type itemOuter struct {
	Sno  int    `json:"sno"`
	Name string `json:"name"`
	// Image_url    string `json:"image_url"`
	// Description  string `json:"description"`
	// Origin       string `json:"origin"`
	OfferPrice   string `json:"offer_price"`
	Price        string `json:"price"`
	Review_count string `json:"review_count"`
	Rating       string `json:"rating"`
	// Weight       string `json:"weight"`
	Item_link string `json:"item_link"`
}

type itemInner struct {
	Image_url   string `json:"image_url"`
	Description string `json:"description"`
	Origin      string `json:"origin"`
	Weight      string `json:"weight"`
}

func main() {

	c := colly.NewCollector(
	// colly.AllowedDomains("https://cococart.in/"),
	// colly.Async(true), // turn on asynchronous request
	// colly.Debugger(&debug.LogDebugger{}), // Attach debugger to collector
	)

	// 	proxySwitcher, err := proxy.RoundRobinProxySwitcher("socks5://188.226.141.127:1080", "socks5://67.205.132.241:1080")
	// if err != nil {
	//   log.Fatal(err)
	// }
	// c.SetProxyFunc(proxySwitcher)

	infoCollector := c.Clone()

	c.Limit(&colly.LimitRule{
		DomainGlob: "*cococart.*",
		// Parallelism: 2,
		RandomDelay: 3 * time.Second,
	})

	// c.SetRequestTimeout(120 * time.Second)
	c.OnRequest(func(r *colly.Request) {
		extensions.RandomUserAgent(c)
		fmt.Println("visiting", r.URL)
	})

	infoCollector.OnRequest(func(r *colly.Request) {
		extensions.RandomUserAgent(c)
		fmt.Println("visiting product URL:", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	// var itemsOuter []itemOuter
	// var itemsInner []itemInner

	var items []item
	itemObj := item{}

	sno := 0

	// jumping inside the open site

	c.OnHTML("div.grid-item__content", func(h *colly.HTMLElement) {
		// itemObj := itemOuter{}

		sno += 1

		itemObj.Sno = sno
		itemObj.Name = h.ChildText("div.grid-product__title")
		// itemObj.Image_url = "x"
		// itemObj.Description = "x"
		// itemObj.Origin = "x"
		// itemObj.Image_url = "x"
		itemObj.OfferPrice = h.ChildText("div.grid-product__price span.grid-product__price--current span.visually-hidden")
		itemObj.Price = h.ChildText("div.grid-product__price span.grid-product__price--original span.visually-hidden")
		itemObj.Review_count = h.ChildAttr("div.jdgm-widget div.jdgm-prev-badge", "data-number-of-reviews")
		itemObj.Rating = h.ChildAttr("div.jdgm-widget div.jdgm-prev-badge", "data-average-rating")
		// itemObj.Weight = "x"
		itemObj.Item_link = h.ChildAttr("a.grid-item__link", "href")

		productUrl := h.ChildAttr("a.grid-item__link", "href")
		productUrl = h.Request.AbsoluteURL(productUrl)

		// if productUrl == "https://cococart.in/collections/shop-all/products/after-eight-mint-chocolate-thins" {
		infoCollector.Visit(productUrl)
		// }
		items = append(items, itemObj)
		itemObj = item{}
	})

	infoCollector.OnHTML("main#MainContent", func(h *colly.HTMLElement) {

		itemObj.Image_url = h.ChildAttr("div.product-image-main div.image-wrap > img", "data-photoswipe-src")
		// itemObj.Description = h.ChildText("div.product-block > div.rte > p:first")
		itemObj.Description = h.ChildText("div.product-block > div.rte > p:nth-child(1)")
		itemObj.Origin = h.ChildText("div.product-block > div.rte > p:nth-child(2)")
		itemObj.Weight = h.ChildText("fieldset > div.variant-input > label.variant__button-label")

	})

	c.OnRequest(func(r *colly.Request) {
		extensions.RandomUserAgent(c)
		fmt.Println(r.URL.String())
	})

	c.OnHTML("[title=Next]", func(h *colly.HTMLElement) {
		nextPage := h.Request.AbsoluteURL(h.Attr("href"))
		fmt.Println(nextPage)
		// if nextPage == "https://cococart.in/collections/shop-all?page=100&sort_by=title-ascending" {
		c.Visit(nextPage)
		// }
	})

	c.Visit("https://cococart.in/collections/shop-all?sort_by=title-ascending")
	// c.Wait()
	// fmt.Println(items)
	// fmt.Println(itemsOuter)
	// fmt.Println(itemsInner)

	content, err := json.MarshalIndent(items, "", "\t")

	if err != nil {
		fmt.Println(err.Error())
	}

	os.WriteFile("chocolate2.json", content, 0644)
	fmt.Println("output", string(content))

}

// type item struct {
// 	Name   string `json:"name"
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
