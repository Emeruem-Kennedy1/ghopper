import { ReactFlow, Node, Edge, Controls, ConnectionMode } from "@xyflow/react";
import { theme } from "antd";
import "@xyflow/react/dist/style.css";

interface GraphData {
  nodes: Node[];
  edges: Edge[];
}

interface MusicGraphProps {
  initialData: GraphData;
}

export const MusicGraph = ({ initialData }: MusicGraphProps) => {
  const { token } = theme.useToken();

  const nodeStyle = {
    style: {
      background: token.colorBgContainer,
      color: token.colorText,
      border: `1px solid ${token.colorBorder}`,
      borderRadius: token.borderRadius,
      padding: token.padding,
    },
  };

  // Update nodes with theme-aware styling
  const themedNodes = initialData.nodes.map((node) => ({
    ...node,
    ...nodeStyle,
  }));

//   const edgeStyle = {
//     style: {
//       stroke: token.colorBorder,
//     },
//   };

  // Update edges with theme-aware styling
  const themedEdges = initialData.edges.map((edge) => ({
    ...edge,
    // ...edgeStyle,
  }));

  return (
    <div
      style={{
        height: "500px",
        border: `1px solid ${token.colorBorder}`,
        borderRadius: token.borderRadius,
        overflow: "hidden",
      }}
    >
      <ReactFlow
        nodes={themedNodes}
        edges={themedEdges}
        connectionMode={ConnectionMode.Strict}

        fitView
        colorMode="dark"
      >
        <Controls orientation="horizontal" />
      </ReactFlow>
    </div>
  );
};
