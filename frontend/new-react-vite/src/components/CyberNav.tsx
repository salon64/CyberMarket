import "./CyberNav.css";
import "../cyberpunk-css-main/cyberpunk.css";
import { Outlet, Link } from "react-router";

function CyberNav() {
    return (
        <div className="top">
        <div className="top-elements">
          <h1 className="cyberpunk-font-og">CyberMarket</h1>
          <div></div> {/*epty*/}
          <nav>
            <ul className="top-navbar">
              <li><Link to="/Marketplace">Marketplace</Link></li>
              <li><Link to="/Inventory">Inventory</Link></li>
              <li><Link to="/Profile">Profile</Link></li>
            </ul>
          </nav>
          <div>Cart</div>
        </div>
      </div>
    );
}

export default CyberNav;