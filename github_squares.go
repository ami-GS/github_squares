package github_squares

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

type NumInfo struct {
	infoStr string
	num     uint16
}

func NewNumInfo(infoStr string, num uint16) *NumInfo {
	return &NumInfo{infoStr, num}
}

func (self NumInfo) String() string {
	return fmt.Sprintf("%s\n", self.infoStr)
}

type Contributions struct {
	rects         [7][54]*Rect
	yearContrib   *NumInfo
	longestStreak *NumInfo
	currentStreak *NumInfo
	month         [12]string
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
	var month [12]string
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
		text := s.Text()
		numText := s.Find("span[class='contrib-number']").Text()
		numResult := strings.Split(numText, " ")
		num, _ := strconv.Atoi(numResult[0])
		yearNum = NewNumInfo(text, uint16(num))
	})

	streakIdx := 0
	doc.Find("div[class='contrib-column table-column']").Each(func(_ int, s *goquery.Selection) {
		text := s.Text()
		numText := s.Find("span[class='contrib-number']").Text()
		numResult := strings.Split(numText, " ")
		num, _ := strconv.Atoi(numResult[0])
		streaks[streakIdx] = NewNumInfo(text, uint16(num))
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

func GetString(contrib *Contributions) (ans string) {
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
	ans += contrib.yearContrib.String()
	ans += "--------------------------------------------------------\n"
	ans += contrib.longestStreak.String()
	ans += "--------------------------------------------------------\n"
	ans += contrib.currentStreak.String()
	ans += "========================================================\n"

	return
}

func ShowSquare(userName string) {
	reqUrl := fmt.Sprintf("http://github.com/%s/", userName)
	contrib := NewContributions(reqUrl)
	str := GetString(contrib)
	fmt.Println(str)
}
