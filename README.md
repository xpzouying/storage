# storage
storage in go


## Usage

### Install

```bash
go get -u github.com/xpzouying/storage
```

### Example

```go

import "github.com/xpzouying/storage"

l, _ := storage.NewLocal("/tmp/oss")

buf := bytes.NewBuffer("hello storage")
l.Put(context.Background(), "file1", buf)


rc, _, := l.Get(context.Background(), "file1")
data, _ := ioutil.ReadAll(rc)
rc.Close()

log.Printf("read data: %s\n", data)
```


### Test

```bash
go test -v github.com/xpzouying/storage
```
