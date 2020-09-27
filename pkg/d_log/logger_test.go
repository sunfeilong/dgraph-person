package d_log

import (
    "github.com/stretchr/testify/assert"
    "sync"
    "testing"
)

func TestFormat(t *testing.T) {
    logger := New()

    keyValues := make(map[string]string)
    keyValues["name"] = "zhangsan"
    keyValues["age"] = "24"

    //字符串匹拼接
    //INFO	d_log/logger_test.go:17	Info name name map map[age:24 name:zhangsan]
    logger.Info("Info ", "name ", "ZS ", "map ", keyValues)
    //JSON 格式
    //Infow 	{"name ": "ZS ", "map ": {"age":"24","name":"zhangsan"}}
    logger.Infow("Infow ", "name ", "ZS ", "map ", keyValues)
    //模板填充
    //Infow name: ZS , map: map[age:24 name:zhangsan]
    logger.Infof("Infof name: %s, map: %s", "ZS ", keyValues)
}

func TestName(t *testing.T) {
    logger := New()

    logger.Info("Info")
    logger.Error("Error")
    defer logger.Sync()
    assert.NotEmpty(t, logger, "日志信息不能为空")
}

func TestMultiSingle(t *testing.T) {
    logger := New()
    times := 102400
    for i := 0; i < times; i++ {
        logger.Infow("测试打印日志", "name", "name")
    }
}

func TestMultiOpen(t *testing.T) {
    waitGroup := sync.WaitGroup{}
    waitGroup.Add(2)
    go func() {
        logger := New()
        times := 102400
        for i := 0; i < times; i++ {
            logger.Infow("1111111111", "name", "name")
        }
        waitGroup.Done()
    }()

    go func() {
        logger := New()
        times := 102400
        for i := 0; i < times; i++ {
            logger.Infow("2222222222", "name", "name")
        }
        waitGroup.Done()
    }()

    waitGroup.Wait()
}
