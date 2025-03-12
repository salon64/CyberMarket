import { Navigate, Outlet } from "react-router"
type ProtectedRouteProps = {
  //user: AuthUser | null;
  redirectPath?: string;
  children: React.ReactNode;
  roles?: string;
}
const ProtectedRoute = ({ redirectPath = "/", roles = "0"}: ProtectedRouteProps) => {
      if (roles === "0") {
        if (localStorage.getItem("token") === "" || localStorage.getItem("token") === null) {
          alert("invalid token")
          return <Navigate to={redirectPath} replace />
        } else {
          //alert("valid")
          return <Outlet />;
        }
      }
      else if (roles === "1") {
        if (localStorage.getItem("role") === "1") {
          return <Outlet />;
        } else {
          alert("You don't have role priviliges to view this page")
          return <Navigate to={"/Marketplace"} replace />
        }
      }

}
export default ProtectedRoute;
