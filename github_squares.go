package github_squares

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ami-GS/soac"
	"strconv"
)

var colorMap map[string]byte = map[string]byte{
	"#d6e685": 156,
	"#8cc665": 112,
	"#44a340": 34,
	"#1e6823": 22,
	"#eeeeee": 237,
}

var Changer *soac.Changer

func init() {
	Changer = soac.NewChanger()
}

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
				Changer.Set256(colorMap[rects[row][col].color])
				ans += Changer.Apply("■")
			} else {
				ans += " "
			}
		}
		ans += "\n"
	}
	return
}

func ShowSquare(userName string) {
	reqUrl := fmt.Sprintf("http://github.com/%s/", userName)
	rects := GetData(reqUrl)
	str := GetString(rects)
	fmt.Println(str)
}
