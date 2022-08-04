import React from "react";
import logo from "./logo.svg";
import "./App.css";
import Navigation from "./Navigation";
import Welcome from "./Welcome";

function App() {
  const [authenticated, setAuthenticated] = React.useState(false);
  const onAuthentication = () => {
    setAuthenticated(true);
  };
  const onUnauthentication = () => {
    setAuthenticated(false);
  };

  const handleLogout = () => {
    setAuthenticated(false);
  };
  return (
    <div className="App">
      {!authenticated && <Welcome onAuthentication={onAuthentication} />}
      {authenticated && (
        <Navigation
          onUnauthentication={onUnauthentication}
          handleLogout={handleLogout}
        />
      )}
    </div>
  );
}

export default App;