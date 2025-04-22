/*
gocu is a curl copycat, a CLI HTTP client focused on simplicity and ease of use
Copyright (C) 2025  Andrés C.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package cache

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const separator = "​" // zwsp 200B
const pairSep = separator + separator
const nameValPairSep = "=" + separator

var varsCache map[string]string

func GetVariable(name string) (string, error) {
	val, exists := Variables()[name]

	if !exists {
		return "", fmt.Errorf("Variable %q does not exist", name)
	}

	return val, nil
}

func AddVariable(name string, value string) error {
	vars := Variables()
	val, exists := vars[name]

	if exists {
		return fmt.Errorf("Variable named %q (value=%s) already exists. Use modify command to override.", name, val)
	}

	setVariable(name, value)

	log.Printf("Added variable %s=%s", name, value)

	return nil
}

func ModifyVariable(name string, value string) error {
	oldVal, err := GetVariable(name)

	if err != nil {
		return err
	}

	setVariable(name, value)

	log.Printf("Modified value of variable %q (%s => %s)", name, oldVal, value)

	return nil
}

func RemoveVariable(name string) error {
	vars := Variables()
	val, exists := vars[name]

	if !exists {
		return fmt.Errorf("Variable named %q does not exist", name)
	}

	delete(Variables(), name)
	save()

	log.Printf("Deleted variable %s=%s", name, val)

	return nil
}

func RemoveAllVariables() error {
	vars := Variables()

	for v := range vars {
		err := RemoveVariable(v)

		if err != nil {
			return err
		}
	}

	return nil
}

func Variables() map[string]string {
	if varsCache == nil {
		loadCache()
	}

	return varsCache
}

func create() {
	p := path()
	_, err := os.Create(p)

	if err != nil {
		log.Fatalf("Couldn't create variables cache: %s", err.Error())
	}

	log.Printf("Cache created (%s)\n", p)
}

func exists() bool {
	p := path()
	_, err := os.Stat(p)

	return !errors.Is(err, os.ErrNotExist)
}

func path() string {
	cacheDir, err := os.UserCacheDir()

	if err != nil {
		log.Fatalf("Couldn't retrieve variables cache path: %s", err.Error())
	}

	return filepath.Join(cacheDir, "gocu.cache")
}

func save() {
	rwMask := 0770
	pairs := make([]string, 0)

	for name, val := range Variables() {
		pair := fmt.Sprintf("%s%s%s", name, nameValPairSep, val)
		pairs = append(pairs, pair)
	}

	res := strings.Join(pairs, pairSep)
	err := os.WriteFile(path(), []byte(res), os.FileMode(rwMask))

	if err != nil {
		log.Fatalf("Couldn't save variables cache: %s", err.Error())
	}
}

func setVariable(name string, value string) {
	Variables()[name] = value
	save()
}

func loadCache() {
	if !exists() {
		create()
	}

	c := cacheFileContent()

	if c != "" {
		varsCache = parseCacheFileContent(c)
	} else {
		varsCache = make(map[string]string)
	}
}

func cacheFileContent() string {
	p := path()
	c, err := os.ReadFile(p)

	if err != nil {
		log.Fatalf("Couldn't read variables cache: %s", err.Error())
	}

	return string(c)
}

func parseCacheFileContent(content string) map[string]string {
	res := make(map[string]string)

	pairs := strings.SplitSeq(content, pairSep)

	for pair := range pairs {
		aux := strings.Split(pair, nameValPairSep)

		name := aux[0]
		val := aux[1]

		res[name] = val
	}

	return res
}
