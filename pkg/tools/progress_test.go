package tools

import "testing"

func TestShowProgress(t *testing.T) {
    length := 1000

    for i := 0; i < length; i++ {
        ShowProgress("TestShowProgress", i+1, length)
    }

}
