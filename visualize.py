import networkx as nx
from pyvis.network import Network
from random import shuffle
import pickle

def main(username='balajis', limit=100):
    G = construct_graph()
    print('calculating pagerank')
    pagerank = calculate_page_rank(G)
    subgraph_usernames = list(G.successors(username))
    subgraph_usernames.extend(list(G.predecessors(username)))
    shuffle(subgraph_usernames)
    subgraph_usernames = subgraph_usernames[:limit-1]
    subgraph_usernames.append(username)
    print('%s nodes in subgraph' % len(subgraph_usernames))
    subgraph = G.subgraph(subgraph_usernames)
    print('subgraph generated')
    # reconstruct subgraph with weights
    subgraph_with_weights = nx.DiGraph()    
    for node in subgraph_usernames:
        subgraph_with_weights.add_node(node, size=float(pagerank[node])*500)
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