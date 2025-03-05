import { Navigate, Outlet } from "react-router"
import { globalAddr } from "../../header";
type ProtectedRouteProps = {
  //user: AuthUser | null;
  redirectPath?: string;
  children: React.ReactNode;
}
const ProtectedRoute = ({ redirectPath = "/" }: ProtectedRouteProps) => {

      if (localStorage.getItem("token") === "" || localStorage.getItem("token") === null) {
        alert("invalid token")
        return <Navigate to={redirectPath} replace />
      } else {
        //alert("valid")
        return <Outlet />;
      }
      return <Navigate to={redirectPath} replace />
    
}
export default ProtectedRoute;
