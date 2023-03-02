package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
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

	infoCollector := c.Clone()

	c.Limit(&colly.LimitRule{
		DomainGlob: "*cococart.*",
		// Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	c.SetRequestTimeout(120 * time.Second)
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL)
	})

	infoCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting product URL:", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	// var items []item
	var itemsOuter []itemOuter
	var itemsInner []itemInner
	sno := 0

	// c.OnHTML("div.grid-item__content", func(h *colly.HTMLElement) {

	// 	// fmt.Println(h.ChildAttr("a.grid-item__link", "href"))
	// 	sno += 1

	// 	// item := item{
	// 	// 	Sno:          sno,
	// 	// 	Name:         h.ChildText("div.grid-product__title"),
	// 	// 	Image_url:    "x",
	// 	// 	Description:  "x",
	// 	// 	Origin:       "x",
	// 	// 	OfferPrice:   h.ChildText("div.grid-product__price span.grid-product__price--current span.visually-hidden"),
	// 	// 	Price:        h.ChildText("div.grid-product__price span.grid-product__price--original span.visually-hidden"),
	// 	// 	Review_count: h.ChildAttr("div.jdgm-widget div.jdgm-prev-badge", "data-number-of-reviews"),
	// 	// 	Rating:       h.ChildAttr("div.jdgm-widget div.jdgm-prev-badge", "data-average-rating"),
	// 	// 	Id:           "x",
	// 	// 	Weight:       "x",
	// 	// 	Item_link:    h.ChildAttr("a.grid-item__link", "href"),
	// 	// }

	// 	// items = append(items, item)
	// })

	// jumping inside the open site

	c.OnHTML("div.grid-item__content", func(h *colly.HTMLElement) {
		itemObj := itemOuter{}
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

		infoCollector.Visit(productUrl)
		itemsOuter = append(itemsOuter, itemObj)
		// items = append(items, itemObj)
	})

	infoCollector.OnHTML("main#MainContent", func(h *colly.HTMLElement) {
		// fmt.Println(h.ChildAttr("div.product-image-main div.image-wrap > img", "data-photoswipe-src")) // Image_url
		// fmt.Println(h.ChildText("div.product-block > div.rte > p:first")) // description
		// fmt.Println(h.ChildText("div.product-block > div.rte > p:nth-child(2)")) // Origin
		// fmt.Println(h.ChildText("fieldset > div.variant-input > label.variant__button-label")) // weight
		itemObj := itemInner{
			Image_url:   h.ChildAttr("div.product-image-main div.image-wrap > img", "data-photoswipe-src"),
			Description: h.ChildText("div.product-block > div.rte > p:first"),
			Origin:      h.ChildText("div.product-block > div.rte > p:nth-child(2)"),
			Weight:      h.ChildText("fieldset > div.variant-input > label.variant__button-label"),
		}

		itemsInner = append(itemsInner, itemObj)

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
	})

	c.OnHTML("[title=Next]", func(h *colly.HTMLElement) {
		nextPage := h.Request.AbsoluteURL(h.Attr("href"))
		fmt.Println(nextPage)
		if nextPage == "https://cococart.in/collections/shop-all?page=100&sort_by=title-ascending" {
			c.Visit(nextPage)
		}
	})

	c.Visit("https://cococart.in/collections/shop-all?sort_by=title-ascending")
	// c.Wait()
	// fmt.Println(items)
	fmt.Println(itemsOuter)
	fmt.Println(itemsInner)

	// content, err := json.MarshalIndent(items, "", "\t")

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// os.WriteFile("products.json", content, 0644)
	// fmt.Println("output", string(content))

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
