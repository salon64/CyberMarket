import React, { useState, useEffect, Component, useRef } from "react";
import MyImage from "./aswedishtiger.png";
import "../cyberpunk-css-main/cyberpunk.css";

function Profile() {
  const nameForm = useRef(null);
  const [data, setData] = useState(null);

  useEffect(() => {
    fetch("http://ronstad.se/users", { method: "GET" })
      .then((response) => response.json())
      .then((data) => setData(data))
      .then((data) => console.log(data))
      .catch((error) => console.error("Error: ", error));

      var tmp = JSON.stringify(data)
      console.log(tmp)
    }, []);

  const handleClickEvent = () => {
    let newName = (document.getElementById("nameChangeFormID") as HTMLInputElement).value;
    console.log(newName);
    
  };
  // function myFunction(event: React.FormEvent<HTMLFormElement>) { 
  //   event.preventDefault(); 
  //   const form = event.currentTarget; 
  //   const input = form.elements.namedItem('q') as HTMLInputElement;
  //   console.log(input.value); 
  // }

  return (
    <>
      <div>
        {data ? (
          JSON.stringify(data)
        ) : (
          <img src={MyImage} alt="User data"></img>
        )}
      </div>

      <div>
        <form ref={nameForm}>
          <div className="cyber-input">
            <input type="text" placeholder="New name..." id="nameChangeFormID"/>
          </div>
        </form>
        <button onClick={handleClickEvent}>Save changes</button>
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
