import { useEffect, useState } from "react";

interface UserItemInterface {
    ItemID: number;
    TypeID: number;
    ItemName: string;
    ItemDescription: string;
    ImgURL: string;
    IsListed: boolean
}
interface sellItem {
  ItemID: number;
  Token: string | null;
  Price: number
}

const ItemTableComponent = () => {
    const [sellPrice, setSellPrice] = useState(0)
    const [userItems, setUserItems] = useState<UserItemInterface[]>([]);
    let fetchString =  "http://ronstad.se/inventory/" + localStorage.getItem("uid")
    useEffect(() => {
        fetch(fetchString, { method: "GET" }) // Replace with your actual API URL
            .then((response) => response.json())
            .then((useritems) => setUserItems(useritems))
            .catch((error) => console.error("Error: ", error));
        }, []);

      
    function rmListing(it: UserItemInterface): void {
        console.log(it.ItemID);
        var adr: string = "http://ronstad.se/Marketplace/removeListing" + it.ItemID
        fetch(adr, { method: "GET" })
        .then((response) => {
          if (response.ok === true) {
            alert("Item removed from marketplace");
          } else {
            console.log("Invalid Request");
            alert("nuh uh");
          }
        })
    }
    function handleSell(itemID: number) {
  
      // Prevent the browser from reloading the page
      //e.preventDefault();
      let tmp: sellItem = {ItemID: itemID, Token: localStorage.getItem("token"), Price: sellPrice}
      //alert("test");
      console.log(JSON.stringify(tmp));
      fetch("http://ronstad.se/Marketplace/addListing", { method: "POST", body: JSON.stringify(tmp) })
      .then((response) => {
        if (response.ok === true) {
          console.log(JSON.stringify(tmp));
          alert("Item succesfully put on marketplace");
        } else {
          console.log(JSON.stringify(tmp));
          console.log("Invalid");
          alert("nuh uh");
        }
      })
  }
    if (userItems == null){
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
                    {/*map either button or text input if item is already listed */}
                    <td>
                        {(<><input type={"number"} onChange={(e) => {setSellPrice(e.target.valueAsNumber || 0)}} name="sellPrice" id="sellPrice"/> <button onClick={() => handleSell(item.ItemID)} className={"cyber-button-small bg-blue fg-yellow"}>Sell</button></>)}
                    </td>
                </tr>
            ))}
        </>
    );
}

export default ItemTableComponent