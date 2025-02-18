import { Navigate, Outlet } from "react-router"
type ProtectedRouteProps = {
  //user: AuthUser | null;
  redirectPath?: string;
  children: React.ReactNode;
}
const ProtectedRoute = ({ children, redirectPath = "/" }: ProtectedRouteProps) => {
    fetch("http://ronstad.se/auth", { method: "POST", body: JSON.stringify(localStorage.getItem("token"))})
    .then(response => response.json())
    .then(data => {console.log(data)
    if (data !== localStorage.getItem("token")) {
      alert("invalid token")
      return <Navigate to={redirectPath} replace />
    } else {
      return <Outlet />;
    }})
    .catch(error => {return <Navigate to={redirectPath} replace />})
    return <Outlet />;
    }
export default ProtectedRoute;