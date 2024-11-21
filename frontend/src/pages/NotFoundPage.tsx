import { useNavigate } from "react-router-dom";
import { Button, Result } from "antd";

const NotFoundPage = () => {
  const navigate = useNavigate();

  return (
    <div
      style={{
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
      }}
    >
      <Result
        status="404"
        title="404"
        subTitle="Sorry, the page you visited does not exist."
        extra={
          <div
            style={{ display: "flex", gap: "12px", justifyContent: "center" }}
          >
            <Button type="primary" onClick={() => navigate("/")}>
              Back Home
            </Button>
            <Button onClick={() => navigate(-1)}>Go Back</Button>
          </div>
        }
      />
    </div>
  );
};

export default NotFoundPage;
