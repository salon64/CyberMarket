import "./Inventory.css";
import "../cyberpunk-css-main/cyberpunk.css";
import ItemTableComponent from "../ItemTable";
import { useState } from "react";
interface addMoney {
  UserID: number;
  Money: number;

}
interface createItemInt {
  UserID: number;
  ItemType: number;

}
const Inventory = () => {
  getMoney()
  function CreateItem(e) {
    // Prevent the browser from reloading the page
    e.preventDefault();
  
    // Read the form data
    let usItmID: number = (
      document.getElementById("usItmId") as HTMLInputElement
    ).valueAsNumber;
    let itType: number = (
      document.getElementById("itmType") as HTMLInputElement
    ).valueAsNumber;
    let tmp: createItemInt = {UserID: usItmID, ItemType: itType}
    console.log(tmp)
    // You can pass formData as a fetch body directly: 
    fetch("http://ronstad.se/Marketplace/CreateItem", { method: "POST", body: JSON.stringify(tmp)})
    .then(response =>  {response.ok ? (console.log("Successfully executed request")):(alert("Invalid input"))
      window.location.reload();
    })
    .catch(error => {alert("Error attempting to create an item :"+error)});
  }
  function AddMoney(e) {
    // Prevent the browser from reloading the page
    e.preventDefault();
    // Read the form data
    let usID: number = (
      document.getElementById("usId") as HTMLInputElement
    ).valueAsNumber;
    let currAmount: number = (
      document.getElementById("cur") as HTMLInputElement
    ).valueAsNumber;

    let tmp: addMoney = {UserID: usID, Money: currAmount}
    console.log(tmp)
    
    fetch("http://ronstad.se/user/AddMoney", { method: "POST", body: JSON.stringify(tmp)})
    .then(response => response.json())
    .then(data => {console.log(data)})
    .catch(error => {}); // kastar error när det funkar?????????????
    window.location.reload();
  }
  
  function getMoney() {
    fetch("http://ronstad.se/user/getMoney/" + localStorage.getItem("uid"), { method: "GET"})
    .then(response => response.json())
    .then(data => {
      const obj = JSON.parse(JSON.stringify(data))
      console.log(obj[0].Amount)
      setWallet(obj[0].Amount) //im going to krill myself
      
    })
    //.catch(error => {alert(error)});
  }
  
  function changeUID() {
      let userID: string = (
        document.getElementById("id") as HTMLInputElement
      ).value;
      localStorage.setItem("uid", userID)
      alert(localStorage.getItem("uid"))
      window.location.reload();
  }
  const [wallet, setWallet] = useState("0")
  
    return (
      <body>
        <header>My Inventory</header> Wallet: {wallet}
    <div className="left-right-container">
      <div className="right">
        {/* store */}
        <input type="text" id="id"></input> <button onClick={() => changeUID()}>Change UID</button>
        <button onClick={() => getMoney()}>getwallet</button>
        <table className="cyber-table store-table">
          <thead style={{backgroundColor: "bisque"}}>
            <tr>
              <th>Item</th>
              <th>Description</th>
              <th>Header 3</th>
              <th>Sell</th>
            </tr>
          </thead>
          <tbody>
            <ItemTableComponent />
          </tbody>
        </table>
      </div>
    </div>

    <h1>Create Item</h1>
      <form method="post" onSubmit={CreateItem}>
      <label>
          UserID: <input name="userID" type="number" id="usItmId" />
        </label>
        <br></br>
        <label>
          ItemType: <input name="itemType" type="number" id="itmType"/>
        </label>
        <br></br>
        <hr />
        <button type="submit">
        Create Item
        </button>
        <br>
        </br>
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
        <hr />
        <button type="submit">
        Add money
        </button>
        <br>
        </br>
      </form>
  </body>
    );
  };
  
  export default Inventory;