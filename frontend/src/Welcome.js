import React from "react";
import { message, Card, Button, Row, Col, Input } from "antd";
import { login, signup } from "./api/auth";
const { Password } = Input;
export default function Welcome(props) {
  const [username, setUsername] = React.useState("");
  const [password, setPassword] = React.useState("");

  const showSuccessMsg = (msg) => {
    message.success({
      content: msg,
      style: { marginTop: "5vh"},
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

  const handleLogin = () => {
    // authenticate user
    login(username, password)
      .then(() => {
        props.onAuthentication();
        setUsername(null);
        setPassword(null);
      })
      .catch((err) => {
        if (err.response) {
          if (err.response.status === 400) {
            showFailureMsg("User already exists");
          }
        } else {
            showFailureMsg(err)
        }
      });
  };

  const handleSignup = () => {
    if (!username || username.length > 8 || username.length < 3) {
        showFailureMsg('Username must be 3 to 8 characters long');
        return;
    }
    if (!password || password.length < 8 || password.length >= 15) {
        showFailureMsg('Password must be 8 to 15 characters long');
        return;
    }
    signup(username, password)
      .then(() => {
        showSuccessMsg("Signup successful. You can now login.")
        setUsername(null);
        setPassword(null);
      })
      .catch((err) => {
        if (err.response) {
            showFailureMsg(err.response.data.errorMsg);
        }
      });
  };

  return (
        <div style={{
            display: "flex",
            height:"100%",
            flexDirection:"column",
            justifyContent:"center",
            alignContent:"middle"
        }}>
        <div>
            <h1 style={{color:"#EE4D2D"}}>My Shopee Favourites</h1>
        </div>
        <div>
            <Input placeholder="Username" onChange={(e) => setUsername(e.target.value)} />
        </div>
        <div>
            <Password placeholder="Password" onChange={(e) => setPassword(e.target.value)} />
        </div>
        <div>
            <Button onClick={handleLogin}>Login</Button>
        </div>
        <div>
            <Button onClick={handleSignup}>Signup</Button>
        </div>
        </div>
  );
}