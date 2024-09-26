package processor

import (
	"fmt"
	"testing"
)

type TestCase struct {
	input    string
	expected string
}

func TestLoadRules(t *testing.T) {
	rules := LoadRules()
	want := true
	got := len(rules) > 0
	if want != got {
		t.Errorf("No rules loaded")
	}
}

func TestCheckEndTagPresent(t *testing.T) {
	testedText := "THis is a test\n\t02/12\n-end of document-\n"
	want := ""
	got := CheckEndTag(testedText)
	if want != got {
		t.Errorf("Wanted empty string, got %s", got)
	}
}

func TestCheckEndTagAbsent(t *testing.T) {
	testedText := "THis is a test\n\t02/12\n-end of document-dfs"
	want := "End tag not found!"
	got := CheckEndTag(testedText)
	if want != got {
		t.Errorf("Wanted %s, got %s", want, got)
	}
}

func TestCheckDatesFormat(t *testing.T) {
	badResponse := `Seems like a suspicious date: %s`
	tests := []TestCase{
		{"I was born on 09/03/1975 ", ""},
		{"I was born on 03.09.1975.", fmt.Sprintf(badResponse, "03.09.1975\n")},
		{"He was born on 24.12.1980", fmt.Sprintf(badResponse, "24.12.1980\n")},
		{"424.12.1980", ""},
		{"I was born on 25/09/1975.", fmt.Sprintf(badResponse, "25/09/1975\n")},
		{"I was born on 09-12-1975.", fmt.Sprintf(badResponse, "09-12-1975\n")},
	}
	for _, test := range tests {
		want := test.expected
		got := CheckDatesFormat(test.input)
		if want != got {
			t.Errorf("Wanted: %s, got: %s", want, got)
		}
	}
}

func TestRuleChecking(t *testing.T) {

	testCases := [][]TestCase{
		{
			{"Mrs. Billy does it willy-nilly", "Mrs."},
			{"Mr Bill is a nuthead", "Mr "},
			{"Mr. Billy does it willy-nilly", ""},
			{"Ms Smith is here", "Ms "},
			{"Jams happen", ""},
		},
		{
			{"The present document is a fake", "present"},
			{"She presented a male child", ""},
		},
		{
			{"The said document is a fake", "said"},
			{"She said goodbye", "said"},
			{"This invoice is outstanding", ""},
		},
		{
			{"This was done during day-time", "ing day-time"},
			{"She danced during the daytime", "ing the daytime"},
			{"During the day nothing happened", ""},
		},
		{
			{"Great amount of physical training", "physical training"},
			{"He received A in Physical Education.", ""},
		},
		{
			{"She presented her natural son", "natural "},
			{"Naturally, he received A", "Naturall"},
			{"He is her natural-born son", ""},
		},
		{
			{"Rate: 3 per cent", "per cent"},
			{"the per-cent rate", "per-cent"},
			{"Percentage is good", ""},
			{"I'll take ten percent", ""},
		},
	}
	rules := LoadRules()
	for i, testCase := range testCases {
		for _, test := range testCase {
			var want string
			if test.expected != "" {
				want = fmt.Sprintf("%s\n\t%s\n\n", rules[i][1], test.expected)
			} else {
				want = ""
			}
			got := CheckRule(rules, i, test.input)
			if want != got {
				t.Errorf("Wanted: %s, got: %s", want, got)
			}
		}
	}
}
