/*
Copyright (c) 2018 XLAB d.o.o

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bufio"
	"math/big"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/fentec-project/gofe/data"
)

// readMatFromFile reads matrix elements from the provided file
// and gives a matrix and the names of the rows
func readMatFromFile(path string) (data.Matrix, []string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error reading matrix from file")
	}

	scanner := bufio.NewScanner(file)
	vecs := make([]data.Vector, 0)

	first := true
	var names []string
	for scanner.Scan() {
		line := scanner.Text()
		values := strings.Split(line, ";")
		if first == true {
			names = values
			first = false
		} else {
			v := make(data.Vector, len(values))
			for i, n := range values {
				v[i], _ = new(big.Int).SetString(n, 10)
			}
			vecs = append(vecs, v)
		}
	}

	mat, err := data.NewMatrix(vecs)
	return mat, names, err
}

// writeVecToFile takes a vector and the names of its inputs and writes
// a file with the first line being names of the inputs and the second line
// its corresponding values in the vector
func writeVecToFile(path string, names []string, vec data.Vector) (error) {
	// open output file
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()
	// make a write buffer
	w := bufio.NewWriter(file)

	// make a write string
	str := ""
	for i, e := range names {
		str += e
		if i < len(vec) - 1 {
			str += ";"
		} else {
			str += "\n"
		}
	}

	for i, e := range vec {
		str += e.String()
		if i < len(vec) - 1 {
			str += ";"
		} else {
			str += "\n"
		}
	}

	// write
	if _, err := w.Write([]byte(str)); err != nil {
		panic(err)
	}

	if err = w.Flush(); err != nil {
		panic(err)
	}
	return err
}