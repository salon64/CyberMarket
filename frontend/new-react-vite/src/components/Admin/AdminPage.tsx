import "../cyberpunk-css-main/cyberpunk.css";
interface addMoney {
  UserID: number;
  Money: number;
}
interface createItemInt {
  UserID: number;
  ItemType: number;
}
interface ItemTypeInformation {
  ItemName: string;
  ItemDescription: string | null;
  ImgURL: string | null;
  ShortDescription: string | null;
}


function AdminPage() {

  function CreateItem(e: any) {
    // Prevent the browser from reloading the page
    e.preventDefault();

    let usItmID: number = (
      document.getElementById("usItmId") as HTMLInputElement
    ).valueAsNumber;
    let itType: number = (
      document.getElementById("itmType") as HTMLInputElement
    ).valueAsNumber;
    let tmp: createItemInt = { UserID: usItmID, ItemType: itType };
    console.log(tmp);
    // You can pass formData as a fetch body directly:
    fetch("http://ronstad.se:5687/Marketplace/CreateItem", {
      method: "POST",
      body: JSON.stringify(tmp),
    })
      .then((response) => {
        response.ok
          ? console.log("Successfully executed request")
          : alert("Invalid input");
        window.location.reload();
      })
      .catch((error) => {
        alert("Error attempting to create an item :" + error);
      });
  }

  function AddMoney(e: any) {
    // Prevent the browser from reloading the page
    e.preventDefault();
    // Read the form data
    let usID: number = (document.getElementById("usId") as HTMLInputElement)
      .valueAsNumber;
    let currAmount: number = (
      document.getElementById("cur") as HTMLInputElement
    ).valueAsNumber;

    let tmp: addMoney = { UserID: usID, Money: currAmount };
    console.log(tmp);

    fetch("http://ronstad.se:5687/user/AddMoney", {
      method: "POST",
      body: JSON.stringify(tmp),
    })
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
      })
      .catch((error) => { console.log(error) }); // kastar error när det funkar?????????????
    window.location.reload();
  }

  function changeUID() {
    let userID: string = (document.getElementById("id") as HTMLInputElement)
      .value;
    localStorage.setItem("uid", userID);
    alert(localStorage.getItem("uid"));
    window.location.reload();
  }

  function CreateItemType(e: any) {
    e.preventDefault();

    let ItemName: string = (
      document.getElementById("ItemName") as HTMLInputElement
    ).value;
    let ShortDescription: string = (
      document.getElementById("ShortDescription") as HTMLInputElement
    ).value;

    let tmp: ItemTypeInformation = { ItemName: ItemName, ItemDescription: null, ImgURL: null, ShortDescription: ShortDescription };
    console.log(tmp);


    fetch("http://ronstad.se:5687/Admin/CreateNewItemType", {
      method: "POST",
      body: JSON.stringify(tmp),
    })
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
      })
      .catch((error) => { console.log(error) }); // kastar error när det funkar?????????????
    window.location.reload();
  }

  return (
    <div>
      <h1>Admin Page</h1>
      <input type="text" id="id"></input>{" "}
      <button onClick={() => changeUID()}>Change UID</button>
      <br />
      <hr />


      <h1>Create Item</h1>
      <form method="post" onSubmit={CreateItem}>
        <label>
          UserID: <input name="userID" type="number" id="usItmId" />
        </label>
        <br></br>
        <label>
          ItemType: <input name="itemType" type="number" id="itmType" />
        </label>
        <br></br>
        <br />
        <button type="submit">Create Item</button>
        <br></br>
        <br />
        <hr />
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
        <br />
        <button type="submit">Add money</button>
        <br></br>
        <br />
        <hr />
      </form>


      <h1>Create Item Type</h1>
      <form method="post" onSubmit={CreateItemType}>
        <label>
          ItemName: <input name="itemname" type="text" id="ItemName" />
        </label>
        <br></br>
        <label>
          ShortDescription: <input name="shortdesc" type="text" id="ShortDescription" />
        </label>
        <br />
        <br></br>
        <button type="submit">create ItemType</button>
        <br></br>
        <br />
        <hr />
      </form>
    </div>
  );
}

export default AdminPage;
