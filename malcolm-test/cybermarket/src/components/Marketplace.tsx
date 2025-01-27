import "./Marketplace.css";
import "./cyberpunk-css-main/cyberpunk.css";


function Marketplace() {
  return (
    
  <body>
    <div className="top">
      <div className="top-elements">
        <h1 className="cyberpunk-font">CyberMarket</h1>
        <div></div> {/*epty*/}
        <nav>
          <ul className="top-navbar">
            <li>Market</li>
            <li>Inventory</li>
            <li>Profile</li>
          </ul>
        </nav>
        <div>Cart</div>
      </div>
    </div>

  <div className="left-right-container">
    <div className="left">
      <div className="left-elements">
        <div className="sorting">
          Sort by:
          <div className="cyber-input">
            <div className="cyber-select">
              <select>
                  <option value="Price">Price</option>
                  <option value="Newest">Newest</option>
                  <option value="Oldest">Oldest</option>
              </select>
            </div>
          </div>
        </div>

        <div className="tags">
          Tags
        </div>
      </div>
      
    </div>

    <div className="right">
      {/* store */}
      <table className="cyber-table store-table">
        <thead>
          <tr>
            <th>Header 1</th>
            <th>Header 2</th>
            <th>Header 3</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>Row 1, Col 1</td>
            <td>Row 1, Col 2</td>
            <td>Row 1, Col 3</td>
          </tr>
          <tr>
            <td>Row 2, Col 1</td>
            <td>Row 2, Col 2</td>
            <td>Row 2, Col 3</td>
          </tr>
          <tr>
            <td>Row 3, Col 1</td>
            <td>Row 3, Col 2</td>
            <td>Row 3, Col 3</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</body>

  );
}


export default Marketplace;