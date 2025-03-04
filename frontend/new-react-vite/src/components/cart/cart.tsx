import { useEffect, useState } from "react";
import "./cart.css";
import { globalAddr } from "../../header";

interface CartItems {
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
}interface userIDInt {
  UserID: number
}
function Cart(){
  const [cartItems, setCartItems] = useState<CartItems[]>([]);
    useEffect(() => {
      var fetchString = `http://${globalAddr}/Marketplace/displayCart/` + localStorage.getItem("uid");
      fetch(fetchString, { method: "GET" })
        .then((response) => response.json())
        .then((cartItems) => {
          setCartItems(cartItems);})
        .catch((error) => console.error("Error: ", error))},[])
  const removeFromCart = (item: CartItems) => {
    var tmpInt: userIDInt = {UserID: Number(localStorage.getItem("uid"))}
    const jsonItem = JSON.stringify(tmpInt)
    fetch("http://"+globalAddr+"/Marketplace/removeFromCart/"+ item.OfferID, { method: "POST", body:  jsonItem}) 
    .then((response) => {response.json()
      window.location.reload();
      alert("Item successfully removed from cart")
      
    })
    .catch((error) => console.error("Error: ", error));
  }
  const cartCheckout = () => {
    fetch("http://"+globalAddr+"/Marketplace/cartCheckout/"+ localStorage.getItem("uid"), { method: "GET"}) 
    .then((response) => {response.json()
      alert("Items successfully checked out from cart")
    })
    .catch((error) => console.error("Error: ", error))
  }
  const market = () => {
            if (cartItems === null) {
              console.log("empty")
              return (
                <tr>
                  <td className="">Item</td>
                  <td className="">Price</td>
                  <td className="">PlaceHolder, Your cart is empty</td>
                  <td className="">Seller</td>
                  <td>
                  </td>
                </tr>
              )
            }
            else {
              return (<>
                  {cartItems.map((item: CartItems) => (
                      <tr key={item.ItemID}>
                      <td className="">{item.ItemName}</td>
                      <td className="">{item.Price}</td>
                      <td className="">{item.ItemDescription}</td>
                      <td className="">{item.Username}</td>
                      <td className=""><input onClick={() => removeFromCart(item)} className="buy-button" type="button" value="Remove From cart" />
                      </td>
                    </tr>
                      
                      ))}
                      <tr>
                        <td>Subtotal: </td>
                        <td>{cartItems.reduce((total, item) => total + item.Price, 0)}</td>
                      </tr>
                      <br>
                      </br>
                      <button onClick={() => cartCheckout()} value="Checkout">Checkout</button>
                      </>
                    )}} 
        
          return (
            
          <body>
            <div className="left-right-container">
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
                    <tr>
                      </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </body>
      )};
export default Cart;