package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// Struct to represent anchor (<a href=" Anchor.Url " ...>Anchor.Name</>) tag.
type Anchor struct {
	Name, Url string
}
type Pre struct {
	RevBy, Message string
}

const anyUNICODE = `\p{L}\p{Z}\p{C}\p{N}\p{S}\p{P}\p{M}`

// Method for unwrapping anchor (<a />) tag. Return a pointer to anchor struct
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

// Parse content from table data tag. Format: <td>tdBody</td>
func UnwrapTd(str string) (tdBody string) {
	return strings.TrimLeft(strings.TrimRight(str, "</td>"), "<td>")
}

// Parse content from pre tag. Returns a pointer to Pre struct. Format: <pre ....> Pre.Message BUG:... Reviewed-by: Pre.RevBy ...
func UnwrapPre(str string) (*Pre, error) {
	pre := &Pre{}
	str = strings.Trim(str, " ")
	reg := regexp.MustCompile(`<pre[ -~]*>([` + anyUNICODE + `]*)BUG`)

	mesg := reg.FindStringSubmatch(str)
	if len(mesg) < 1 {
		// Preference for extracting the exact commit message. In case it is a revert commit, send the content of complete pre as message body
		mesg = regexp.MustCompile(`<pre[ -~]*>([` + anyUNICODE + `]*)</pre>`).FindStringSubmatch(str)
		if len(mesg) < 1 {
			fmt.Println(mesg)
			return nil, fmt.Errorf("invalid branch name/url: error while parsing commit message (No term 'BUG')")
		}
	}
	pre.Message = mesg[1]

	reg2 := regexp.MustCompile(`[ -~]*Reviewed-by:([ -~]*)`)
	revby := reg2.FindStringSubmatch(str)
	if len(revby) < 1 {
		fmt.Println(mesg)
		return nil, fmt.Errorf("invalid branch name/url: error while parsing commit message (No Reviewed-By)")
	}
	pre.RevBy = revby[1]

	return pre, nil
}
