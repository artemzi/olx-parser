package helpers

import (
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

var month = map[string]string{
	"января":   "Jan",
	"февраля":  "Feb",
	"марта":    "Mar",
	"апреля":   "Apr",
	"мая":      "May",
	"июня":     "Jun",
	"июля":     "Jul",
	"августа":  "Aug",
	"сентября": "Sep",
	"октября":  "Oct",
	"ноября":   "Nov",
	"декабря":  "Dec",
}

// PrettifyString removes repeated whitespaces from string
func PrettifyString(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func ParseMeta(s string) (time.Time, string) {
	data := strings.Split(s, ",")[1:] // [date, adverts_id]

	date := strings.Split(strings.TrimSpace(data[0]), " ")
	date[1] = month[date[1]]

	const shortForm = "02-Jan-2006"
	t, err := time.Parse(shortForm, strings.Join(date, "-"))
	if err != nil {
		log.Fatal(err)
	}
	return t, strings.TrimSpace(strings.Split(data[1], ":")[1])
}
