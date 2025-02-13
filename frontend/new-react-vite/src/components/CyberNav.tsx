import "./CyberNav.css";
import "./cyberpunk-css-main/cyberpunk.css";
import {  Link } from "react-router";

function CyberNav() {
    return (
        <div className="top">
        <div className="top-elements">
          <h1 className="cyberpunk-font-og">CyberMarket</h1>
          
          <nav>
            <ul className="top-navbar">
              <li><Link to="/Marketplace"><button className="cyber-button-small bg-blue fg-yellow">Marketplace</button></Link></li>
              <li><Link to="/Inventory"><button className="cyber-button-small bg-blue fg-yellow">Inventory</button></Link></li>
              <li><Link to="/Profile"><button className="cyber-button-small bg-blue fg-yellow">Profile</button></Link></li>
            </ul>
          </nav>
          <nav className="top-navbar">
           <li><Link to="/cart"><button className="cyber-button-small bg-red fg-yellow">Cart</button></Link></li>
          </nav>
        </div>
      </div>
    );
}

export default CyberNav;