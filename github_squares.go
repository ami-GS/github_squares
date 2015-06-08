package github_squares

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
)

type Rect struct {
	color string
	count byte
	date  string
}

func GetData(reqUrl string) (results [7][54]Rect) {
	doc, _ := goquery.NewDocument(reqUrl)
	column := 0

	doc.Find("rect").Each(func(_ int, s *goquery.Selection) {
		yTmp, _ := s.Attr("y")
		y, _ := strconv.Atoi(yTmp)
		color, _ := s.Attr("fill")
		countTmp, _ := s.Attr("data-count")
		count, _ := strconv.Atoi(countTmp)
		date, _ := s.Attr("data-date")
		results[y/13][column] = Rect{color, byte(count), date}
		if y == 78 {
			column++
		}
	})
	return
}

func GetString(rects [7][54]Rect) (ans string) {
	for row := 0; row < 7; row++ {
		for col := 0; col < 54; col++ {
			if rects[row][col].date != "" {
				if rects[row][col].count != 0 {
					ans += "■"
				} else {
					ans += "□"
				}
			} else {
				ans += " "
			}
		}
		ans += "\n"
	}
	return
}

func ShowSquare(reqUrl string) {
	rects := GetData(reqUrl)
	str := GetString(rects)
	fmt.Println(str)
}
