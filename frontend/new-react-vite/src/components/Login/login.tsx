import { Outlet, Link } from "react-router";
import { useState } from "react";
import "./login.css";
import CRTScreen from "../CRTScreen";

function handleSubmit(e) {
  // Prevent the browser from reloading the page
  e.preventDefault();

  // Read the form data
  const form = e.target;
  const formData = new FormData(form);
  
  // You can pass formData as a fetch body directly:
  fetch("http://ronstad.se/users", { method: "GET"})
  .then(response => response.json())
  .then(data => console.log(data));

  // Or you can work with it as a plain object:
  const formJson = Object.fromEntries(formData.entries());
  console.log(formJson);
}
function handleSubmit2(e) {
  // Prevent the browser from reloading the page
  e.preventDefault();

  // Read the form data
  const form = e.target;
  const formData = new FormData(form);
  
  // You can pass formData as a fetch body directly:
  fetch("http://ronstad.se/users", { method: "GET"})
  .then(response => response.json())
  .then(data => console.log(data));

  // Or you can work with it as a plain object:
  const formJson = Object.fromEntries(formData.entries());
  console.log(formJson);
}
const Login = () => {
  const [logButton, setLogButton] = useState("login");
  return (
    <CRTScreen>
      <h1 className="cyberpunk-font-og">Cybermarket</h1>
      <form method="post" onSubmit={handleSubmit}>
        <label>
          Username: <input name="userName" type="text" />
        </label>
        <br></br>
        <label>
          Password: <input name="password" type="password" />
        </label>
        <hr />
        <button type="submit">
        Login
        {/* <Link to="/Marketplace">Log In</Link> */}
        </button>
        <br>
        </br>
        <Link to="/Register">Register your account</Link>
      </form>
    </CRTScreen>
  );
};

export default Login;
