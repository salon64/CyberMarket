import "./Marketplace.css";
import "../cyberpunk-css-main/cyberpunk.css";
import CRTScreen from "../CRTScreen";

// I am really cool
function Marketplace() {
  return (
    
  <body>


  <div className="left-right-container">
    <div className="left">
      <div className="left-elements">
        <div className="sorting">
          Sort by:
          <div className="cyber-input">
            <div className="cyber-select ac-purple">
              <select>
                  <option style={{color:"black"}}value="Price">Price</option>
                  <option style={{color:"black"}}value="Newest">Newest</option>
                  <option style={{color:"black"}}value="Oldest">Oldest</option>
              </select>
            </div>
          </div>
        </div>

        <div className="tags">
          <div className="tagtext">
            <header>Tags</header>
            <ul>
              <input type="checkbox" className="cyber-check-tags" ></input> Cyberware <br></br>
              <input type="checkbox" className="cyber-check-tags" ></input> Consumables<br></br>
              <input type="checkbox" className="cyber-check-tags" ></input> Weapons<br></br>
              <input type="checkbox" className="cyber-check-tags" ></input> Drugs<br></br>
              <input type="checkbox" className="cyber-check-tags" ></input> Armor<br></br>
            </ul>
          </div>
        </div>
      </div>
      
    </div>

    <div className="right">
      
      {/* store */}
      <table className="cyber-table store-table">
        <thead>
          <tr className="thead">
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