import axios from "axios";
const endpoint = process.env.NODE_ENV === "development" ? process.env.REACT_APP_ENDPOINT_LOCAL : process.env.REACT_APP_ENDPOINT_PROD;

const login = (username, password) => {
  return axios
    .post(`${endpoint}/user/login`, { username: username, password: password })
    .then((res) => {
      return res.data;
    })
    .catch((err) => {
      throw err;
    });
};

const signup = (username, password) => {
  return axios
  .post(`${endpoint}/user/signup`, { username: username, password: password })
    .then((res) => {
      return res.data;
    })
    .catch((err) => {
      throw err;
    });
}

export { login, signup};