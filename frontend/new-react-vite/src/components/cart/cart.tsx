import "./cart.css";
function Cart(){
    return (
    <div className="right">      <table className="cyber-table store-table ac-custom">
    <thead>
      <tr className="thead">
        <th>Item</th>
        <th>Price</th>
        <th>Description</th>
        <th>Seller</th>
        <th>Remove</th>
      </tr>
    </thead>
    <tbody>
        <tr className="">
            <td className="">{}(ItemPlaceholder)</td>
            <td className="">{}(PricePlaceholder)</td>
            <td className="">{}(DescriptionPlaceholder)</td>
            <td className="">{}<image className="sellIcon"></image></td>
            <td className="">{}<button className="rm-button">X</button></td>
        </tr>
        <tr className="">
            <td>Total:</td>
            <td className="">{}</td>
        </tr>
      
    </tbody>
    </table>
    <button>Checkout</button>
    </div>
   

    );
}
export default Cart;