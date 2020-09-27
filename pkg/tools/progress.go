package tools

import (
    "fmt"
    "github.com/xiaotian/dgraph-person/pkg/d_log"
)

var logger = d_log.New()

func ShowProgress(title string, currLength int, totalLength int) {
    progress := fmt.Sprintf("%6.2f", float64(currLength)/float64(totalLength)*100)
    logger.Infof("%s, progress: [%s%s]", title, progress, "%")
}
