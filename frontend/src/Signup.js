import React from "react";
import { message, Card, Button, Row, Col, Input, Form } from "antd";
import { WarningOutlined } from "@ant-design/icons";
import { signup } from "./api/auth";
export default function Signup(props) {
  const [errorMessage, setErrorMessage] = React.useState("")

  const onFinish = (values) => {
    username = values.username
    password = values.password
    if (!username || !password) {
      setErrorMessage("Username and/or password can not be blank.")
    }
    if (username.length > 15 || username.length < 3) {
      setErrorMessage("Username must be 3 to 15s characters long!")
    } else if (password.length < 8 || password.length > 20) {
      setErrorMessage("Password should be 8 to 20 characters long!")
    }
    signup(values.username, values.password)
      .then((res) => {
        if (!res.errorCode) {
          setErrorMessage("Wrong username and/or password.")
          return
        }
        if (res.errorCode && res.errorCode !== -1) {
          switch (res.errorCode) {
            case 240011:
              return setErrorMessage("User already exists. Please login!")
            case 240013:
              return setErrorMessage("Wrong username and/or password.")
            default:
              return setErrorMessage("Unexpected error occured. Please try again later!")
          }
        }
        setErrorMessage("Successfully signed up. Please login.")
      })
      .catch((err) => {
        console.log(err)
        setErrorMessage("Unexpected error occured. Please try again later!")
      });
  };

  return (
    <Row type="flex" justify="center" align="middle" style={{
      height: "100vh",
      margin: "auto"
    }}>
      <Col style={{
        margin: "auto"
      }}>
          <Form
            name="basic"
            labelCol={{ span: 8 }}
            wrapperCol={{ span: 16 }}
            initialValues={{ remember: true }}
            onFinish={onFinish}
            autoComplete="off"
            scrollToFirstError
          >
            <Form.Item
              label="Username"
              name="username"
              rules={[{ required: true, message: 'Please input your username!' }]}
            >
              <Input />
            </Form.Item>

            <Form.Item
              label="Password"
              name="password"
              rules={[{ required: true, message: 'Please input your password!' }]}
            >
              <Input.Password />
            </Form.Item>

            {errorMessage !== "" && (
              <p className="errorMessage">
                <WarningOutlined />
                &nbsp;{errorMessage}
              </p>
            )}


            <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
              <Button type="primary" htmlType="submit">
                Signup
              </Button>
            </Form.Item>
          </Form>
      </Col>
    </Row>
  );
}