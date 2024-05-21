# Prerequisites

1. You need to install. Check https://go.dev/doc/install for more details. Check by running `go version` from command-line. 
2. (optionaly) Download some text file (around 3-5mb). You can use [Project Gutenberg](https://www.gutenberg.org/)

# Building

```
cd olostan/step6
go build .
```

# Running

```
./babyai --help
```

It should output:

```
Usage of ./babyai:
  -file filename
        filename to parse (default "books/comb.txt")
  -i    interactive mode
  -length int
        number of characters to generate (default 100)
  -random
        randomize output (default true)
  -tokens length
        token length (how many characters to take into account) (default 10)
```

In non-interactive mode you can run:

```
./babyai hello there
```
(where "hello there" - a text to generate text based on)

You can also disable output randomization by using `--random=false` parameter to produce more predictable text.