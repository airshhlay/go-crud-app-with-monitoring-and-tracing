import React from "react";
import {
  Select,
  Input,
  Popover,
  Tag,
  Row,
  Col,
  Button,
  Card,
  message,
} from "antd";
import { DeleteFilled, EditFilled, EnterOutlined } from "@ant-design/icons";
import { getAll, deleteItem } from "./api/api";

const { Option } = Select;
export default function ItemList(props) {
  const [data, setData] = React.useState();

  // fetch user's favourite items upon page load
  React.useEffect(() => {
    onRefresh();
  }, []);

  // shows a success message to user upon successful adding of item
  const showSuccessMsg = (msg) => {
    message.success({
      content: msg,
      style: { marginTop: "5vh" },
      duration: 1,
    });
  };

  // refreshes user's favourites list
  const onRefresh = () => {
    getItemList()
      .then((res) => {
        setData(res);
      })
      .catch((err) => {
        if (err.response && err.response.status === 401 ) {
            showFailureMsg("Session timeout. You have been logged out");
            props.onUnauthentication();
            return;
        }
        showFailureMsg("Unexpected error occured");
      });
  };

  const showFailureMsg = (msg) => {
    message.error({
      content: msg,
      style: { marginTop: "5vh" },
      duration: 1,
    });
  };

  const onDelete = (tag, index) => {
    deleteItem(tag, index)
      .then((res) => {
        showSuccessMsg("it's gone");
        onRefresh();
      })
      .catch((err) => {
        if (err.response && err.response.status === 401) {
            showFailureMsg("Session timeout. You have been logged out");
            props.onUnauthentication();
            return;
        }
        showFailureMsg("Unexpected error occured");
      });
  };

  return (
    <Row type="flex" justify="center" style={{ height: "100%" }}>
      <Col>
        <Button style={{ marginLeft: "20vw" }} onClick={onRefresh}>
          Refresh
        </Button>
        <div style={{ overflowY: "auto", height: "100%" }}>
          {data &&
            data.map((item) => {
                return (
                  <Card key={`${item.itemId}`}>
                    <Row>
                      <Col span={19}>
                        <p>{item.name}</p>
                      </Col>
                      <Col span={2}>
                        <Button
                          style={{ backgroundColor: "lightblue" }}
                          onClick={() => onDelete(item.itemId, item.shopId)}
                        >
                          <DeleteFilled />
                        </Button>
                      </Col>
                    </Row>
                  </Card>
                );
            })}
        </div>
      </Col>
    </Row>
  );
}