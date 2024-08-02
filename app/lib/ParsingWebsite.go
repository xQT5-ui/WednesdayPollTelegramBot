package lib

import (
	"log"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func GetDataFromWebsite(url string) string {
	// Get website data
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Ошибка получения информации:\n%v", err)
		return ""
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("Ошибка доступности информации: %d %s", res.StatusCode, res.Status)
		return ""
	}

	// Read website data
	body, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Ошибка чтения информации с сайта:\n%v", err)
		return ""
	}

	fact := body.Find("div#fact").Text()
	if fact != "" {
		pattern := `\.[\w|а-яА-Я]+`

		re := regexp.MustCompile(pattern)
		fact = re.ReplaceAllString(fact, ".")
	}

	log.Printf("Информация с сайта '%s' успешно получена:\n%s", url, fact)

	return fact
}
