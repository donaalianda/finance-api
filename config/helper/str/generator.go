package str

import (
	"regexp"
	"strings"
	"math"
)

func PhoneConvertToAbbv(phone string) string {
	isMatch, _ := regexp.MatchString(`^[0{1}]`, phone)
	if isMatch {
		re := regexp.MustCompile(`^[0{1}]`)
		s := re.ReplaceAllString(phone, `+62`)

		return s
	}

	return phone
}

func Replacer(source string, replacer *strings.Replacer) string {
	return replacer.Replace(source)
}

func CustomRound(x , correction float64) float64 {
    t := math.Trunc(x)
    if math.Abs(x-t) >= correction {
        return t + math.Copysign(1, x)
    }
    return t
}
