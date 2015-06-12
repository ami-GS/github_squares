package github_squares

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/ami-GS/soac"
	"strconv"
	"strings"
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

type Contributions struct {
	rects         [7][54]Rect
	yearNum       uint16
	longestStreak uint16
	currentStreak uint16
	month         [12]string
}

func (self Contributions) Get(row, column int) Rect {
	return self.rects[row][column]
}

type Rect struct {
	color string
	count byte
	date  string
}

func GetData(reqUrl string) (contrib Contributions) {
	doc, _ := goquery.NewDocument(reqUrl)
	column := 0
	rects := [7][54]Rect{}
	doc.Find("rect").Each(func(_ int, s *goquery.Selection) {
		yTmp, _ := s.Attr("y")
		y, _ := strconv.Atoi(yTmp)
		color, _ := s.Attr("fill")
		countTmp, _ := s.Attr("data-count")
		count, _ := strconv.Atoi(countTmp)
		date, _ := s.Attr("data-date")
		rects[y/13][column] = Rect{color, byte(count), date}
		if y == 78 {
			column++
		}
	})

	m := 0
	var month [12]string
	doc.Find("text").Each(func(_ int, s *goquery.Selection) {
		attr, exists := s.Attr("class")
		if exists && attr == "month" {
			month[m] = s.Text()
			m++
		}
	})

	var yearNum uint16
	var streaks [2]uint16
	doc.Find("div[class='contrib-column contrib-column-first table-column']").Each(func(_ int, s *goquery.Selection) {
		text := s.Find("span[class='contrib-number']").Text()
		result := strings.Split(text, " ")
		num, _ := strconv.Atoi(result[0])
		yearNum = uint16(num)
	})

	streakIdx := 0
	doc.Find("div[class='contrib-column table-column']").Each(func(_ int, s *goquery.Selection) {
		text := s.Find("span[class='contrib-number']").Text()
		result := strings.Split(text, " ")
		num, _ := strconv.Atoi(result[0])
		streaks[streakIdx] = uint16(num)
		streakIdx++
	})

	contrib = Contributions{rects, yearNum, streaks[0], streaks[1], month}
	return
}

func GetString(contrib Contributions) (ans string) {
	ans = "  " + string(contrib.month[0][0])
	m := 1
	rect := contrib.Get(6, 0) // investigate first column month
	mStr := strings.Split(rect.date, "-")
	prev := mStr[1]
	for col := 1; col < 54; col++ {
		rect = contrib.Get(0, col)
		mStr = strings.Split(rect.date, "-")
		if len(mStr) >= 2 && mStr[1] != prev {
			ans += string(contrib.month[m][0])
			prev = mStr[1]
			m++
			if m == 12 {
				break
			}
		} else {
			ans += " "
		}
	}
	ans += "\n"

	for row := 0; row < 7; row++ {
		switch {
		case row == 1:
			ans += "M "
		case row == 3:
			ans += "W "
		case row == 5:
			ans += "F "
		default:
			ans += "  "
		}

		for col := 0; col < 54; col++ {
			rect := contrib.Get(row, col)
			if rect.date != "" {
				Changer.Set256(colorMap[rect.color])
				ans += Changer.Apply("■")
			} else {
				ans += " "
			}
		}
		ans += "\n"
	}

	ans += "========================================================\n"
	ans += fmt.Sprintf("Contributions in the last year\n   %d total\ndummy -dummy\n", contrib.yearNum)
	ans += "--------------------------------------------------------\n"
	ans += fmt.Sprintf("Longest streak\n   %d days\ndummy -dummy\n", contrib.longestStreak)
	ans += "--------------------------------------------------------\n"
	ans += fmt.Sprintf("Current streak\n   %d days\ndummy -dummy\n", contrib.longestStreak)
	ans += "========================================================\n"

	return
}

func ShowSquare(userName string) {
	reqUrl := fmt.Sprintf("http://github.com/%s/", userName)
	contrib := GetData(reqUrl)
	str := GetString(contrib)
	fmt.Println(str)
}
