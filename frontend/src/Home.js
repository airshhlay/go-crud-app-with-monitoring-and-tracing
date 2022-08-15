import React from "react";
import { Row, Col, Card, Button } from 'antd'
export default function Home(props) {
  return (<Row justify="center" align="middle" style={{ height: "100%" }}>
    <Col>
      <Card style={{ backgroundColor: "#EE4D2D", color: "white", fontSize: "4em"}}>
        My Shopee Favourites
      </Card>
      <Row style={{marginTop:"10px"}}>
        <Col span={10}>
          <Button onClick={props.switchToLogin} style={{width: "100%"}}>Login</Button>
        </Col>
        <Col span={4}></Col>
        <Col span={10}>
          <Button onClick={props.switchToSignup} style={{width: "100%"}}>Signup</Button>
        </Col>
      </Row>
    </Col>
  </Row>)
}