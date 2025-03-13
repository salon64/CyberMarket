import "./Inventory.css";
import "../cyberpunk-css-main/cyberpunk.css";
import ItemTableComponent from "../ItemTable";
import { useState } from "react";
import { globalAddr } from "../../header";

const Inventory = () => {
  getMoney()

  function getMoney() {
    fetch("http://" + globalAddr + "/user/getMoney/" + localStorage.getItem("uid"),
      {
        method: "GET",
        headers: new Headers({
          "Authorization": "Bearer " + localStorage.getItem("token")
        }),
      })
      .then((response) => response.ok ? response.json() : response.text().then((r) => alert(r)))
      .then(data => {
        console.count(data)
        const obj = JSON.parse(JSON.stringify(data))
        console.log(obj.Amount)
        setWallet(obj.Amount) //im going to krill myself
      }
    )
  }

  const [wallet, setWallet] = useState("0")

  return (
    <body>
      <header>My Inventory</header> Wallet: {wallet}
      <div className="left-right-container">
        <div className="right">
          {/* store */}
          <table className="cyber-table store-table">
            <thead style={{ backgroundColor: "bisque" }}>
              <tr>
                <th>Item</th>
                <th>ItemType</th>
                <th>ItemName</th>
                <th>Add/Remove from market</th>
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
