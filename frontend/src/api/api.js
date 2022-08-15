import axios from "axios";
const endpoint = process.env.REACT_APP_ENDPOINT
// process.env.NODE_ENV === "development" ? process.env.REACT_APP_ENDPOINT_LOCAL : process.env.REACT_APP_ENDPOINT_PROD;
const GET_LIST = process.env.REACT_APP_GET_LIST ? process.env.REACT_APP_GET_LIST : "/api/item/get/list"
const ADD_ITEM = process.env.REACT_APP_ADD_ITEM ? process.env.REACT_APP_ADD_ITEM : "/api/item/add/fav"
const DELETE_ITEM = process.env.REACT_APP_DELETE_ITEM ? process.env.REACT_APP_DELETE_ITEM : "/api/item/delete/fav"
axios.defaults.withCredentials = true
axios.interceptors.request.use(
  (config) => {
    config.withCredentials = true
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)
const getItemList = (page) => {
  return axios
    .get(`${endpoint}${GET_LIST}`, {withCredentials: true, params: {page}})
    .then((res) => {
      return res.data;
    })
    .catch((err) => {
      throw err;
    });
};

const submitItem = (itemID, shopID) => {
  return axios
    .post(
      `${endpoint}${ADD_ITEM}`, 
      { itemID, shopID },
      {withCredentials: true},
    )
    .then((res) => {
      return res.data;
    })
    .catch((err) => {
      throw err;
    });
};

const deleteItem = (itemID, shopID) => {
  console.log(`deleting item ${itemID}${shopID}`)
  return axios
    .delete(`${endpoint}${DELETE_ITEM}`, {withCredentials: true, params: {itemID, shopID}})
    .then((res) => {
      return res.data;
    })
    .catch((err) => {
      throw err;
    });
};

export {
    getItemList,
    submitItem,
    deleteItem
};