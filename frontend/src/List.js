import React from "react";
import {
  Select,
  Input,
  Popover,
  Tag,
  Row,
  Col,
  Button,
  Card,
  message,
  Pagination
} from "antd";
import { DeleteFilled, EditFilled, EnterOutlined } from "@ant-design/icons";
import { getItemList, deleteItem, submitItem }
  from "./api/api";
import { bignumber } from "mathjs"

const { Option } = Select;
export default function ItemList(props) {
  const [data, setData] = React.useState([]);
  const [shopeeLink, setShopeeLink] = React.useState("")
  const [totalPages, setTotalPages] = React.useState(0)
  const [currentPage, setCurrentPage] = React.useState(1)
  // fetch user's favourite items upon page load
  const [errorMessage, setErrorMessage] = React.useState("")
  React.useEffect(() => {
    onRefresh();
  }, []);

  // shows a success message to user upon successful adding of item
  const showSuccessMsg = (msg) => {
    message.success({
      content: msg,
      style: { marginTop: "5vh" },
      duration: 1,
    });
  };

  const showFailureMsg = (msg) => {
    message.error({
      content: msg,
      style: { marginTop: "5vh" },
      duration: 1,
    });
  };

  // refreshes user's favourites list
  const onRefresh = () => {
    console.log(`totalPages: ${totalPages}`)
    console.log(`current page: ${currentPage}`)
    getItemList(currentPage - 1)
      .then((res) => {
        if (res.errorCode && res.errorCode !== -1) {
          showFailureMsg("Error occured when fetching items")
        } else if (!res.items || res.items.length === 0 || !res.totalPages) {
          showFailureMsg("Nothing in favourites - submit something!")
        } else {
          setData(res.items)
          setTotalPages(res.totalPages)
        }
      })
      .catch((err) => {
        if (err.response && err.response.status === 401) {
          showFailureMsg("Session timeout. You have been logged out");
          props.onUnauthentication();
          return;
        }
        showFailureMsg("Unexpected error occured");
      });
  };

  const onDelete = (tag, index) => {
    deleteItem(tag, index)
      .then((res) => {
        if (!res.errorCode) {
          return setErrorMessage("Unexpected error occured. Please try again later!")
        }
        if (res.errorCode && res.errorCode !== -1) {
          switch (res.errorCode) {
            default:
              return setErrorMessage("Unexpected error occured. Please try again later!")
          }
        }
        // show success message TODO
        onRefresh();
      })
      .catch((err) => {
        if (err.response && err.response.status === 401) {
          showFailureMsg("Session timeout. You have been logged out");
          props.onUnauthentication();
          return;
        }
        showFailureMsg("Unexpected error occured");
      });
  };

  const onPageChange = (page) => {
    // console.log(`current page: ${page}`)
    setCurrentPage(page)
    onRefresh()
  }

  const SHOPEE_LINK_REGEX = /^https:\/\/shopee\.sg\/[\S]+-i\.([0-9]+)\.([0-9]+)(\?[\S]+)?$/

  const onAddToFavourites = () => {
    var matches = SHOPEE_LINK_REGEX.exec(shopeeLink)
    if (!matches || !matches[1] || !matches[2]) {
      console.log("Not shopee link format")
      showFailureMsg("Incorrect format!")
    } else {
      var itemId = matches[2]
      var shopId = matches[1]
      submitItem(itemId, shopId)
        .then((res) => {
          if (res.errorCode && res.errorCode != -1) {
            switch (res.errorCode) {
              case 340011:
                return showFailureMsg("Item already in favourites, find something else!")
              default:
                return showFailureMsg(
                  "Something went wrong, please try again later!"
                )
            }
          }
          if (res.item) {
            setData([res.item, ...data].slice(0, 5))
          }
        })
        .catch(err => {
          console.log(err)
          if (err.response && err.response.status === 401) {
            showFailureMsg("Session timeout. You have been logged out");
            props.onUnauthentication();
            return;
          }
          showFailureMsg("Unexpected error occured")
        })
    }
  }

  return (
    <Row type="flex" justify="center" style={{ height: "100%" }}>
      <Col>
        <Input placeholder="Submit Shopee link" onChange={(e) => setShopeeLink(e.target.value)}></Input>
        <Button onClick={onAddToFavourites}>Add to Favourites!</Button>
        <Button onClick={onRefresh}>
          Refresh
        </Button>
        <div style={{ overflowY: "auto", height: "100%" }}>
          {data &&
            data.map((item) => {
              return (
                <Card key={`${item.itemId}`} style={{
                  borderStyle: "solid"
                }}>
                  <Row>
                    <Col span={19}>
                      <p>{item.name}</p>
                      <p>${item.price*1.0/100000.0}</p>
                    </Col>
                    <Col span={2}>
                      <Button
                        style={{ backgroundColor: "lightblue" }}
                        onClick={() => onDelete(item.itemId, item.shopId)}
                      >
                        <DeleteFilled />
                      </Button>
                    </Col>
                  </Row>
                </Card>
              );
            })}
        </div>
        <Pagination current={currentPage} onChange={onPageChange} pageSize={5} total={totalPages * 5} />;
      </Col>
    </Row>
  );
}