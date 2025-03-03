import { Link } from "react-router";
import "./login.css";
import CRTScreen from "../CRTScreen";
import { useNavigate } from "react-router";
import { globAddr } from "../../header";

function Login() {
  const navigate = useNavigate()
  function handleSubmit(e: any) {
    // Prevent the browser from reloading the page
    e.preventDefault();
  
    // Read the form data
   
    const form = e.target;
    const formData = new FormData(form);
    const formJson = Object.fromEntries(formData.entries());
    
    // You can pass formData as a fetch body directly:
    fetch("http://"+globAddr+"/login", { method: "POST", body: JSON.stringify(formJson)})
    .then(response => response.json())
    .then(data => {
      if (data?.Token) {  // if token exists
        const obj = JSON.parse(JSON.stringify(data))
        console.log(JSON.stringify(data))
        console.log(obj)
        localStorage.setItem("token", obj.Token)
        localStorage.setItem("uid", obj.UserID)
        navigate("/Marketplace")
      } else{
        alert("nuh uh")
      }
    })
    .catch(error => {alert(error)});
      
  
    // Or you can work with it as a plain object:
  
    console.log(formJson);
  }
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
