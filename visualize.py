import networkx as nx
from pyvis.network import Network
from random import shuffle
import pickle

def main(username='balajis', limit=40):
    G = construct_graph()
    print('calculating pagerank')
    pagerank = calculate_page_rank(G)
    neighbors = list(G.successors(username))
    neighbors.extend(list(G.predecessors(username)))
    
    # sort by highest PR
    highest = sorted([(neighbor, pagerank[neighbor]) for neighbor in neighbors], key=lambda t: t[1], reverse=True)[:limit-1]
    highest_nodes = [t[0] for t in highest]
    highest_nodes.append(username)
    print('%s nodes in subgraph' % len(highest_nodes))
    subgraph = G.subgraph(highest_nodes)
    print('subgraph generated')
    # reconstruct subgraph with weights
    subgraph_with_weights = nx.DiGraph()    
    for node in highest_nodes:
        print(pagerank[node])
        subgraph_with_weights.add_node(node, size=float(pagerank[node])*1000)
    for edge in subgraph.edges():
        subgraph_with_weights.add_edge(*edge)
    print('subgraph with weights generated')
    net = Network(notebook=True, height="2000px", width="2000px")
    net.from_nx(subgraph_with_weights)
    print('network generated')
    net.show("example.html")

def calculate_page_rank(G, filename='pagerank.pickle'):
    try:
        with open(filename, 'rb') as f:
            return pickle.load(f)
    except:
        print('building page rank from scratch')
    pagerank = nx.pagerank(G)
    with open(filename, 'wb') as f:
        pickle.dump(pagerank, f)
    return pagerank

def construct_graph(filename='graph.pickle'):
    try:
        with open(filename, 'rb') as f:
            return pickle.load(f)
    except:
        print('building full graph from scratch')

    g = nx.DiGraph()
    with open('edges.txt', 'r') as f:
        for line in f:
            # followed, follower
            a, b = line.split()
            g.add_edge(b, a)
    with open(filename, 'wb') as f:
        pickle.dump(g, f)
    return g


if __name__ == '__main__':
    main()