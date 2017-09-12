# queue

[![License][License-Image]][License-Url] [![ReportCard][ReportCard-Image]][ReportCard-Url] [![Build][Build-Status-Image]][Build-Status-Url] [![Coverage][Coverage-Image]][Coverage-Url] [![GoDoc][GoDoc-Image]][GoDoc-Url]

## Get

``` bash
go get -u github.com/LyricTian/queue
```

## Usage

``` go
package main

import (
    "fmt"

    "github.com/LyricTian/queue"
)

func main() {
    queue.Run(10, 100)

    job := queue.NewSyncJob("hello", func(v interface{}) (interface{}, error) {
        return fmt.Sprintf("%s,world", v), nil
    })
    queue.Push(job)

    result := <-job.Wait()
    if err := job.Error(); err != nil {
        panic(err)
    }

    fmt.Println(result)

    // output: hello,world
}
```

## MIT License

``` text
    Copyright (c) 2017 Lyric
```

[License-Url]: http://opensource.org/licenses/MIT
[License-Image]: https://img.shields.io/npm/l/express.svg
[Build-Status-Url]: https://travis-ci.org/LyricTian/queue
[Build-Status-Image]: https://travis-ci.org/LyricTian/queue.svg?branch=master
[ReportCard-Url]: https://goreportcard.com/report/github.com/LyricTian/queue
[ReportCard-Image]: https://goreportcard.com/badge/github.com/LyricTian/queue
[GoDoc-Url]: https://godoc.org/github.com/LyricTian/queue
[GoDoc-Image]: https://godoc.org/github.com/LyricTian/queue?status.svg
[Coverage-Url]: https://coveralls.io/github/LyricTian/queue?branch=master
[Coverage-Image]: https://coveralls.io/repos/github/LyricTian/queue/badge.svg?branch=master