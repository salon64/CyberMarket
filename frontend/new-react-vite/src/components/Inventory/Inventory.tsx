import "./Inventory.css";
import "../cyberpunk-css-main/cyberpunk.css";
import ItemTableComponent from "../ItemTable";
import { useState } from "react";

const Inventory = () => {
  getMoney()
  
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

  const [wallet, setWallet] = useState("0")
  
    return (
      <body>
        <header>My Inventory</header> Wallet: {wallet}
    <div className="left-right-container">
      <div className="right">
        {/* store */}
        
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
  </body>
    );
  };
  
  export default Inventory;