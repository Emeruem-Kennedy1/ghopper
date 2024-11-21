// AnalysisPage.tsx
import { Space, Spin } from "antd";
import { Node, Edge } from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import { AnalysisForm } from "../components/AnalysisForm";
import { MusicGraph } from "./MusicGraph";
import { useGraphDataQuery } from "../hooks/useGraphData";

interface GraphData {
  nodes: Node[];
  edges: Edge[];
}

const defaultData: GraphData = {
  nodes: [],
  edges: [],
};

export const AnalysisPage = () => {
  // We don't need local state anymore as React Query will handle the data
  return (
    <Space
      direction="vertical"
      size="large"
      style={{
        width: "100%",
      }}
    >
      <AnalysisForm />
      <MusicGraphContainer />
    </Space>
  );
};

// Create a separate container component to handle the graph data fetching
const MusicGraphContainer = () => {
  const { data: graphData, isLoading } = useGraphDataQuery();

  if (isLoading) {
    return (
      <div
        style={{ display: "flex", justifyContent: "center", padding: "2rem" }}
      >
        <Spin size="large" />
      </div>
    );
  }

  return <MusicGraph initialData={graphData ?? defaultData} />;
};
