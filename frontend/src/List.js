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
import { bignumber, compositionDependencies } from "mathjs"

const { Option } = Select;
export default function ItemList(props) {
  const [data, setData] = React.useState([]);
  const [shopeeLink, setShopeeLink] = React.useState("")
  const [totalPages, setTotalPages] = React.useState(0)
  const [currentPage, setCurrentPage] = React.useState(1)
  // fetch user's favourite items upon page load
  React.useEffect(() => {
    onRefresh(currentPage);
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
  const onRefresh = (currentPage) => {
    getItemList(currentPage - 1)
      .then((res) => {
        if (res.errorCode && res.errorCode !== -1) {
          showFailureMsg("Error occured when fetching items")
        } else if (res.items && res.totalPages) {
          setData(res.items ? res.items : [])
          setTotalPages(res.totalPages ? res.totalPages : 0)
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

  const onDelete = (itemID, shopID) => {
    deleteItem(itemID, shopID)
      .then((res) => {
        if (!res.errorCode) {
          return showFailureMsg("Unexpected error occured. Please try again later!")
        }
        if (res.errorCode && res.errorCode !== -1) {
          switch (res.errorCode) {
            default:
              return showFailureMsg("Unexpected error occured. Please try again later!")
          }
        }
        const dataLength = data.length
        // remove the item from the current data being displayed
        setData(
          data.filter((item) => {
            console.log(item)
            return item.itemID !== itemID && item.shopID !== shopID
          })
        )
        // if no more items to display on this page
        if (dataLength === 1) {
          if (currentPage > 1) {
            const newPage = currentPage - 1
            setCurrentPage(newPage)
            onRefresh(newPage)
          } else {
            // current page is already 1
            onRefresh(1)
          }
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

  const onPageChange = (page) => {
    if (page !== currentPage) {
      setCurrentPage(page)
      onRefresh(page)
    }
  }

  const SHOPEE_LINK_REGEX = /^https:\/\/shopee\.sg\/[\S]+-i\.([0-9]+)\.([0-9]+)(\?[\S]+)?$/

  const onAddToFavourites = () => {
    var matches = SHOPEE_LINK_REGEX.exec(shopeeLink)
    if (!matches || !matches[1] || !matches[2]) {
      showFailureMsg("Incorrect format!")
    } else {
      var itemID = matches[2]
      var shopID = matches[1]
      submitItem(itemID, shopID)
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
            if (currentPage === 1) {
              // only add in the item if the current page is the first page, since we display newest items first
              setData([res.item, ...data].slice(0, 5))
            }
            showSuccessMsg("Successfully added!")
          } else {
            showFailureMsg("Something went wrong with fetching your item. Please try again!")
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
        <Button onClick={() => onRefresh(currentPage)}>
          Refresh
        </Button>
        <div style={{ overflowY: "auto", height: "100%" }}>
          {!data || data.length === 0 && "It's empty here."}
          {data &&
            data.map((item) => {
              return (
                <Card key={`${item.itemID}`} style={{
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
                        onClick={() => onDelete(item.itemID, item.shopID)}
                      >
                        <DeleteFilled />
                      </Button>
                    </Col>
                  </Row>
                </Card>
              )
            })}
        </div>
        <Pagination
          current={currentPage}
          onChange={onPageChange}
          pageSize={5}
          total={totalPages * 5}
        />
      </Col>
    </Row>
  );
}