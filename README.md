# snowflake
Snowflake算法Golang实现

# 安装

```
go get github.com/Anderson-Lu/snowflake
```

# 使用

```go
import "github.com/Anderson-Lu/snowflake"

func main() {
    workerID := 111
    x, e := snowflake.NewIDGenerator(workerID)
    if e != nil {
        panic(e)
    }
    id := x.x.GenerateID()
    //...
}
```