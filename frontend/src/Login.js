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
          setErrorMessage("Wrong username and/or password.")
          return
        }
        if (res.errorCode && res.errorCode !== -1) {
          switch (res.errorCode) {
            case 240011:
              return setErrorMessage("Wrong username and/or password.")
            case 240012:
              return setErrorMessage("Wrong username and/or password.")
            default:
              return setErrorMessage("Unexpected error occured. Please try again later!")
          }
        }
        props.onAuthentication();
      })
      .catch((err) => {
        console.log(err)
        setErrorMessage("Unexpected error occured. Please try again later!")
      });
  };


  // const handleLogin = () => {
  //   // authenticate user
  //   login(username, password)
  //     .then(() => {
  //       props.onAuthentication();
  //       setUsername(null);
  //       setPassword(null);
  //     })
  //     .catch((err) => {
  //       if (err.response) {
  //         if (err.response.status === 400) {
  //           showFailureMsg("User already exists");
  //         }
  //       } else {
  //         showFailureMsg(err)
  //       }
  //     });
  // };

  // const handleSignup = () => {
  //   if (!username || username.length > 8 || username.length < 3) {
  //     showFailureMsg('Username must be 3 to 8 characters long');
  //     return;
  //   }
  //   if (!password || password.length < 8 || password.length >= 15) {
  //     showFailureMsg('Password must be 8 to 15 characters long');
  //     return;
  //   }
  //   signup(username, password)
  //     .then(() => {
  //       showSuccessMsg("Signup successful. You can now login.")
  //       setUsername(null);
  //       setPassword(null);
  //     })
  //     .catch((err) => {
  //       if (err.response) {
  //         showFailureMsg(err.response.data.errorMsg);
  //       }
  //     });
  // };

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