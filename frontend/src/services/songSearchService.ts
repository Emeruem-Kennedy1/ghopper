import axios from "axios";
import { ArtistSearchData } from "../types";
import { getToken } from "../utils/auth";
import { GraphData, SongSearchResponseData } from "../types/songSearch";
import dagre from "dagre";

export async function getSongs(data: ArtistSearchData): Promise<GraphData> {
  const token = getToken();
  const response = await axios.post<SongSearchResponseData>(
    "/api/api/search",
    data,
    {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }
  );

  return transformToGraphData(response.data);
}

export function transformToGraphData(
  responseData: SongSearchResponseData
): GraphData {
  // Create nodes without positions first
  const initialNodes = Object.entries(responseData.nodes).map(([id, song]) => {
    const artistNames = song.artists.map((artist) => artist.name).join(", ");

    return {
      id: id,
      data: {
        label: `Song: ${
          song.title
        }\nArtist: ${artistNames}\nGenres: ${song.genres.join(", ")}`,
      },
      position: { x: 0, y: 0 }, // Initial position will be updated by dagre
      type: "default",
    };
  });

  // Create edges with deduplication
  const edgeSet = new Set<string>();
  const initialEdges: GraphData["edges"] = [];

  Object.entries(responseData.adjacencyList).forEach(([sourceId, targets]) => {
    Object.keys(targets).forEach((targetId) => {
      // Create a unique edge identifier that's the same regardless of direction
      const edgeNodes = [sourceId, targetId].sort();
      const edgeId = `e${edgeNodes[0]}-${edgeNodes[1]}`;

      // Only add the edge if we haven't seen it before
      if (!edgeSet.has(edgeId)) {
        edgeSet.add(edgeId);
        initialEdges.push({
          id: edgeId,
          source: edgeNodes[0],
          target: edgeNodes[1],
        });
      }
    });
  });

  // Set up the dagre graph
  const g = new dagre.graphlib.Graph({
    directed: false, // Specify that the graph is undirected
  });

  g.setGraph({
    rankdir: "TB",
    nodesep: 20,
    ranksep: 200,
    marginx: 50,
    marginy: 50,
  });

  g.setDefaultEdgeLabel(() => ({}));

  // Add nodes to the dagre graph
  initialNodes.forEach((node) => {
    g.setNode(node.id, {
      width: 250,
      height: 80,
    });
  });

  // Add edges to the dagre graph
  initialEdges.forEach((edge) => {
    g.setEdge(edge.source, edge.target);
  });

  // Calculate the layout
  dagre.layout(g);

  // Update node positions based on dagre calculations
  const nodes = initialNodes.map((node) => {
    const nodeWithPosition = g.node(node.id);
    return {
      ...node,
      position: {
        x: nodeWithPosition.x - 125,
        y: nodeWithPosition.y - 40,
      },
    };
  });

  return {
    nodes,
    edges: initialEdges,
  };
}
