package scans

import (
	"strconv"
	"strings"

	"github.com/Lutz-Pfannenschmidt/yagll"
)

// TransformPortRange transforms a port range string into a format that nmap can understand
//
// "100~4" -> "100-104"
func TransformPortRange(input string) string {
	input = strings.ReplaceAll(input, " ", "")
	ranges := strings.Split(input, ",")
	for i, rangeIn := range ranges {
		if strings.Contains(rangeIn, "~") {
			j := strings.Index(rangeIn, "~")
			start, _ := strconv.Atoi(rangeIn[:j])
			end, _ := strconv.Atoi(rangeIn[j+1:])
			ranges[i] = rangeIn[:j] + "-" + strconv.Itoa(start+end)
		}
	}
	yagll.Debugf("Transformed port range '%s' to '%s'", input, strings.Join(ranges, ", "))
	return strings.Join(ranges, ", ")
}
