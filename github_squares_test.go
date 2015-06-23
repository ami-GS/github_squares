package github_squares

import (
	"fmt"
	"reflect"
	"testing"
)

var numInfo *NumInfo

func init() {
	infoStr := "Contributions in the last year\n920 total\nJun 23, 2014 – Jun 23, 2015"
	numInfo = NewNumInfo(infoStr)

}

func TestNewRect(t *testing.T) {
	color := "#d6e685"
	count := byte(100)
	date := "2015-06-15"
	actual := NewRect(color, count, date)
	expect := &Rect{color, count, date}
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("got %v\nwant %v", actual, expect)
	}
}

func TestNewNumInfo(t *testing.T) {
	infoStr := "Contributions in the last year\n920 total\nJun 23, 2014 – Jun 23, 2015"
	title := "Contributions in the last year"
	num := "920 total"
	days := "Jun 23, 2014 – Jun 23, 2015"
	actual := NewNumInfo(infoStr)
	expect := &NumInfo{title, num, days}
	if !reflect.DeepEqual(actual, expect) {
		t.Errorf("got %v\nwant %v", actual, expect)
	}
}

func TestNewContributions(t *testing.T) {
	user := "ami-GS"
	reqUrl := fmt.Sprintf("http://github.com/%s/", user)
	actual := NewContributions(reqUrl)
	infoStr := "Contributions in the last year\n920 total\nJun 23, 2014 – Jun 23, 2015"
	var rects [7][54]*Rect
	var month [14]string
	expect := &Contributions{rects, NewNumInfo(infoStr), NewNumInfo(infoStr), NewNumInfo(infoStr), month}
	if actual == expect {
		t.Errorf("got %v\nwant %v", actual, expect)
	}
}

func TestString(t *testing.T) {
	actual := numInfo.String()
	expect := "\tContributions in the last year\n\t\t920 total\n\tJun 23, 2014 – Jun 23, 2015\n"
	if actual != expect {
		t.Errorf("got %v\nwant %v", actual, expect)
	}
}
