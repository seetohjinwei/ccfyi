package init_op // because of name clash with `init`

import (
	"fmt"
	"os"
	"path/filepath"
)

type operation struct {
	path string
}

type Options struct {
	Quiet bool
}

func (op *operation) getPath() error {
	path, err := filepath.Abs(op.path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "git-init: failed to get absolute path for '%v'\n", op.path)
		return err
	}

	op.path = path + "/.git/"

	return nil
}

func (op *operation) makeDir() error {
	err := os.MkdirAll(op.path, os.ModePerm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "git-init: failed to create directory '%v'\n", op.path)
		return err
	}

	return nil
}

func (op *operation) makeHead() error {
	// TODO:

	return nil
}

func (op *operation) makeConfig() error {
	const content = `[core]
	repositoryformatversion = 0
	filemode = true
	bare = false
	logallrefupdates = true
	ignorecase = true
	precomposeunicode = true
`
	const filepath = "config"

	err := makeFile(op.path+filepath, content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "git-init: failed to create config file\n")
		return err
	}

	return nil
}

func (op *operation) makeDescription() error {
	const content = "Unnamed repository; edit this file 'description' to name the repository.\n"
	const filepath = "description"

	err := makeFile(op.path+filepath, content)
	if err != nil {
		fmt.Fprintf(os.Stderr, "git-init: failed to create description file\n")
		return err
	}

	return nil
}

func (op *operation) makeHooks() error {
	// TODO:

	return nil
}

func (op *operation) makeInfo() error {
	// TODO:

	return nil
}

func (op *operation) makeObjects() error {
	// TODO:

	return nil
}

func (op *operation) makeRef() error {
	// TODO:

	return nil
}

func makeFile(path, content string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

type step func() error

func Perform(path string, opts Options) error {
	op := operation{
		path: path,
	}

	steps := []step{
		op.getPath,
		op.makeDir,
		op.makeHead,
		op.makeConfig,
		op.makeDescription,
		op.makeHooks,
		op.makeInfo,
		op.makeObjects,
		op.makeRef,
	}

	for _, s := range steps {
		err := s()
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(os.Stdout, "git-init: Initialized empty Git repository in %v\n", op.path)

	return nil
}
