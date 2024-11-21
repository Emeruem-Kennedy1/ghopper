import { useQuery } from "@tanstack/react-query";
import { GraphData } from "../types/songSearch";

const startX = 250;
const startY = 100;

const defaultData: GraphData = {
  nodes: [
    {
      id: "1",
      data: { label: "Song: Blinding Lights \n Genres: Pop, HipHop" },
      position: { x: startX, y: startY },
      //   type: "default",
    },
    {
      id: "2",
      data: { label: "Song: Levitating \n Genres: Pop, Dance" },
      position: { x: startX + 200, y: startY + 100 },
      type: "default",
    },
    {
      id: "3",
      data: { label: "Song: Butter \n Genres: Pop, HipHop" },
      position: { x: startX + 400, y: startY + 200 },
      type: "default",
    },
    {
      id: "4",
      data: { label: "Song: Stay \n Genres: Pop, HipHop" },
      position: { x: startX + 600, y: startY + 300 },
      type: "default",
    },
  ],
  edges: [
    { id: "e1-2", source: "1", target: "2" },
    { id: "e3-1", source: "3", target: "1" },
    { id: "e4-1", source: "4", target: "2" },
  ],
};;

export const useGraphDataQuery = () => {
  return useQuery({
    queryKey: ["graphData"],
    queryFn: () => Promise.resolve(defaultData), // Initially returns empty data
    enabled: false, // Don't fetch automatically
  });
};
