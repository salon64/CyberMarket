import { Outlet, Link } from "react-router";

const Layout = () => {
  return (
    <>
      <nav>
        
          
            <Link to="/">Login</Link>
          
            <Link to="/Marketplace">Marketplace</Link>
     
            <Link to="/contact">Contact</Link>
        
    
      </nav>

      <Outlet />
    </>
  )
};

export default Layout;