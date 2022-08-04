import { Tabs, Button, Row, Col } from "antd";
import ItemList from "./List";
const { TabPane } = Tabs;

export default function Navigation(props) {
  return (
    <Row type="flex" justify="center" align="top" style={{ height: "100%" }}>
      <Col span={2} />
      <Col span={20}>
        <Tabs defaultActiveKey="1" centered>
          <TabPane tab="My Favourites" key="2">
            <ItemList onUnauthentication={props.onUnauthentication} />
          </TabPane>
        </Tabs>
      </Col>
      <Col span={2}>
        <div>
          <Button onClick={props.handleLogout}>Logout</Button>
        </div>
      </Col>
    </Row>
  );
}