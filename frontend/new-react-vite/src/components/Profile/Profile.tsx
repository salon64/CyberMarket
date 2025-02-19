import { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";

import MyImage from "./aswedishtiger.png";
import "../cyberpunk-css-main/cyberpunk.css";
const getToken = () => {
  alert("token is: " + localStorage.getItem("token"))
}


async function handleSubmit(e) {
  // Prevent the browser from reloading the page
  e.preventDefault();

  // Read the form data
  const form = e.target;
  const formData = new FormData(form);
  console.log(formData);

  // You can pass formData as a fetch body directly:
  let userID: string = (
    document.getElementById("nameChangeFormID") as HTMLInputElement
  ).value;
  userID = "http://ronstad.se/users/" + userID;
  console.log(userID);

  // Or you can work with it as a plain object:
  const formJson = Object.fromEntries(formData.entries());
  console.log(JSON.stringify(formJson));
  fetch(userID, { method: "POST", body: JSON.stringify(formJson) })
    .then((response) => {
      if (response.ok === true) {
        console.log("Valid");
        alert("Account has been succesfully updated");
      } else {
        console.log("Invalid");
        alert("nuh uh");
      }
    })
    .then((data) => console.log(data));
  fetch("http://ronstad.se/users", { method: "GET" })
    .then((response) => response.json())
    .then((res) => console.log(res));
}

function Profile() {
  const nameForm = useRef(null);
  const [data, setData] = useState(null);
  let navigate = useNavigate()
  useEffect(() => {
    fetch("http://ronstad.se/users", { method: "GET" })
      .then((response) => response.json())
      .then((data) => setData(data))
      // .then((data) => console.log(data))
      .then(data => console.log(data))
      .catch((error) => console.error("Error: ", error));

    // console.log(JSON.stringify(data));
    console.log("seal was here")
  }, []);


 

  // const handleClickEvent = () => {

  // };
  // function myFunction(event: React.FormEvent<HTMLFormElement>) {
  //   event.preventDefault();
  //   const form = event.currentTarget;
  //   const input = form.elements.namedItem('q') as HTMLInputElement;
  //   console.log(input.value);
  // }
  const [uid] = useState(localStorage.getItem("uid"))
  return (
    <>
      <div>
        <p className="oxanium-font">UID: {uid} </p>
        <br></br>
        {data ? (
          JSON.stringify(data)
        ) : (
          <img src={MyImage} alt="User data"></img>
        )}
      </div>

      

      <div>
        <div className="cyber-input">
          <input
            type="text"
            placeholder="Enter ID..."
            id="nameChangeFormID"
          ></input>
          <form method="post" onSubmit={handleSubmit}>
            <label>
              <input name="token" type="text" placeholder="token..." />
            </label>
            <label>
              <input
                name="new_name"
                type="text"
                placeholder="Enter new username..."
              />
            </label>
            <br></br>
            <label>
              <input
                name="new_pswd"
                type="password"
                placeholder="Enter new password..."
              />
            </label>
            <button>Save changes</button> <br>
            </br>
            <button onClick={getToken}>Check token</button>
            <button onClick={() => {
              localStorage.setItem("token", "")
                navigate('/')
              }}>Log Out</button>
          </form>
        </div>
      </div>

      {/* <div className="cyber-input">
        
      </div>
      <button className="cyber-button bg-* fg-balck" onClick={handleClickEvent}>
        Save changes
        <span className=""></span>
        <span className="tag">User</span>
      </button> */}
    </>
  );
}

export default Profile;
