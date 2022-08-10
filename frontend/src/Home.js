import React from "react";
import {Row, Col, Card, Button} from 'antd'
export default function Home(props) {
  return <Row>
    <Col>
      <Card>
        My Shopee Favourites
      </Card>
    </Col>
    <Col>
    <Button onClick={props.switchToLogin}>Login</Button>
    <Button onClick={props.switchToSignup}>Signup</Button>
    </Col>
  </Row>
}