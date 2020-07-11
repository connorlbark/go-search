# go-search

Implementation of search algorithms in golang.

Inspired by [go-astar](https://github.com/beefsack/go-astar)
but with the goal of adding more algorithms and flexibility.

## CLI Usage

```bash
$ go run ./cmd/go-search --on <environment> --with <algorithm>
```

For example, with the `maze` environment

Legend:
- `x`: impassable
- `.`: passable
- `●`: path):

```bash
$ go run ./cmd/go-search --on maze with 'a*'
Found node (29,9) in 135 interations.
Steps (61): start, right, right, right, down, right, right, right, up, right, right, right, right, right, right, down, down, down, down, down, right, right, right, up, up, right, right, right, down, right, right, down, right, right, right, up, up, left, up, up, right, right, down, right, right, up, up, right, right, right, down, down, down, left, down, down, right, down, down, down, down
Total cost of solution: 61
●●●●x.●●●●●●●xxxxxxxxxxxxx●●●●
.xx●●●●xx.xx●xxxxxxx..●●●x●xx●
.x..xx....xx●xxx.xxx..●x●●●xx●
..xxx.x.xx.x●x.●●●●xxx●●xxxx●●
.xx.....xx.x●xx●xx●●●xx●xxxx●x
..x.xxxxxx..●●●●xxxx●●●●xxxx●●
x.x..xx..x.xxxx...xxxxxx...xx●
...x.xxx...xxxxxx..xx....x.xx●
xx...xxxxxxxxxxxxx.xxx.xx..xx●
```

## Package Usage

Basic usage, using premade
algorithms and environments: 

```go
package main

import (
    "github.com/porgull/go-search/pkg/environments"
    "github.com/porgull/go-search/pkg/algorithms"
    "github.com/porgull/go-search/pkg/search"

)


func main() {
    algorithm, err := algorithms.GetAlgorithm("a*")
    if err != nil {
        panic(err)
    }

    // You can also make your own environments.
    // See the godoc for the details on the interface.
    env, err := environments.GetEnvironment("maze")
    if err != nil {
        panic(err)
    }

    result, err := algorithm.Run(env)
    if err != nil {
        panic(err)
    }

    result.Print() // print statistics out
}
```

However, you can also implement your own algorithms and
environments.

See the `environments.Environment` interface and the
`algorithms.Algorithm` interface for details.
