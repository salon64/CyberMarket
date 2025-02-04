import "./Inventory.css";
import "../cyberpunk-css-main/cyberpunk.css";

const Inventory = () => {
    return (
      <body>
        <header>My Inventory</header>
    <div className="left-right-container">
      {/* <div className="left">
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
        
      </div> */}
  
      <div className="right">
        {/* store */}
        <table className="cyber-table store-table">
          <thead>
            <tr>
              <th>Item</th>
              <th>Description</th>
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
  };
  
  export default Inventory;