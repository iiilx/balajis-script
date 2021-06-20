import networkx as nx
import matplotlib.pyplot as plt

OUTPUT_FILE_NAME = "bitclout_social_graph.png"

def main():
    g = construct_graph()
    nx.draw(g)
    plt.savefig(OUTPUT_FILE_NAME)
    print('Saved graph to %s' % OUTPUT_FILE_NAME)

def construct_graph():
    g = nx.Graph()
    with open('edges.txt', 'r') as f:
        for line in f:
            a, b = line.split()
            g.add_edge(a, b)
    return g


if __name__ == '__main__':
    main()