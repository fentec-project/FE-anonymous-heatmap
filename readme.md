# Creating privacy preserving heatmap with functional encryption - example

This repository demonstrates how the [GoFE](https://github.com/fentec-project/gofe) functional encryption (FE)
library can be used for creating heatmap from location data of users in a way that the data of each individual
is anonymous and encrypted. The server is able to perform computations on encrypted data together with keys given
by all the users decrypting only the sum of all the locations, i.e. creating a heatmap.

Particularly, we simulate users of London Underground giving encrypted location information to the server
allowing the server to create only the sum of all their locations while the individual's locations remain
anonymous to the server as well as to the other users. With such data the server gains information about
the usage of each station making it able to detect increased traffic, congestions, etc., while preserving the
privacy of the users.

### Input data
In this repository one can find data about the stations of London Underground (`london_stations.csv`) and 
train lines connecting them (`london_connections.csv`). The data is collected with Python3 script
`london_create_data.py`, which additionally creates location information about users. The latter is created
in a random way, each user randomly chooses a starting station and by a mix of shortest paths and random choices
travels to a final station. For example, data of one user can be visualized in the following image showing the path
traveled by one user:

![Alt text](one_user.png?raw=true "Visualization of data of a user.")

To be more precise, the data of each user is a vector of zeros and ones, each input of the vector indicating
if the user has traveled via the corresponding station. The data of all such vectors is saved in `london_paths.txt`.
Note: this is a private information of every user which in a practical scenario would not be collected.


## Running the example

1. Build the example by running:
````bash
$ go get github.com/fentec-project/FE-anonymous-heatmap
````
2. This will produce the `FE-anonymous-heatmap` executable in your `$GOPATH/bin`.
If you have `$GOPATH/bin` in your `PATH` environment variable, you
will be able to run the example by running command `FE-anonymous-heatmap` from the
root of this repository.

Otherwise just call:
```bash
$ $GOPATH/bin/FE-anonymous-heatmap
```
The code will do the following:
* It will read the data from `london_paths.txt` and simulate the interaction between clients and the server.
Each client has a vector describing his location data, which he does not wants to disclose. The goal of the
server is to obtain the heatmap of all the locations.
* Clients first need to agree on a common shared secret. This is needed since there is no central server which
can assign secret keys to the users; if such a server existed the scheme would not be anonymous any more. Creating
a common secret is done by each user publishing a public key. Then each user can collect the public keys of all the
others and creates a shared secret with his private key.
* Each user then encrypts his vector of locations with his secret.
* Each user creates a partial key which can be used only for the decryption of the sum of all the the location vectors.
* Finally, the server collects partial keys of all the users. With all the keys it is able to decrypt the sum of all
location vectors giving it the data for the heatmap. No individual data can be decrypted with these keys.

Running the above script should output the following:
````
reading the data; numer of clients: 100
clients agreed on secret keys
simulating encryption of 100 clients
clients encrypted the data
clients created keys for decrypting heatmap
heatmap decrypted:
[1 2 2 7 1 2 1 3 4 5 5 4 4 13 2 3 5 1 0 12 3 4 1 4 4 0 5 4 9 3 5 3 4 1 8 2 5 2 3 0 6 4 3 5 6 1 2 0 12 4 7 1 6 1 3 6 3 8 4 2 2 1 1 2 4 2 2 3 1 2 4 5 4 1 6 1 2 3 1 3 3 11 1 11 2 7 1 1 0 0 2 2 5 1 2 5 5 2 3 0 2 2 5 2 0 3 3 1 2 1 1 3 6 5 1 0 1 5 1 3 2 6 7 8 1 0 7 2 0 5 2 1 0 5 0 3 4 10 0 1 1 16 3 6 3 4 3 3 6 3 13 2 2 2 0 3 7 12 1 8 0 0 6 5 6 1 6 10 6 1 4 2 6 1 2 2 4 10 5 3 4 11 2 4 2 1 6 3 3 10 2 5 0 1 5 1 15 2 4 12 7 1 2 2 6 2 9 2 5 5 7 3 4 1 9 7 5 1 2 1 0 1 5 4 4 5 3 3 4 5 0 4 1 1 1 1 0 4 4 5 4 2 1 0 2 4 1 1 3 3 9 6 5 0 13 3 1 3 3 1 2 3 2 7 4 4 9 3 4 11 9 5 4 1 0 6 2 2 3 7 3 22 5 3 1 3 1 2 4 4 2 5 3 1 2 1 1 4 15 4 3 4]
````
Explanation: the code simulates the above described process for 100 users. There are 302 London Underground stations
included in this example, hence the output is a 302-dimensional vector describing how many users traveled through
each station.

### Visualization
The obtained heatmap vector can be visualized on the network of London Underground with colors indicating
the traffic through each station. This can be achieved by running a provided Python3 script. First install
the dependencies:
````bash
$ pip3 install networkx matplotlib
```` 
Navigate to the root of this repository and run the visualization script:
````bash
$ cd $GOPATH/src/github.com/fentec-project/FE-anonymous-heatmap
$ python3 london_visualize.py
````
The image is now saved in `heatmap.png`. It should look like:
![Alt text](heatmap.png?raw=true "Visualization of data of a user.")

### Generating training data
If one wishes to re-create the random data, this can be cone by running the following
Python3 script:
````bash
$ python3 london_create_data.py
````