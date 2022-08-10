import React from "react";
import logo from "./logo.svg";
import "./App.css";
import Navigation from "./Navigation";
import Login from "./Login";
import Home from "./Home"
import Signup from "./Signup"

function App() {
  const [authenticated, setAuthenticated] = React.useState(localStorage.getItem("authenticated") ? localStorage.getItem("authenticated") : "");
  const [page, setPage] = React.useState("")
  const onAuthentication = () => {
    setAuthenticated("yes");
    localStorage.setItem("authenticated", "yes")
  };
  const onUnauthentication = () => {
    setAuthenticated("");
  };

  const handleLogout = () => {
    setAuthenticated("");
    localStorage.removeItem("authenticated")
  };

  const switchToLogin = () => {
    setPage("login")
  }

  const switchToSignup = () => {
    setPage("signup")
  }

  const isAuthenticated = () => {
    return authenticated === "yes"
  }

  return (
    <div className="App">
      {!isAuthenticated() && page === "" && <Home switchToLogin={switchToLogin} switchToSignup={switchToSignup} />}
      {!isAuthenticated() && page === "login" && <Login onAuthentication={onAuthentication} />}
      {!isAuthenticated() && page === "signup" && <Signup onAuthentication={onAuthentication} />}
      {isAuthenticated() && (
        <Navigation
          onUnauthentication={onUnauthentication}
          handleLogout={handleLogout}
        />
      )}
    </div>
  );
}

export default App;