package main

import (
	"fmt"
	"log"
	"os"

	yq "bla/yq"

	"github.com/itchyny/gojq"
)

var (
	outputIndent = 2
)

func process(iter yq.InputIter, code *gojq.Code) error {
	var err error
	for {
		v, ok := iter.Next()
		if !ok {
			return err
		}
		if er, ok := v.(error); ok {
			printError(er)
			err = &yq.EmptyError{Err: er}
			// log.Fatalln(er) //todo ?
			continue
		}
		// TODO: if er := cli.printValues(code.Run(v, cli.argvalues...)); er != nil {
		if er := printValues(code.Run(v)); er != nil {
			printError(er)
			err = &yq.EmptyError{Err: er}
		}
	}
}

// cli.printValues
// https://github.com/itchyny/gojq/blob/main/cli/cli.go#L356
func printValues(iter gojq.Iter) error {
	m := createMarshaler()
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return err
		}

		if err := m.Marshal(v, os.Stdout); err != nil {
			return err
		}
	}
	return nil
}

// cli.printError
// https://github.com/itchyny/gojq/blob/main/cli/cli.go#L422
func printError(err error) {
	if er, ok := err.(interface{ IsEmptyError() bool }); !ok || !er.IsEmptyError() {
		if er, ok := err.(interface{ IsHaltError() bool }); !ok || !er.IsHaltError() {
			fmt.Fprintf(os.Stderr, "%s: %s\n", "cmdYq", err)
		} else if er, ok := err.(gojq.ValueError); ok {
			v := er.Value()
			if str, ok := v.(string); ok {
				os.Stderr.Write([]byte(str))
			} else {
				bs, _ := gojq.Marshal(v)
				os.Stderr.Write(bs)
				os.Stderr.Write([]byte{'\n'})
			}
		}
	}
}

// cli.createMarshaler
// https://github.com/itchyny/gojq/blob/main/cli/cli.go#L392
func createMarshaler() yq.Marshaler {
	return yq.YamlFormatter(&outputIndent)
}

func main() {
	query, err := gojq.Parse("del(.x.y.z)")
	if err != nil {
		log.Fatalln(err)
	}
	code, err := gojq.Compile(query)
	if err != nil {
		log.Fatalln(err)
	}

	file, err := os.Open("test.yml")
	if err != nil {
		log.Fatalln(err)
	}

	iter := yq.NewYAMLInputIter(file, "baba")

	err = process(iter, code)
	if err != nil {
		log.Fatalln(err)
	}
}
