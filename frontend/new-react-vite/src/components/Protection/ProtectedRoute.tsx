import { Navigate, Outlet } from "react-router"
type ProtectedRouteProps = {
  user: AuthUser | null;
  redirectPath?: string;
  children: React.ReactNode;
}
const ProtectedRoute = ({user, children, redirectPath = "/" }: ProtectedRouteProps) => {
    fetch("http://ronstad.se/auth", { method: "POST", body: JSON.stringify(localStorage.getItem("token"))})
    .then(response => response.json())
    .then(data => {console.log(data);
    if ("data" === localStorage.getItem("token")) {
      return <Navigate to={redirectPath} replace />;
    }})
    
    return children ? children : <Outlet />;
    }
export default ProtectedRoute;