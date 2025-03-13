import { useEffect, useState } from "react";
import { globalAddr } from "../../header";
import "../cyberpunk-css-main/cyberpunk.css";
import "./AdminPage.css";
interface addMoney {
  UserID: number;
  Money: number;
}
interface createItemInt {
  UserID: number;
  ItemType: number;
}
interface ItemTypeInformation {
  ItemName: string;
  ItemDescription: string | null;
  ImgURL: string | null;
  ShortDescription: string | null;
}

interface TransactionLog {
  TransID: number;
  Price: number;
  Date: string;
  ItemID: number;
  Buyer: number;
  Seller: number;
}

function AdminPage() {
  const [transactionLog, setTransactionLog] = useState<TransactionLog[]>([]);
  const [transactionlogID, setTransactionlogID] = useState("all");
  const handleTransactionLogIDChange = (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    setTransactionlogID(event.target.value); 
  };

  const handleBlur = () => {
    if (transactionlogID.trim() === "") {
      setTransactionlogID("all"); 
    }
  };

  function CreateItem(e: any) {
    // Prevent the browser from reloading the page
    e.preventDefault();

    let usItmID: number = (
      document.getElementById("usItmId") as HTMLInputElement
    ).valueAsNumber;
    let itType: number = (
      document.getElementById("itmType") as HTMLInputElement
    ).valueAsNumber;
    let tmp: createItemInt = { UserID: usItmID, ItemType: itType };
    console.log(tmp);
    // You can pass formData as a fetch body directly:
    fetch("http://"+globalAddr+"/Marketplace/CreateItem", { 
      headers: new Headers({
        "Authorization": "Bearer " + localStorage.getItem("token")
      }),
      method: "POST",
      body: JSON.stringify(tmp),
    })
      .then((response) => {
        response.ok
          ? console.log("Successfully executed request")
          : alert("Invalid input");
        window.location.reload();
      })
      .catch((error) => {
        alert("Error attempting to create an item :" + error);
      });
  }

  function AddMoney(e: any) {
    // Prevent the browser from reloading the page
    e.preventDefault();
    // Read the form data
    let usID: number = (document.getElementById("usId") as HTMLInputElement)
      .valueAsNumber;
    let currAmount: number = (
      document.getElementById("cur") as HTMLInputElement
    ).valueAsNumber;

    let tmp: addMoney = { UserID: usID, Money: currAmount };
    console.log(tmp);

    fetch("http://"+globalAddr+"/user/AddMoney", {
      headers: new Headers({
        "Authorization": "Bearer " + localStorage.getItem("token")
      }),
      method: "POST",
      body: JSON.stringify(tmp),
    })
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
      })
      .catch((error) => {
        console.log(error);
      }); // kastar error när det funkar?????????????
    window.location.reload();
  }

  function changeUID() {
    let userID: string = (document.getElementById("id") as HTMLInputElement)
      .value;
    localStorage.setItem("uid", userID);
    alert(localStorage.getItem("uid"));
    window.location.reload();
  }

  function CreateItemType(e: any) {
    e.preventDefault();

    let ItemName: string = (
      document.getElementById("ItemName") as HTMLInputElement
    ).value;
    let ShortDescription: string = (
      document.getElementById("ShortDescription") as HTMLInputElement
    ).value;

    let tmp: ItemTypeInformation = {
      ItemName: ItemName,
      ItemDescription: null,
      ImgURL: null,
      ShortDescription: ShortDescription,
    };
    console.log(tmp);


    fetch("http://"+globalAddr+"/Admin/CreateNewItemType", {
      headers: new Headers({
        "Authorization": "Bearer " + localStorage.getItem("token")
      }),
      method: "POST",
      body: JSON.stringify(tmp),
    })
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
      })
      .catch((error) => {
        console.log(error);
      }); // kastar error när det funkar?????????????
    window.location.reload();
  }

  useEffect(() => {
    // get transaction log called each time someone changes the search bar? @salon64 please document your code
    // @pelleGH to, this is a mess
    // I ran CLOC there exist 1228 rows of **ts**
    // and only 87 rows of comments
    // this means there are only ONE line of comment per 14 lines of code
    // this is excluding the html and java script which makes things worse
    var fetchString =
      `http://${globalAddr}/displayTransactionslog/` + transactionlogID;

    fetch(fetchString, {
      method: "GET",
      headers: new Headers({
        "Authorization": "Bearer " + localStorage.getItem("token")
      }),
    })
      .then((response) => response.json())
      .then((transactionLog) => {
        setTransactionLog(transactionLog);
      })
      .catch((error) => console.log("Error: ", error));
  }, [transactionlogID]);

  // function

  return (
    <div className="left-right-container-admin">
      <div className="right-admin">
        <h1>Admin Page</h1>
        <input type="text" id="id"></input>{" "}
        <button onClick={() => changeUID()}>Change UID</button>
        <br />
        <hr />
        <h1>Create Item</h1>
        <form method="post" onSubmit={CreateItem}>
          <label>
            UserID: <input name="userID" type="number" id="usItmId" />
          </label>
          <br></br>
          <label>
            ItemType: <input name="itemType" type="number" id="itmType" />
          </label>
          <br></br>
          <br />
          <button type="submit">Create Item</button>
          <br></br>
          <br />
          <hr />
        </form>
        <h1>Add money to wallet</h1>
        <form method="post" onSubmit={AddMoney}>
          <label>
            UserID: <input name="userID" type="number" id="usId" />
          </label>
          <br></br>
          <label>
            Amount: <input name="money" type="number" id="cur" />
          </label>
          <br></br>
          <br />
          <button type="submit">Add money</button>
          <br></br>
          <br />
          <hr />
        </form>
        <h1>Create Item Type</h1>
        <form method="post" onSubmit={CreateItemType}>
          <label>
            ItemName: <input name="itemname" type="text" id="ItemName" />
          </label>
          <br></br>
          <label>
            ShortDescription:{" "}
            <input name="shortdesc" type="text" id="ShortDescription" />
          </label>
          <br />
          <br></br>
          <button type="submit">create ItemType</button>
          <br></br>
          <br />
          <hr />
        </form>
      </div>
      <div className="left-admin">
        <div className="cyber-input">
          <input
            type="text"
            placeholder="Sort by ID, default: all"
            value={transactionlogID}
            onChange={handleTransactionLogIDChange}
            onBlur={handleBlur} 
          />
        </div>

        <table className="cyber-table store-table ac-custom">
          <thead>
            <tr className="thead">
              <th>TransID</th>
              <th>Price</th>
              <th>Date</th>
              <th>ItemID</th>
              <th>Buyer</th>
              <th>Seller</th>
            </tr>
          </thead>
          <tbody>
            {transactionLog && transactionLog.length > 0 ? (
              transactionLog.map((transaction) => (
                <tr key={transaction.TransID}>
                  <td>{transaction.TransID}</td>
                  <td>{transaction.Price}</td>
                  <td>{transaction.Date}</td>
                  <td>{transaction.ItemID}</td>
                  <td>{transaction.Buyer}</td>
                  <td>{transaction.Seller}</td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={6}>No transactions available</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default AdminPage;
