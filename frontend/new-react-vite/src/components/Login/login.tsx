import { Link } from "react-router";
import "./login.css";
import CRTScreen from "../CRTScreen";

function handleSubmit(e) {
  // Prevent the browser from reloading the page
  e.preventDefault();

  // Read the form data
  const form = e.target;
  const formData = new FormData(form);
  const formJson = Object.fromEntries(formData.entries());
  
  // You can pass formData as a fetch body directly:
  fetch("http://ronstad.se/auth", { method: "POST", body: JSON.stringify(formJson)})
  .then(response => response.json())
  .then(data => {
    console.log(data)
    localStorage.setItem("token", JSON.stringify(data))
    alert("Token is: "+ localStorage.getItem("token"))});
    

  // Or you can work with it as a plain object:

  console.log(formJson);
}

const Login = () => {
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
        <Link to="/Marketplace">Log In</Link>
        </button>
        <br>
        </br>
        <Link to="/Register">Register your account</Link>
      </form>
    </CRTScreen>
  );
};

export default Login;
