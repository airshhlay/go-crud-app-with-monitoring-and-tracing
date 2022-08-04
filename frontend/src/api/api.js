import axios from "axios";
const endpoint = process.env.NODE_ENV === "development" ? process.env.REACT_APP_ENDPOINT_LOCAL : process.env.REACT_APP_ENDPOINT_PROD;

const getItemList = () => {
  return axios
    .get(`${endpoint}/item/get/list`)
    .then((res) => {
      return res.data;
    })
    .catch((err) => {
      console.log(err);
      throw err;
    });
};

const submitItem = (itemId, shopId) => {
  return axios
    .post(
      `${endpoint}/item/add`,
      { itemId, shopId },
    )
    .then((res) => {
      return "ok";
    })
    .catch((err) => {
      console.log(err);
      throw err;
    });
};

const deleteItem = (itemId, shopId) => {
  return axios
    .delete(`${endpoint}/item/delete/${itemId}/${shopId}`)
    .then((res) => "ok")
    .catch((err) => {
      console.log(err);
      throw err;
    });
};

export {
    getItemList,
    submitItem,
    deleteItem
};