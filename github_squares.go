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

type Month struct {
	x    byte
	name string
}

func (self Month) Initial() string {
	return string(self.name[0])
}

type Contributions struct {
	rects         [7][54]*Rect
	yearContrib   *NumInfo
	longestStreak *NumInfo
	currentStreak *NumInfo
	month         [14]*Month
}

func NewContributions(userID string) *Contributions {
	contribDoc, _ := goquery.NewDocument(
		fmt.Sprintf("https://github.com/users/%s/contributions", userID))

	column := 0
	rects := [7][54]*Rect{}
	contribDoc.Find("rect").Each(func(_ int, s *goquery.Selection) {
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
	var months [14]*Month // sometimes 12 or 13
	contribDoc.Find("text").Each(func(_ int, s *goquery.Selection) {
		attr, exists := s.Attr("class")
		if exists && attr == "month" {
			strX, _ := s.Attr("x")
			x, _ := strconv.Atoi(strX)
			months[m] = &Month{byte(x / 13), s.Text()}
			m++
		}
	})

	doc, _ := goquery.NewDocument(
		fmt.Sprintf("https://github.com/%s", userID))
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
	return &Contributions{rects, yearNum, streaks[0], streaks[1], months}
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

func (self Contributions) GetString(symbol string) (ans string) {
	ans = "  "
	m := 0
	for col := 0; col < 54 && self.month[m] != nil; col++ {
		if col == int(self.month[m].x) {
			ans += self.month[m].Initial()
			m++
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
				ans += Changer.Apply(symbol)
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
	argNum := len(os.Args)
	if argNum >= 2 {
		contrib := NewContributions(os.Args[1])
		symbol := "■"
		if argNum >= 4 && os.Args[2] == "-c" {
			symbol = string(os.Args[3][0])
		}
		str := contrib.GetString(symbol)
		fmt.Println(str)
	}
}
