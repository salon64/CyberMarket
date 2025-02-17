import { useEffect, useState } from "react";

interface UserItemInterface {
    ItemID: number;
    TypeID: number;
    ItemName: string;
    ItemDescription: string;
    ImgURL: string;
  }
  async function handleSell(itemID) {
    // Prevent the browser from reloading the page
    //e.preventDefault();
    
    let state = {
        sellP: (document.getElementById("sellPrice") as HTMLInputElement).value,
        token: localStorage.getItem("token"),
        itemID: itemID
    }
    //alert("test");
    console.log(JSON.stringify(state));
    fetch("http://ronstad.se/Marketplace", { method: "POST", body: JSON.stringify(state) })
    .then((response) => {
      if (response.ok === true) {
        console.log(JSON.stringify(state));
        alert("Item succesfully put on marketplace");
      } else {
        console.log(JSON.stringify(state));
        console.log("Invalid");
        alert("nuh uh");
      }
    })
}
const ItemTableComponent = () => {
    const [userItems, setUserItems] = useState<UserItemInterface[]>([]);
    useEffect(() => {
        fetch("http://ronstad.se/inventory/1", { method: "GET" }) // Replace with your actual API URL
            .then((response) => response.json())
            .then((useritems) => setUserItems(useritems))
            .catch((error) => console.error("Error: ", error));
        }, []);

      
    return (
        <>
            {userItems.map((item) => (
                <tr key={item.ItemID} className="">
                    <td className="">{item.ItemID}</td>
                    <td className="">{item.TypeID}</td>
                    <td className="">{item.ItemName}</td>
                    <td><input type={"text"} name="sellPrice" id="sellPrice"/> <button onClick={() => handleSell(item.ItemID)} className={"cyber-button-small bg-blue fg-yellow"}>Sell</button></td>
                </tr>
            ))}
        </>
    );
}

export default ItemTableComponent