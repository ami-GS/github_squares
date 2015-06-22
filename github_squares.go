package github_squares

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"strconv"
	"strings"
)

type NumInfo struct {
	title string
	num   string
	days  string
}

func NewNumInfo(infoStr string) *NumInfo {
	text := strings.Replace(infoStr, "–\n              ", "– ", 1)
	days := ""
	title := ""
	num := ""
	for _, v := range strings.Split(text, "\n") {
		if len(v) != 0 && len(title) == 0 {
			title = strings.TrimSpace(v)
		} else if len(title) != 0 && len(num) == 0 {
			num = strings.TrimSpace(v)
		} else if strings.Contains(v, "–") {
			days += strings.TrimSpace(v)
		}
	}
	return &NumInfo{title, num, days}
}

func (self NumInfo) String() string {
	return fmt.Sprintf("\t%s\n\t\t%s\n\t%s\n", self.title, self.num, self.days)
}

type Contributions struct {
	rects         [7][54]*Rect
	yearContrib   *NumInfo
	longestStreak *NumInfo
	currentStreak *NumInfo
	month         [14]string
}

func NewContributions(reqUrl string) *Contributions {
	doc, _ := goquery.NewDocument(reqUrl)
	column := 0
	rects := [7][54]*Rect{}
	doc.Find("rect").Each(func(_ int, s *goquery.Selection) {
		yTmp, _ := s.Attr("y")
		y, _ := strconv.Atoi(yTmp)
		color, _ := s.Attr("fill")
		countTmp, _ := s.Attr("data-count")
		count, _ := strconv.Atoi(countTmp)
		date, _ := s.Attr("data-date")
		rects[y/13][column] = NewRect(color, byte(count), date)
		if y == 78 {
			column++
		}
	})

	m := 0
	var month [14]string // sometimes 12 or 13
	doc.Find("text").Each(func(_ int, s *goquery.Selection) {
		attr, exists := s.Attr("class")
		if exists && attr == "month" {
			month[m] = s.Text()
			m++
		}
	})

	var yearNum *NumInfo
	var streaks [2]*NumInfo
	doc.Find("div[class='contrib-column contrib-column-first table-column']").Each(func(_ int, s *goquery.Selection) {
		yearNum = NewNumInfo(s.Text())
	})

	streakIdx := 0
	doc.Find("div[class='contrib-column table-column']").Each(func(_ int, s *goquery.Selection) {
		streaks[streakIdx] = NewNumInfo(s.Text())
		streakIdx++
	})
	return &Contributions{rects, yearNum, streaks[0], streaks[1], month}
}

func (self Contributions) Get(row, column int) *Rect {
	return self.rects[row][column]
}

type Rect struct {
	color string
	count byte
	date  string
}

func NewRect(color string, count byte, date string) *Rect {
	return &Rect{color, count, date}
}

func (self Contributions) GetString() (ans string) {
	ans = "  " + string(self.month[0][0])
	m := 1
	rect := self.Get(6, 0) // investigate first column month
	mStr := strings.Split(rect.date, "-")
	prev := mStr[1]
	for col := 1; col < 54; col++ {
		rect = self.Get(0, col)
		mStr = strings.Split(rect.date, "-")
		if len(mStr) >= 2 && mStr[1] != prev {
			ans += string(self.month[m][0])
			prev = mStr[1]
			m++
			if self.month[m] == "" {
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
			rect := self.Get(row, col)
			if rect != nil && rect.date != "" {
				Changer.Set256(colorMap[rect.color])
				ans += Changer.Apply("■")
			} else {
				ans += " "
			}
		}
		ans += "\n"
	}

	ans += "========================================================\n"
	ans += self.yearContrib.String()
	ans += "--------------------------------------------------------\n"
	ans += self.longestStreak.String()
	ans += "--------------------------------------------------------\n"
	ans += self.currentStreak.String()
	ans += "========================================================\n"

	return
}

func ShowSquare() {
	if len(os.Args) >= 2 {
		reqUrl := fmt.Sprintf("http://github.com/%s/", os.Args[1])
		contrib := NewContributions(reqUrl)
		str := contrib.GetString()
		fmt.Println(str)
	}
}
