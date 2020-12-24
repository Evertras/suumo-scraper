package suumo

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	colly "github.com/gocolly/colly/v2"

	"github.com/evertras/suumo-scraper/internal/geocode"
)

func ListWards(prefecture Prefecture) ([]Ward, error) {
	wards := make([]Ward, 0)

	c := colly.NewCollector()

	c.OnHTML("#js-areaSelectForm", func(e *colly.HTMLElement) {
		e.ForEach("li", func(i int, e *colly.HTMLElement) {
			name := e.ChildText("label span:first-of-type")
			code := e.ChildAttr("input[type=\"checkbox\"]", "value")

			if code != "" {
				wards = append(wards, Ward{
					Name: name,
					Code: code,
					Prefecture: prefecture,
				})
			}
		})
	})

	err := c.Visit(fmt.Sprintf("https://suumo.jp/chintai/%s/city/", prefecture.URLPath))

	if err != nil {
		return nil, err
	}

	return wards, nil
}

type Geocoder interface {
	GetNeighborhoodLocation(ctx context.Context, neighborhood string) (*geocode.Location, error)
}

func WardListings(ctx context.Context, geocoder Geocoder, ward Ward) ([]Listing, error) {
	var mu sync.Mutex
	listings := make([]Listing, 0)
	linksCollector := colly.NewCollector()
	completed := 0
	totalRequests := 0

	detailCollector := colly.NewCollector(
		colly.Async(),
	)
	detailCollector.Limit(&colly.LimitRule{
		DomainGlob:  "*suumo.jp*",
		Parallelism: 5,
		Delay:       time.Millisecond * 50,
	})

	detailCollector.OnRequest(func(r *colly.Request) {
		mu.Lock()
		totalRequests++
		mu.Unlock()
	})

	detailCollector.OnResponse(func(r *colly.Response) {
		mu.Lock()
		completed++
		log.Printf("    Listings > %d / %d", completed, totalRequests)
		mu.Unlock()
	})

	linksCollector.OnRequest(func(r *colly.Request) {
		log.Printf("%s - Visiting page %s - %s", ward.Name, r.URL.Query().Get("page"), r.URL)
	})

	linksCollector.OnHTML(".pagination", func(e *colly.HTMLElement) {
		e.ForEach("p.pagination-parts:last-of-type a", func(i int, a *colly.HTMLElement) {
			if a.Text == "次へ" {
				a.Request.Visit(a.Attr("href"))
			}
		})
	})

	linksCollector.OnHTML(
		"tr.js-cassette_link td:last-of-type a",
		func(a *colly.HTMLElement) {
			//link := a.Attr("href")

			//detailCollector.Visit("https://suumo.jp" + link)
		},
	)

	linksCollector.OnHTML(".l-cassetteitem > li", func(li *colly.HTMLElement) {
		name := li.ChildText(".cassetteitem_content-title")
		neighborhood := li.ChildText(".cassetteitem_detail-col1")
		yearsRaw := li.ChildText(".cassetteitem_detail-col3 > div:first-of-type")

		parsedYears, err := extractAgeYears(yearsRaw)

		if err != nil {
			panic(fmt.Sprint("bad year:", yearsRaw, err))
		}

		li.ForEach("tr.js-cassette_link", func(i int, tr *colly.HTMLElement) {
			rawFloor := tr.ChildText("td:nth-child(3)")

			parsedFloor, err := extractFloor(rawFloor)

			if err != nil {
				panic(fmt.Sprint("bad floor:", rawFloor, err))
			}

			rawPrice := tr.ChildText("td:nth-child(4) .cassetteitem_price--rent")

			parsedPrice, err := extractPriceYen(rawPrice)

			if err != nil {
				panic(fmt.Sprintf("bad price: %q %v", rawPrice, err))
			}

			layout := tr.ChildText(".cassetteitem_madori")

			rawSquareMeters := tr.ChildText(".cassetteitem_menseki")

			parsedSquareMeters, err := extractSquareMeters(rawSquareMeters)

			if err != nil {
				panic(fmt.Sprintf("bad square meters: %q %v", rawSquareMeters, err))
			}

			loc, err := geocoder.GetNeighborhoodLocation(ctx, neighborhood)

			if err != nil {
				log.Fatalf("Failed to get location for neighborhood %q: %v", neighborhood, err)
			}

			if loc == nil {
				log.Fatal(loc)
			}

			listing := Listing{
				Title:            name,
				Neighborhood:     neighborhood,
				AgeYears:         parsedYears,
				Floor:            parsedFloor,
				PricePerMonthYen: parsedPrice,
				Layout:           layout,
				SquareMeters:     parsedSquareMeters,
				Ward:             ward,
				Location:         loc,
			}

			mu.Lock()
			listings = append(listings, listing)
			mu.Unlock()
		})
	})

	url := fmt.Sprintf(
		"https://suumo.jp/jj/chintai/ichiran/FR301FC001/?ar=030&bs=040&ta=%s&sc=%s&cb=0.0&ct=9999999&mb=0&mt=9999999&et=9999999&cn=9999999&shkr1=03&shkr2=03&shkr3=03&shkr4=03&sngz=&po1=25&pc=10&page=1",
		ward.Prefecture.Code,
		ward.Code,
	)

	err := linksCollector.Visit(url)

	if err != nil {
		return nil, err
	}

	detailCollector.Wait()

	return listings, nil
}
