import { Button } from "antd";
import { LeftOutlined, RightOutlined } from "@ant-design/icons";

interface CarouselArrowProps {
  type: "prev" | "next";
  onClick?: () => void;
}

export const CarouselArrow = ({ type, onClick }: CarouselArrowProps) => (
  <Button
    type="text"
    icon={type === "prev" ? <LeftOutlined /> : <RightOutlined />}
    onClick={onClick}
    style={{
      position: "absolute",
      top: "45%",
      transform: "translateY(-50%)",
      zIndex: 2,
      [type === "prev" ? "left" : "right"]: 0,
      background: "rgba(255, 255, 255, 0.1)",
      border: "0px solid #d9d9d9",
      borderRadius: "50%",
      width: "32px",
      height: "32px",
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
      padding: 0,
      fontSize: "12px",
    }}
  />
);
