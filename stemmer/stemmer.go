package stemmer

import (
    "github.com/kljensen/snowball"
    "strings"
)

func getStem(originalWord, language string) (string, error) {
    stemmedWord, err := snowball.Stem(originalWord, language, true)
    if err != nil {
        return "", err
    } else {
        return stemmedWord, nil
    }
}

func stemWord(originalWord, language string) (string, error) {
    stem, err := getStem(originalWord, language)
    if err != nil {
        return "", err
    } else {
        insert := stem + "|"
        split := strings.Replace(originalWord, stem, insert, 1)
        return split, nil
    }
}

func stemPhrase(originalPhrase, language string) (string, error) {
    words := strings.Split(originalPhrase, " ")
    splitWords := make([]string, len(words))
    for _, word := range words {
        splitWord, err := stemWord(word, language)
        if err != nil {
            return "", err
        } else {
            splitWords = append(splitWords, splitWord)
        }
    }
    return strings.Join(splitWords, " "), nil
}

func stemLineWithTabs(line, language string) (string, error) {
    lineParts := strings.Split(line, "\t")
    newLineParts := make([]string, 0)
    stemmedPart, err := stemPhrase(lineParts[0], language)
    if err != nil {
        return "", err
    } else {
        newLineParts = append(newLineParts, stemmedPart)
        newLineParts = append(newLineParts, lineParts[1])
        return strings.Join(newLineParts, "\t"), nil
    }
}

func StemText(text, language string) string {
    sep := "\n" // may need to be different for web, \r\n or something like that
    lines := strings.Split(text, sep)
    newLines := make([]string, 0)
    for _, line := range lines {
        stemmedLine, err := stemLineWithTabs(line, language)
        if err != nil {
            return err.Error()
        } else {
            newLines = append(newLines, stemmedLine)
        }
    }
    return strings.Join(newLines, sep)
}

// TODO:
// 1. Implement removal of cases like "jogger|", i.e. with a pipe at word end
// 2. Make sure tabs do not appear at the beginning of each line in the resulting text
