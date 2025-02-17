import { Link } from "react-router";
import "./login.css";
import CRTScreen from "../CRTScreen";
interface uLogin {
  uID: string;
  tok: string;
}
function handleSubmit(e) {
  // Prevent the browser from reloading the page
  e.preventDefault();

  // Read the form data
  const form = e.target;
  const formData = new FormData(form);
  const formJson = Object.fromEntries(formData.entries());
  
  // You can pass formData as a fetch body directly:
  fetch("http://ronstad.se/login", { method: "POST", body: JSON.stringify(formJson)})
  .then(response => response.json())
  .then(data => {
    const obj = JSON.parse(JSON.stringify(data))
    //let obj: uLogin = JSON.parse(JSON.stringify(data))
    console.log(JSON.stringify(data))
    console.log(obj)
    localStorage.setItem("token", obj.Token)
    localStorage.setItem("uid", obj.UserID)
    alert("Token is: "+ localStorage.getItem("token"));
    alert("uid is: "+ localStorage.getItem("uid"))});

  // Or you can work with it as a plain object:

  console.log(formJson);
}

const Login = () => {
  return (
    <CRTScreen>
      <h1 className="cyberpunk-font-og">Cybermarket</h1>
      <form method="post" onSubmit={handleSubmit}>
        <label>
          Username: <input name="name" type="text" />
        </label>
        <br></br>
        <label>
          Password: <input name="pswd" type="password" />
        </label>
        <hr />
        <button type="submit">
        
        </button>
        <br>
        </br>
        <Link to="/Register">Register your account</Link>
      </form>
    </CRTScreen>
  );
};

export default Login;
