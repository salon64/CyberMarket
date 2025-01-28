import { Outlet, Link } from "react-router";
import "../Login/login.css";
import CRTScreen from "../CRTScreen";

async function handleSubmit(e) {
    // Prevent the browser from reloading the page
    e.preventDefault();
  
    // Read the form data
    
    const form = e.target;
    const formData = new FormData(form);
    console.log(formData)
    // You can pass formData as a fetch body directly:

  
    // Or you can work with it as a plain object:
    const formJson = Object.fromEntries(formData.entries());
    console.log(JSON.stringify(formJson));
    fetch("http://ronstad.se/user", { method: "POST", body: JSON.stringify(formJson), mode: "no-cors"})
    .then(response => response.text())
    // .then(data => console.log(data));
  }


  const Register = () => {
    return (
      <CRTScreen>
        <h1>Cybermarket</h1>
        <form method="post" onSubmit={handleSubmit}>
          <label>
            Username: <input name="Name" type="text" />
          </label>
          <br></br>
          <label>
            Password: <input name="Pswd" type="password" />
          </label>
          <hr />
          <button type="submit">
          Register
          {/* <Link to="/Marketplace">Log In</Link> */}
          </button>
          <br>
        </br>
        <Link to="/">Return to login</Link>
        </form>
      </CRTScreen>
    );
  };
  
  export default Register;