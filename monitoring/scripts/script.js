import http from 'k6/http';
import { check, group, sleep } from "k6";

const BASE_URL = "http://docker.for.mac.host.internal:80/api"
const SIGNUP = "/user/signup"
const LOGIN = "/user/login"
const GET_FAV_LIST = "/item/get/list"
const ADD_FAV = "/item/add/fav"
const DELETE_FAV = "/item/delete/fav"


export let options = {
  stages: [
      // Ramp-up from 1 to 10 virtual users (VUs) in 5s
      { duration: "5s", target: 10 },

      // Stay at rest on 5 VUs for 10s
      { duration: "10s", target: 10 },

      // Ramp-down from 10 to 0 VUs for 5s
      { duration: "5s", target: 0 }
  ]
};

const params = {
  headers: {
    Accepts: "application/json"
  }
}


// list of itemIDs and shopIDs
const items = [
  {itemID: "9168461947", shopID: "361953971"},
  {itemID: "3049602146", shopID: "361953971"},
  {itemID: "5378451122", shopID: "315440631"},
  {itemID: "16777222699", shopID: "286117947"},
]

// get a random integer in the range [min, max]
// used to generate random itemIDs and shopIDs (to test invalid item)
function getRndInteger(min, max) {
  return Math.floor(Math.random() * (max - min) ) + min;
}

export default function () {
  const userPayload = {
    username: `testuser-${__VU}`,
    password: "password"
  }
  // create a new user (this will fail if the user already exists, so just check that 200 status ok)
  const signupRes = http.post(`${BASE_URL}${SIGNUP}`, JSON.stringify(userPayload), params)
  check(signupRes, {
    resStatus200: (r) => r.status === 200
  })

  // store received cookies from login here
  const jar = http.cookieJar()

  const loginRes = http.post(`${BASE_URL}${LOGIN}`, JSON.stringify(userPayload), params)
  check(loginRes, {
    resStatus200: (r) => r.status === 200,
    // noErrors: (r) => {
    //   var resBody = JSON.parse(r.body)
    //   return resBody.errorCode && resBody.errorCode === -1
    // }
  });

  // add received cookies to cookie jar
  new Map(Object.entries(loginRes.cookies)).forEach((v, k) => {
    jar.set(BASE_URL, v[0].name, v[0].value)
  });
  sleep(.300);

  // get list
  group("get_fav_list_success", () => {
    const getListRes = http.get(`${BASE_URL}${GET_FAV_LIST}?page=0`);
    check(getListRes,
      {
        resStatus200: (r) => r.status === 200,
        noErrors: (r) => {
          var resBody = JSON.parse(r.body)
          return resBody.errorCode && resBody.errorCode === -1
        }
    });
  })
  sleep(.300);

  // add items successfully
  group("add_fav_success", () => {
    for (var item of items) {
      const addItemRes = http.post(`${BASE_URL}${ADD_FAV}`, JSON.stringify(item), params);
      check(addItemRes,
        {
          resStatus200: (r) => r.status === 200,
          // noErrors: (r) => {
          //   var resBody = JSON.parse(r.body)
          //   console.log("Add fav error: " + resBody.errorCode)
          //   return resBody.errorCode && resBody.errorCode === -1
          // }
      });
      sleep(0.5)
    }
  })
  sleep(0.3)

  // delete the items added
  group("delete_fav_success", () => {
    for (var item of items) {
      const getItemRes = http.del(`${BASE_URL}${DELETE_FAV}?itemID=${item.itemID}&shopID=${item.shopID}`);
      check(getItemRes,
        {
          resStatus200: (r) => r.status === 200,
          // noErrors: (r) => {
          //   var resBody = JSON.parse(r.body)
          //   console.log("Delete fav error: " + resBody.errorCode)
          //   return resBody.errorCode && resBody.errorCode === -1
          // }
      });
      sleep(0.5)
    }
  })

  sleep(0.3)

  // add random invalid items
  group("add_item_nonexisting", () => {
    for (let i = 0; i < 5; i++) {
      let itemID = getRndInteger(1, 1000).toString()
      let shopID = getRndInteger(1, 1000).toString()
      let item = {itemID, shopID}
      const addItemRes = http.post(`${BASE_URL}${ADD_FAV}`, JSON.stringify(item), params);
      check(addItemRes,
        {
          resStatus200: (r) => r.status === 200,
          hasErrors: (r) => {
            var resBody = JSON.parse(r.body)
            return resBody.errorCode && resBody.errorCode !== -1
          }
      });
      sleep(0.5)
    }
  })
};