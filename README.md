# go-search

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/github.com/porgull/go-search)

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
Found node (29,9) in 135 iterations.
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

## Provided Search Algorithms

Uninformed Search Algorithms:
- TODO

Informed Search Algorithms
- [Greedy Best First Search](https://en.wikipedia.org/wiki/Best-first_search#Greedy_BFS) (key: `greedy_best_first`): Uses only heuristic
- [A*](https://en.wikipedia.org/wiki/A*_search_algorithm) (key: `a*`): Uses both cost and heuristic

Local Search Algorithms:
- TODO

## Provided Environments

Some of the environments are defined
by loading in JSON files. See
`assets/environments` for the
files that define the pre-made
environments if you wish to make
your own.

### GridEnvironment

Here's an example definition of
a grid environment:

```
*...x........xxxxxxxxxxxxx....
.xx....xx.xx.xxxxxxx.....x.xx.
.x..xx....xx.xxx.xxx...x...xx.
..xxx.x.xx.x.x.....xxx..xxxx..
.xx.....xx.x.xx.xx...xx.xxxx.x
..x.xxxxxx......xxxx....xxxx..
x.x..xx..x.xxxx...xxxxxx...xx.
...x.xxx...xxxxxx..xx....x.xx.
xx...xxxxxxxxxxxxx.xxx.xx..xx.
xxxxxxxxxxxxxxxxxx.....xx...x!
```

Legend:
- `*`: Start
- `!`: End
- `.`: Passable, cost 1
- `,`: Passable, cost 2
- `#`: Passable, cost 3
- `x`: Impassable

The heuristic for this environment is 
the [Manhattan Distance](https://en.wikipedia.org/wiki/Taxicab_geometry)
to the goal node.

Pre-made Grid environments:
- `corners`: Simply has to traverse to the corner
- `maze`: Basic maze

### StateEnvironment

A state environment is an environment
where every state can be loaded
into memory.

In other words, the states can 
be enumerated like so:

```json5
...
"arad": {
    "heuristic": 366,
    "children": { // key = name, val = cost to node
        "zerind": 75, 
        "timisoara": 118,
        "sibiu": 140
    }
},
"zerind": {
    "heuristic": 374,
    "children": {
        "oradea": 71,
        "arad": 75
    }
},
...
```

Pre-made State environments:
- `bucharest`: From the 3rd Edition of
AI: A Modern Approach by Stuart J.
Russell and Peter Norvig