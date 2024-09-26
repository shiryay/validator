package processor

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var Rules [][]string

const RulesFile string = "./rules.json"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func LoadRules() [][]string {
	// Read rules file and return the structure.
	var rules [][]string
	file, err := os.Open(RulesFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&rules); err != nil {
		panic(err)
	}
	return rules
}

func CheckDatesFormat(checked string) string {
	returnString := ""
	wrongDatePattern := `\b(\d|\d\d)([\./-])(\d|\d\d)([\./-])(\d\d\d\d|\d\d)\b`
	wrongDateRegex, err := regexp.Compile(wrongDatePattern)
	check(err)
	matches := wrongDateRegex.FindAllStringSubmatch(checked, -1)
	for _, match := range matches {
		// The following line must be corrected
		month, err := strconv.Atoi(match[1])
		check(err)
		if string(match[2]) != "/" || string(match[4]) != "/" || month > 12 {
			returnString = returnString + fmt.Sprintf("Seems like a suspicious date: %s\n", match[0])
		}
	}
	return returnString
}

func CheckEndTag(checked string) string {
	trimmedString := strings.Replace(checked, "\r\n", "", -1)
	trimmedString = strings.Replace(trimmedString, "\n", "", -1)
	if strings.HasSuffix(trimmedString, "-end of document-") {
		return ""
	} else {
		return "Check the end of document tag!\n"
	}
}

func CheckRule(rules [][]string, ruleIndex int, checked string) string {
	returnString := ""
	stopPhrasePattern := ""
	var stopPhraseRegex *regexp.Regexp
	// var err error

	stopPhrasePattern = `(?i)` + rules[ruleIndex][0]
	stopPhraseRegex = regexp.MustCompile(stopPhrasePattern)
	matches := stopPhraseRegex.FindAllString(checked, -1)
	if len(matches) != 0 {
		returnString = rules[ruleIndex][1] + "\n\t" + strings.Join(matches, "\n\t") + "\n"
	}
	return returnString
}

func CheckText(text string) string {
	reportString := ""
	Rules = LoadRules()
	reportString += CheckDatesFormat(text)
	reportString += CheckEndTag(text)
	for i := range Rules {
		reportString = reportString + CheckRule(Rules, i, text)
	}
	return reportString
}
