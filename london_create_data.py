import csv
import random
import networkx as nx
import matplotlib.pyplot as plt


def random_walk(G):
    walk = []
    in_st = random.choice(list(G.nodes()))
    walk.append(in_st)
    length = random.randint(2, int(nx.diameter(G) / 2))
    for _ in range(length):
        neigh = [x for x in G[walk[-1]] if x not in walk]
        if len(neigh) == 0:
            return walk
        next_st = random.choice(neigh)
        walk.append(next_st)
    return walk


o = open("london_stations.csv")
reader = list(csv.reader(o, delimiter=','))
positions = {}
id_to_station = {}
for row in reader:
    if row[0] == "id":
        continue
    positions[row[3]] = [float(row[2]), 2 * float(row[1])]
    id_to_station[row[0]] = row[3]
o.close()

o2 = open("london_connections.csv")
reader2 = list(csv.reader(o2, delimiter=','))
G = nx.Graph()
for row in reader2:
    if row[0] == "station1":
        continue
    if G.has_edge(id_to_station[row[0]], id_to_station[row[1]]):
        continue
    G.add_edge(id_to_station[row[0]], id_to_station[row[1]])

stations = list(G.nodes())
num_clients = 100
vectors = []
for i in range(num_clients):
    flip = random.random()
    if flip > 0.7:
        in_st = random.choice(stations)
        out_st = random.choice(stations)
        path = nx.shortest_path(G, in_st, out_st)
    else:
        path = random_walk(G)
    vec = []
    for i in range(len(stations)):
        if stations[i] in path:
            vec.append(1)
        else:
            vec.append(0)
    vectors.append(vec)

w = open("london_paths.txt", "w")
w.write(";".join(stations) + "\n")
for i in range(num_clients):
    w.write(";".join([str(x) for x in vectors[i]]) + "\n")
w.close()

nx.draw(G, node_size=30, pos=positions, node_color=vectors[0], cmap="YlGnBu",
               vmin=-2, vmax=2)
plt.savefig('one_user.png', dpi=1000)
