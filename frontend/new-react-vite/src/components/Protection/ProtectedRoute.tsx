import { Navigate } from "react-router"

const ProtectedRoute = () => {
    fetch("http://ronstad.se/auth", { method: "POST", body: JSON.stringify(localStorage.getItem("token"))})
    .then(response => response.json())
    .then(data => {
      console.log(data)
      if (localStorage.getItem("token") !== data) {
        <Navigate to="/" replace />
      }
    }
}