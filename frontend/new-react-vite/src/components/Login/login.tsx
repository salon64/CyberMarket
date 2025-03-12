import { Link } from "react-router";
import "./login.css";
import CRTScreen from "../CRTScreen";
import { useNavigate } from "react-router";
import { globalAddr } from "../../header";

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
    fetch("http://"+globalAddr+"/login", { method: "POST", body: JSON.stringify(formJson) })
      .then(response => response.json())
      .then(data => {
        if (data?.Token) {  // if token exists
          console.log(data)
          localStorage.setItem("token", data.Token)
          localStorage.setItem("uid", data.UserID)
          localStorage.setItem("role", data.Role)
          navigate("/Marketplace")
        } else {
          alert("nuh uh")
        }
      })
      .catch(error => { alert("testing " + error) });


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
            Login
        </button>
        <br>
        </br>
        <Link to="/Register">Register your account</Link>
      </form>
    </CRTScreen>
  );
};

export default Login;
