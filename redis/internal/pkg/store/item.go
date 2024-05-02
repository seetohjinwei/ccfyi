package store

type Item interface {
	Do(command string, args []string) (string, bool)
}
