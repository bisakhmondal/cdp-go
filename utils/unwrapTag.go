package utils

import (
	"regexp"
	"strings"
)

type Anchor struct {
	Name, Url string
}
type Pre struct {
	RevBy, Message string
}

func UnwrapA(str string) *Anchor {
	str = strings.TrimRight(str, "</a>")
	splits := strings.SplitAfter(str, ">")
	name := splits[1]

	atag := strings.Replace(splits[0], " ", "", -1)
	urlStart := strings.SplitAfter(atag, "href=\"")
	url := strings.Split(urlStart[1], `"`)[0]

	return &Anchor{
		Name: name,
		Url:  url,
	}
}

func UnwrapTd(str string) string {
	return strings.TrimLeft(strings.TrimRight(str, "</td>"), "<td>")
}

func UnwrapPre(str string) *Pre {
	pre := &Pre{}
	str = strings.Trim(str, " ")
	reg := regexp.MustCompile(`<pre[ -~]*>([ -~\n]*)BUG`)

	mesg := reg.FindStringSubmatch(str)
	pre.Message = mesg[1]

	reg2 := regexp.MustCompile(`[ -~]*Reviewed-by:([ -~]*)`)
	revby := reg2.FindStringSubmatch(str)
	pre.RevBy = revby[1]

	return pre
}
