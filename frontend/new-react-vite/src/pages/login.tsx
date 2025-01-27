import { Outlet, Link } from "react-router";
import "./pageCss/login.css"
function handleSubmit(e) {
    // Prevent the browser from reloading the page
    e.preventDefault();

    // Read the form data
    const form = e.target;
    const formData = new FormData(form);

    // You can pass formData as a fetch body directly:
    fetch('/some-api', { method: form.method, body: formData });

    // Or you can work with it as a plain object:
    const formJson = Object.fromEntries(formData.entries());
    console.log(formJson);
  }
const Login = () => {
    return (
        <pre>
        <h1>Cybermarket</h1>
        <form method="post" onSubmit={handleSubmit}>
      <label>
        Username: <input name="userName" type="text"/>
      </label>
      <br>
      </br>
      <label>
        Password: <input name="userName" type="password"></input>
      </label>
      <hr />
      <button type="submit"><Link to="/Marketplace">Log In</Link></button>
      
    </form>
    </pre>
    );
  };
  
  export default Login;