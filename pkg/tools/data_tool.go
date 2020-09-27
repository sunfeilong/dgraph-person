package tools

import "regexp"

var numberRegex = "\\d+"

func IsNumber(value string) bool {

    match, err := regexp.Match(numberRegex, [] byte(value))
    if err != nil {
        return false
    }

    return match
}
