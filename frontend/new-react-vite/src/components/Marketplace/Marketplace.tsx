import "./Marketplace.css";
import "../cyberpunk-css-main/cyberpunk.css";
import CyberpunkWindow from "../cyberpunkWindow";
import { useEffect, useState } from "react";
import { globalAddr } from "../../header.tsx"
import PopUpComments from "../PopUpComments.tsx";


// I am really cool
interface MarketplaceState {
  sortBy: string;
  search: string;
}
interface userIDInt {
  UserID: number
}


interface MarketplaceItems {
  ItemID: number;
  TypeID: number;
  UserID: number;

  ItemName: string;
  ItemDescription: string;
  ImgURL: string;
  
  OfferID: number;
  Price: number;
  CreationDate: string;

  Username: string;
  InCheckout: number;
}




function Marketplace() {
  const [sortState, setSortState] = useState<MarketplaceState>({sortBy: "Newest", search: ""})
  const [marketplaceItems, setMarketplaceItems] = useState<MarketplaceItems[]>([]);
  
  const [showComments, setShowComments] = useState(false)
  const [fuckTS, setFuckTS] = useState<number | null>(null);

  // console.log("showComments state:", showComments);
  const [cartStatus, setCartStatus] = useState<{ [key: number]: boolean }>({});
  onchange = s => { // <-- wtf is this
    console.log(sortState.sortBy)
    console.log(s)
  }
  

  useEffect(() => {
    var fetchString = `http://${globalAddr}/Marketplace/displayMarket`;
    fetch(fetchString, { method: "POST", body: JSON.stringify(sortState) })
      .then((response) => response.json())
      .then((marketplaceItems) => {
        setMarketplaceItems(marketplaceItems);
        
        // Run checkCart for each item
        marketplaceItems.forEach((item: MarketplaceItems) => {
          checkCart(item);
        });
      })
      .catch((error) => console.error("Error: ", error));
  }, [sortState]);

  const buyItem = (item: MarketplaceItems) => {
    console.log("OfferID: ?", item.OfferID);
    var tmpInt: userIDInt = {UserID: Number(localStorage.getItem("uid"))}
    const jsonItem = JSON.stringify(tmpInt)
    console.log(tmpInt)
    console.log(jsonItem)
    var fetchString = `http://`+globalAddr+`/Marketplace/buy/` + item.ItemID
    
    fetch(fetchString, { method: "POST", body:  jsonItem}) 
            .then((response) => response.json())
            .then((marketplaceItems) => {setMarketplaceItems(marketplaceItems)
              window.location.reload();
            })
            .catch((error) => console.error("Error: ", error));

  };

  const addToCart = (item: MarketplaceItems) => {
    var tmpInt: userIDInt = {UserID: Number(localStorage.getItem("uid"))}
    const jsonItem = JSON.stringify(tmpInt)
    fetch("http://"+globalAddr+"/Marketplace/addToCart/"+ item.OfferID, { method: "POST", body:  jsonItem}) 
    .then((response) => {response.json()
      alert("Item successfully added to cart")
      window.location.reload();
    })
    .catch((error) => console.error("Error: ", error));
  }
  const removeFromCart = (item: MarketplaceItems) => {
    var tmpInt: userIDInt = {UserID: Number(localStorage.getItem("uid"))}
    const jsonItem = JSON.stringify(tmpInt)
    fetch("http://"+globalAddr+"/Marketplace/removeFromCart/"+ item.OfferID, { method: "POST", body:  jsonItem}) 
    .then((response) => {response.json()
      alert("Item successfully removed from cart")
      window.location.reload();
    })
    .catch((error) => console.error("Error: ", error));
  }

  const checkCart = async (item: MarketplaceItems) => {
    const tmpInt: userIDInt = { UserID: Number(localStorage.getItem("uid")) };
    const jsonItem = JSON.stringify(tmpInt);
  
    try {
      const response = await fetch(`http://${globalAddr}/Marketplace/checkCart/${item.OfferID}`, {
        method: "POST",
        body: jsonItem,
      });
  
      const data = await response.json();
      console.log(data.InCheckout);
  
      setCartStatus(prevState => ({
        ...prevState,
        [item.ItemID]: data.InCheckout || false
      }));
    } catch (error) {
      console.error("Error checking cart:", error);
    }
  };
  const market = () => {
    if (marketplaceItems === null) {
      console.log("empty")
      return (
        <tr>
          <td className="">itemname</td>
          <td className="">Price</td>
          <td className="">PlaceHolder, Your marketplace is empty</td>
          <td className="">Seller</td>
          <td>
            <input className='buy-button' type='button' value='Buy' />
          </td>
        </tr>
      )
    }
    else {
      return (
          marketplaceItems.map((item: MarketplaceItems) => (
              <tr key={item.ItemID}>
                {/* <td className="">
                  <button onClick={() => {
                    console.log("Opening comments window"); // Debugging
                    setShowComments(true);
                  }}>
                    {item.ItemName}
                  </button>
                </td> */}
                <td style={{ cursor: 'pointer' }} onClick={() => {setFuckTS(item.TypeID); setShowComments(true); } } >{item.ItemName}</td>
                <td className="">{item.Price}</td>
                <td className="">{item.ItemDescription}</td>
                <td className="">{item.Username}</td>
                <td>
                  <input onClick={() => buyItem(item)} className='buy-button' type='button' value='Buy' />
                  {cartStatus[item.ItemID] ? 
                        (<input onClick={() => removeFromCart(item)} className="buy-button" type="button" value="Remove From cart" />
                      )
                         : (<input onClick={() => addToCart(item)} className="buy-button" type="button" value="Add To Cart" />
                        )}
                </td>
              </tr>)
        ) 
      )
    }
  } 

  return (
  
  <div>
    {/* {showComments && <PopUpComments onClose={() => setShowComments(false)} />} */}
    {showComments && (
      <PopUpComments
        onClose={() => setShowComments(false)}
        itemId={fuckTS} // Pass item.ItemID as a prop
      />
    )}
    <div className="left-right-container">
      <div className="left">
        <CyberpunkWindow>
        <div className="left-elements">
          <div className="sorting oxanium-font">
            Sort by:
            <div className="cyber-input">
              <div className="cyber-select ac-purple ">
                <select className="oxanium-font"  onChange={(s) => setSortState({sortBy: s.target.value, search: sortState.search})}>
                    <option style={{color:"black"}}className="oxanium-font" value="Newest"           >Newest</option>
                    <option style={{color:"black"}}className="oxanium-font" value="Oldest"           >Oldest</option>
                    <option style={{color:"black"}}className="oxanium-font" value="Price_Ascending"  >Price Ascending</option>
                    <option style={{color:"black"}}className="oxanium-font" value="Price_Descending" >Price Descending</option>
                    <option style={{color:"black"}}className="oxanium-font" value="Alphabetically_Ascending" >Alphabetically A-Ö</option>
                    <option style={{color:"black"}}className="oxanium-font" value="Alphabetically_Descending" >Alphabetically Ö-A</option>
                </select>
              </div>
            </div>
          </div>

          <div className="tags">  
                    
            <div className="tagtext">
            <header>Tags</header> 
              <ul>
                <input type="checkbox" className="cyber-check-tags" ></input> Cyberware <br></br>
                <input type="checkbox" className="cyber-check-tags" ></input> Consumables<br></br>
                <input type="checkbox" className="cyber-check-tags" ></input> Weapons<br></br>
                <input type="checkbox" className="cyber-check-tags" ></input> Drugs<br></br>
                <input type="checkbox" className="cyber-check-tags" ></input> Armor<br></br>
              </ul>
            </div>
          </div>
        </div>
        </CyberpunkWindow>
        
      </div>

      <div className="right">
        
        {/* store */}
        <table className="cyber-table store-table ac-custom">
          <thead>
            <tr className="thead">
              <th>Item</th>
              <th>Price</th>
              <th>Description</th>
              <th>Seller</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {market()}
          </tbody>
        </table>
      </div>
    </div>
  
  </div>
  );
}


export default Marketplace;