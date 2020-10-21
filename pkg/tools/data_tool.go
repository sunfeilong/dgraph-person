package tools

import "regexp"

var numberRegex = "^\\d{11}$"

func IsNumber(value string) bool {

    match, err := regexp.Match(numberRegex, [] byte(value))
    if err != nil {
        return false
    }

    return match
}
