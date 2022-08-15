import React from "react";
import { message, Card, Button, Row, Col, Input, Form } from "antd";
import { WarningOutlined } from "@ant-design/icons";
import { login, signup } from "./api/auth";
export default function Login(props) {
  const [errorMessage, setErrorMessage] = React.useState("")
  const showSuccessMsg = (msg) => {
    message.success({
      content: msg,
      style: { marginTop: "5vh" },
      duration: 1,
    });
  };

  const showFailureMsg = (msg) => {
    console.log(msg)
    message.error({
      content: msg,
      style: { marginTop: "5vh" },
      duration: 1,
    });
  };

  const onFinish = (values) => {
    login(values.username, values.password)
      .then((res) => {
        if (!res.errorCode) {
          return setErrorMessage("Unexpected error occured. Please try again later!")
        }
        if (res.errorCode && res.errorCode !== -1) {
          switch (res.errorCode) {
            case 240011:
              return setErrorMessage("Wrong username and/or password.")
            case 240012:
              return setErrorMessage("Wrong username and/or password.")
            case 240013:
              return setErrorMessage("Wrong username and/or password.")
            default:
              return setErrorMessage("Unexpected error occured. Please try again later!")
          }
        }
        props.onAuthentication();
      })
      .catch((err) => {
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
                Login
              </Button>
            </Form.Item>
          </Form>
      </Col>
    </Row>
  );
}