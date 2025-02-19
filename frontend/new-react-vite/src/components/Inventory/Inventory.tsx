import "./Inventory.css";
import "../cyberpunk-css-main/cyberpunk.css";
import ItemTableComponent from "../ItemTable";

const Inventory = () => {
  function handleSubmit(e) {
    // Prevent the browser from reloading the page
    e.preventDefault();
  
    // Read the form data
   
    const form = e.target;
    const formData = new FormData(form);
    const formJson = Object.fromEntries(formData.entries());
    
    // You can pass formData as a fetch body directly:
    fetch("http://ronstad.se/Marketplace/CreateItem", { method: "POST", body: JSON.stringify(formJson)})
    .then(response =>  {response.ok ? (alert("ok")):(alert("not ok"))})
    .catch(error => {alert("nuh uh")});
  }
  function handleSubmit2(e) {
    // Prevent the browser from reloading the page
    e.preventDefault();
    // Read the form data
   
    const form = e.target;
    const formData = new FormData(form);
    const formJson = Object.fromEntries(formData.entries());
    
    // You can pass formData as a fetch body directly:
    fetch("http://ronstad.se/Marketplace/AddMoney", { method: "POST", body: JSON.stringify(formJson)})
    .then(response =>  {response.ok ? (alert("ok")):(alert("not ok"))})
    .catch(error => {alert("nuh uh")});
  }
  function changeUID() {
      let userID: string = (
        document.getElementById("id") as HTMLInputElement
      ).value;
      localStorage.setItem("uid", userID)
      alert(localStorage.getItem("uid"))
      window.location.reload();
  }

    return (
      <body>
        <header>My Inventory</header>
    <div className="left-right-container">
      <div className="right">
        {/* store */}
        <input type="text" id="id"></input> <button onClick={() => changeUID()}>Change UID</button>
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
      <form method="post" onSubmit={handleSubmit}>
      <label>
          UserID: <input name="userID" type="text" />
        </label>
        <br></br>
        <label>
          ItemType: <input name="itemType" type="text" />
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
      <form method="post" onSubmit={handleSubmit2}>
      <label>
          UserID: <input name="userID" type="text" />
        </label>
        <br></br>
        <label>
          Amount: <input name="money" type="text" />
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