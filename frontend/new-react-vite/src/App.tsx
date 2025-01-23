import { useState } from 'react'
import validator from 'validator'
//import reactLogo from './assets/react.svg'
//import viteLogo from '/vite.svg'
import './App.css'

const App = () => {
  const [emailError, setEmailError] = useState("");
  const validateEmail = (e) => {
      const email = e.target.value;

      if (validator.isEmail(email)) {
          setEmailError(" ");
      } else {
          setEmailError("Enter valid Email!");
      }
  };

  return (

      <div style={{
        margin: "auto",
        display: "block",
      }}>
          <pre>
              <h1>Cybermarket</h1>
              <span style={{color:"aqua"}}>Enter Email: </span>
              <input
                  type="text"
                  id="userEmail"
                  onChange={(e) => validateEmail(e)}
              ></input>{" "}
              
              <span
                  style={{
                      fontFamily: "Rajdhani",
                      
                      color: "red"
                  }}
              >
                  {emailError}
                  <br />
              </span>
              <span style={{color:"aqua"}}>Enter Password: </span>
              <input
                  
                  type="password"
                  id="userEmail"
                  onChange={(e) => validateEmail(e)}
              ></input>{" "}
              
              <span
                  style={{
                      fontFamily: "Rajdhani",
                      
                      color: "red"
                  }}
              >
                  {emailError}
                  <br />
              </span>
              <br />
              <button>
                Login
              </button>
          </pre>
      </div>
  );
};

export default App;
