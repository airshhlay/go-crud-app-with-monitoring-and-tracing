import axios from "axios";
const endpoint = process.env.REACT_APP_ENDPOINT
const LOGIN = process.env.LOGIN ? process.env.LOGIN : "/api/user/login"
const SIGNUP = process.env.SIGNUP ? process.env.SIGNUP : "/api/user/signup"
axios.defaults.withCredentials = true

const login = (username, password) => {
  return axios
    .post(`${endpoint}${LOGIN}`, { username, password }, {withCredentials: true})
    .then((res) => {
      return res.data;
    })
    .catch((err) => {
      throw err;
    });
};

const signup = (username, password) => {
  return axios
  .post(`${endpoint}${SIGNUP}`, {username, password}, {withCredentials: true})
    .then((res) => {
      return res.data;
    })
    .catch((err) => {
      throw err;
    });
}

export { login, signup};