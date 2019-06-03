import csv
import networkx as nx
import matplotlib.pyplot as plt

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
G = nx. Graph()
for row in reader2:
    if row[0] == "station1":
        continue
    if G.has_edge(id_to_station[row[0]], id_to_station[row[1]]):
        continue
    G.add_edge(id_to_station[row[0]], id_to_station[row[1]])

o = open("london_heatmap.txt")
stations = o.readline().strip().split(";")
heatmap = [int(x) for x in o.readline().split(";")]
o.close()
heatmap_dict = {stations[i]: heatmap[i] for i in range(len(stations))}
colors = [heatmap_dict[x] for x in G.nodes()]
print(heatmap_dict)

nx.draw(G, node_size=30, pos=positions, node_color=colors, cmap="YlGnBu",
               vmin=-10, vmax=max(colors))
plt.savefig("heatmap.png", dpi=1000)
