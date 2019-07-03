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

	"fmt"
	"github.com/fentec-project/gofe/innerprod/fullysec"
	"github.com/fentec-project/gofe/data"
	"github.com/pkg/errors"
	"math/big"
	"github.com/fentec-project/bn256"
)

// This is a demonstration of how can the decentralized inner-product
// functional encryption scheme implemented in GoFE library be used
// for creating anonymous heatmaps. An anonymous heatmap in our case
// is a summary of location data of many users where each user encrypts
// its data independently and send it to a server. The server is given
// keys by the users to be able to preform computations on the
// encrypted data to obtain only the heatmap of all the location data
// while the data of each individual user remains anonymous. To be more
// precise we will create
func main() {
	// the demonstration preforms all the computations in a single
	// execution while in a practical scenario computations would be
	// preformed by many clients

	// first we read the data that each user owns
	pathVectors, stations, err := readMatFromFile("london_paths.txt")
	if err != nil {
		panic(errors.Wrap(err, "error reading data"))
	}
	numClients := len(pathVectors)
	vecDim := len(pathVectors[0])
	fmt.Println("reading the data; numer of clients:", numClients)

	clients := make([]*fullysec.DMCFEClient, numClients)
	pubKeys := make([]*bn256.G1, numClients)

	// create clients and make a slice of their public values
	for i := 0; i < numClients; i++ {
		c, err := fullysec.NewDMCFEClient(i)
		if err != nil {
			panic(errors.Wrap(err, "could not instantiate fullysec"))
		}
		clients[i] = c
		pubKeys[i] = c.ClientPubKey
	}

	// based on public values of each client create private matrices T_i summing to 0
	for i := 0; i < numClients; i++ {
		err = clients[i].SetShare(pubKeys)
		if err != nil {
			panic(errors.Wrap(err, "could not create private values"))

		}
	}
	fmt.Println("clients agreed on secret keys")

	// each client encrypts his locations
	fmt.Println("simulating encryption of", numClients, "clients")
	ciphers := make([][]*bn256.G1, vecDim)
	for i := 0; i < vecDim; i++ {
		ciphers[i] = make([]*bn256.G1, numClients)
		for k := 0; k < numClients; k++ {
			label := stations[i]

			c, err := clients[k].Encrypt(pathVectors[k][i], label)
			if err != nil {
				panic(errors.Wrap(err, "could not encrypt"))
			}
			ciphers[i][k] = c
		}
	}
	fmt.Println("clients encrypted the data")


	// each client gives his key share corresponding to the vector of
	// ones; only knowing all the key shares one can decrypt the
	// sum of all locations of the clients
	keyShares := make([]data.VectorG2, numClients)
	oneVec := data.NewConstantVector(numClients, big.NewInt(1))
	for k := 0; k < numClients; k++ {
		keyShare, err := clients[k].DeriveKeyShare(oneVec)
		if err != nil {
			panic(errors.Wrap(err, "could not generate key shares"))
		}
		keyShares[k] = keyShare
	}
	fmt.Println("clients created keys for decrypting heatmap")


	heatmap := make([]*big.Int, vecDim)
	for i := 0; i < vecDim; i++ {
		label := stations[i]

		heatmap[i], err = fullysec.DMCFEDecrypt(ciphers[i], keyShares, oneVec, label, big.NewInt(int64(numClients)))
		if err != nil {
			panic(errors.Wrap(err, "could not decrypt"))
		}
	}

	fmt.Println("heatmap decrypted:")
	fmt.Println(heatmap)

	writeVecToFile("london_heatmap.txt", stations, heatmap)
}
