# Spell Check

1. build the dictionary
2. run the spell check

```sh
$ spellcheck [-build <dict.txt>] [-dict <dict.sc>] <words> ...
$ cat <file> | spellcheck [-build <dict.txt>] [-dict <dict.sc>]
```

## compiling

```sh
$ zig build
```

## building the dictionary

The dictionary must be built first into a `sc` file first (file extension doesn't actually matter).

The destination defaults to `dict.sc`. However, another destination can be optionally specified.

Note that `-build` flag must be the **first** argument to the program. If this flag is provided, the program will **not** run spell check.

Sample `dict.txt` and `dict.sc` have been provided in the repository.

```sh
$ cp /usr/share/dict/words dict.txt # available on Unix-like systems
$ spellcheck -build dict.txt  # or any other list of words that are newline separated
$ spellcheck -build dict.txt -dict dict.sc
```

The dictionary file defaults to `dict.sc`. However another dictionary file can be optionally specified.

Note that `-dict` flag must be the **first** argument to the program.

```sh
$ spellcheck -dict dict.sc
```

## inputs

It accepts inputs via command line arguments.

```sh
# check all arguments
$ spellcheck "hi hellooo" word concurrency coding challenges tpyo up
>Misspelt:
>  hellooo
>  coding
>  challenges
>  tpyo
```

It accepts piped input, if no `words` are passed to it.

```sh
$ cat file.txt | spellcheck
```

## tests

```sh
zig build test
```

https://codingchallenges.fyi/challenges/challenge-bloom/

## extensions

- accept files / directories
    - print the file / line / column information
    - try parallelism in zig ([std.Thread](https://ziglang.org/documentation/master/std/#std.Thread)?)
        - doesn't seem to have many other options currently
    - this should be the only way? to accept input (might want to nuke the current options)
