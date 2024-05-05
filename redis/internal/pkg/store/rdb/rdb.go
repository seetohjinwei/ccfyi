package rdb

// *partially* follows the file format listed here: https://rdb.fnordig.de/file_format.html
// liberties taken because this is a Go program, some of these optimisations don't matter

const redis = "REDIS"
const version = "LITE" // intentionally not an integer to not collide with redis version numbers
const magicString = redis + version
