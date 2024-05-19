# Watcher

Watches a directory / file and runs a series of commands when changes happen. Easily configurable from a config file.

Use cases:

- auto-recompilation of file

## Config

Setup a `watcher.yml` file in the current working directory. Alternatively, pass in the config file via `--config-file <path/to/config.yml>`

```yaml
include:
    - src
exclude:
    - src/exclude.zig
commands:
    - zig build
```

Files watched = files that match include *and* do not match exclude

## Platforms Supported

- macos / darwin (OS X v10.5 and later)

## References

### Darwin

Uses the File System Events API

- https://developer.apple.com/library/archive/documentation/Darwin/Conceptual/FSEvents_ProgGuide/Introduction/Introduction.html#//apple_ref/doc/uid/TP40005289
