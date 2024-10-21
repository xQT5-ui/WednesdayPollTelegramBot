package lib

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"

	lg "app.go/app/lib/logger"
)

func DataFromWebsite(url string, log *lg.Logger) string {
	// Get website data
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err, "Ошибка подключения к сайту:")
		return ""
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatal(fmt.Errorf("ошибка доступности информации: %d %s", res.StatusCode, res.Status), "")
		return ""
	}

	// Read website data
	body, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err, "Ошибка чтения информации с сайта:")
		return ""
	}

	fact := body.Find("div#fact").Text()
	if fact != "" {
		log.Info("RAW-data с сайта: " + fact)

		pattern := `\.[\w|а-яА-Я]+нтересно`

		re := regexp.MustCompile(pattern)
		fact = re.ReplaceAllString(fact, ".")
	}

	log.Info(fmt.Sprintf("Информация с сайта '%s' успешно получена:\n%s", url, fact))

	return fact
}
