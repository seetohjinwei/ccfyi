# Spell Check

1. build the dictionary
2. run the spell check

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
$ spellcheck "hi hello" word concurrency coding challenges tpyo up
> Misspelt:
>   coding
>   challenges
>   tpyo
```

It accepts piped input.

```sh
$ cat file.txt | spellcheck
```

## tests

```sh
zig build test
```

https://codingchallenges.fyi/challenges/challenge-bloom/
