import { useEffect, useState } from "react";
import { globalAddr } from "../header";

interface UserItemInterface {
  ItemID: number;
  TypeID: number;
  ItemName: string;
  IsListed: number;
  ItemDescription: string;
  ImgURL: string;
}
interface sellItem {
  ItemID: number;
  Token: string | null;
  Price: number
}

const ItemTableComponent = () => {
  const [sellPrice, setSellPrice] = useState(0)
  const [userItems, setUserItems] = useState<UserItemInterface[]>([]);
  let fetchString = "http://"+globalAddr+"/inventory/" + localStorage.getItem("uid")
  useEffect(() => {
    fetch(fetchString, { method: "GET" }) // Replace with your actual API URL
      .then((response) => response.json())
      .then((useritems) => {
        setUserItems(useritems)
        console.log(useritems)
      })
      .catch((error) => console.error("Error: ", error));
  }, []);


  function rmListing(it: UserItemInterface): void {
    console.log(it.ItemID);
    var adr: string = "http://"+globalAddr+"/Marketplace/removeListing/" + it.ItemID
    console.log(adr)
    fetch(adr, { method: "GET" })
      .then((response) => {
        if (response.ok === true) {
          alert("Item removed from marketplace");
        } else {
          console.log("Invalid Request");
          alert("nuh uh");
        }
      })
      .catch(error => console.log(error))
    window.location.reload();
  }
  function handleSell(itemID: number) {

    // Prevent the browser from reloading the page
    //e.preventDefault();
    let tmp: sellItem = { ItemID: itemID, Token: localStorage.getItem("token"), Price: sellPrice }
    //alert("test");
    console.log(JSON.stringify(tmp));
    console.log("start")
    fetch("http://"+globalAddr+"/Marketplace/addListing", { method: "POST", body: JSON.stringify(tmp) })
      .then((response) => {
        if (response.ok === true) {
          console.log(JSON.stringify(tmp));
          console.log("Item succesfully put on marketplace")
          alert("Item succesfully put on marketplace");
          window.location.reload();
        } else {
          console.log(JSON.stringify(tmp));
          console.log("Invalid");
          alert("nuh uh");
        }
      })
      .catch(error => {
        console.log(error)
        console.log("catch")
      })
  }
  if (userItems == null) {
    return (
      <>
      </>
    )
  }

  return (
    <>
      {userItems.map((item) => (
        <tr key={item.ItemID} className="">
          <td className="">{item.ItemID}</td>
          <td className="">{item.TypeID}</td>
          <td className="">{item.ItemName}</td>
          <td className="">{item.IsListed}</td>
          {/*map either button or text input if item is already listed */}
          <td>
            {item.IsListed ? (<button onClick={() => rmListing(item)}>Remove Listing</button>) : (<><input type={"number"} onChange={(e) => { setSellPrice(e.target.valueAsNumber || 0) }} name="sellPrice" id="sellPrice" /> <button onClick={() => handleSell(item.ItemID)} className={"cyber-button-small bg-blue fg-yellow"}>Sell</button></>)}
          </td>
        </tr>
      ))}
    </>
  );
}

export default ItemTableComponent
