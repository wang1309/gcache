
<img width="298" alt="image" src="https://github.com/wang1309/gcache/assets/20272951/6f811e8c-8dc8-4340-9493-292f15841bd2">

#### What is gcache？

gcache is a high-performance Go language implementation of a cache library. It is easy to use and  understand

#### Design principle?

- LRU elimination algorithm
- Use singleflight to avoid cache pinging
- Support custom load method
- Support cache change listener



##### TODO

Map shard

#### How to install ?

go get -u  github.com/wang1309/gcache



#### use case

**init cache object**

```go
import (
	"sync"
	"time"
	"golang.org/x/sync/singleflight"
	"github.com/wang1309/gcache"
)

const maxItems = 100

func main() {
	// create cache object
	c := cache.NewCache(loader, maxItems)
}

func loader(key string) interface{} {
	// get data by key
	return "value"
}

```



**Get**

```go
value := c.Get("key")
```



**Put**

```go
c.Put("key", "value")
```



#### Benchmark

##### ENV

goos: darwin
goarch: arm64

```
BenchmarkCache_GetParallel/Get-8                 3000000               314.4 ns/op
BenchmarkCache_Put-8     3000000               190.4 ns/op
```


