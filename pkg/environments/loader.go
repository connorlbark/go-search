package environments

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
)

var (
	environmentTypes = map[string]reflect.Type{}
)

func addEnvironmentType(name string, e Environment) {
	environmentTypes[name] = reflect.TypeOf(e)
}

// MustLoadEnvironmentFrom loads an environment into memory, panicking on an error
func MustLoadEnvironmentFrom(f io.Reader) Environment {
	env, err := LoadEnvironmentFrom(f)
	if err != nil {
		panic(err)
	}
	return env
}

// LoadEnvironmentFrom loads an environment into memory
func LoadEnvironmentFrom(f io.Reader) (Environment, error) {

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	t, err := getType(b)
	if err != nil {
		return nil, err
	}

	env := newEnvOfType(t)

	err = json.Unmarshal(b, env)
	if err != nil {
		return nil, err
	}

	return env, nil
}

func newEnvOfType(t reflect.Type) Environment {
	v := reflect.New(t.Elem())

	//	for v.Kind() == reflect.Ptr {
	//		v = v.Elem()
	//	}

	//	return v.Addr().Interface().(Environment)
	return v.Interface().(Environment)
}

func getType(b []byte) (reflect.Type, error) {
	t := &typeField{}

	if err := json.Unmarshal(b, t); err != nil {
		return nil, err
	}

	if t, ok := environmentTypes[t.Type]; ok {
		return t, nil
	}

	return nil, fmt.Errorf("could not find environment type %s; valid include %s", t.Type, strings.Join(getValidEnvironmentTypes(), ", "))
}

type typeField struct {
	Type string `json:"type"`
}

func getValidEnvironmentTypes() []string {
	names := make([]string, len(environmentTypes))

	i := 0
	for name := range environmentTypes {
		names[i] = name
		i++
	}
	return names
}
