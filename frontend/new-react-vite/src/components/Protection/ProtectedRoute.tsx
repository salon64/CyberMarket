import { Navigate, Outlet } from "react-router"
import { globalAddr } from "../../header";
type ProtectedRouteProps = {
  //user: AuthUser | null;
  redirectPath?: string;
  children: React.ReactNode;
}
const ProtectedRoute = ({ redirectPath = "/" }: ProtectedRouteProps) => {
  fetch("http://"+globalAddr+"/auth", { method: "POST", body: JSON.stringify(localStorage.getItem("token")) })
    .then(response => response.json())
    .then(data => {
      console.log(data)
      if (data !== localStorage.getItem("token")) {
        alert("invalid token")
        return <Navigate to={redirectPath} replace />
      } else {
        return <Outlet />;
      }
    })
    .catch(error => {
      console.log(error)
      return <Navigate to={redirectPath} replace />
    })
  return <Outlet />;
}
export default ProtectedRoute;
