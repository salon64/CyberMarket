import { Outlet, Link } from "react-router";
import CyberNav from "./CyberNav";

const Layout = () => {
  return (
    <>
      <CyberNav></CyberNav>
      <Outlet />
    </>
  )
};

export default Layout;