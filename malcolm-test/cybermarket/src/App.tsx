import "./App.css";
// import "././cyberpunk-css-main/cyberpunk.css";

function App() {
  return (
    <div className="tripple">
      <div className="left">left</div> 

      <div className="cyberpunk-window">
        {/* <div className="stripes"> mid</div> */}
        <div className="crt-screen">
          mid
        </div>
      </div>

      <div className="right">right</div>
    </div>

    // <div className="0"> 
    //   <div className="1">text</div>
    //   <div className="2">-----</div>
    // </div>
  );
}

export default App;

// function App() {
//   return (
    
//   <body>
//     <div className="top">
//       <div className="top-elements">
//         <h1 className="cyberpunk-font">CyberMarket</h1>
//         <div></div> {/*epty*/}
//         <nav>
//           <ul className="top-navbar">
//             <li>Market</li>
//             <li>Inventory</li>
//             <li>Profile</li>
//           </ul>
//         </nav>
//         <div>Cart</div>
//       </div>
//     </div>

//   <div className="left-right-container">
//     <div className="left">
//       <div className="left-elements">
//         <div className="sorting">
//           Sort by:
//           <div className="cyber-input">
//             <div className="cyber-select">
//               <select>
//                   <option value="Price">Price</option>
//                   <option value="Newest">Newest</option>
//                   <option value="Oldest">Oldest</option>
//               </select>
//             </div>
//           </div>
//         </div>

//         <div className="tags">
//           Tags
//         </div>
//       </div>
      
//     </div>

//     <div className="right">
//       {/* store */}
//       <table className="cyber-table store-table">
//         <thead>
//           <tr>
//             <th>Header 1</th>
//             <th>Header 2</th>
//             <th>Header 3</th>
//           </tr>
//         </thead>
//         <tbody>
//           <tr>
//             <td>Row 1, Col 1</td>
//             <td>Row 1, Col 2</td>
//             <td>Row 1, Col 3</td>
//           </tr>
//           <tr>
//             <td>Row 2, Col 1</td>
//             <td>Row 2, Col 2</td>
//             <td>Row 2, Col 3</td>
//           </tr>
//           <tr>
//             <td>Row 3, Col 1</td>
//             <td>Row 3, Col 2</td>
//             <td>Row 3, Col 3</td>
//           </tr>
//         </tbody>
//       </table>
//     </div>
//   </div>
// </body>

//   );
// }


// export default App;

// <body>
    //   <header className="cyber-h cyberpunk-font" >
    //     CyberMarket
    //   </header>

    //   <footer>
    //     <div className="cyber-razor-bottom">footer</div>
    //   </footer>
    // </body>


// function App() {
//   return (
//       <div className="crt-screen">
//         <div className="crt-content">
//           <h1>Cyberpunk is like really preem</h1>
//           <p>But retro might be cooler, idk?</p>
//         </div>
//       </div>
//   );
// }

// import Alert from "./components/Alert";
// import ListGroup from "./components/ListGroup";

// function App() {
//   let devs = ["Malcolm", "shaya", "Olle"];

//   const handleSelectItem = (item: string) => {
//     console.log(item)
//   }

//   return <div><ListGroup items={devs} heading="Devs" onSelectItem={handleSelectItem}/></div>;
// }
// function App() {
//   return (
//     <div>
//       <Alert> 
//         Hello <span>World</span>
//       </Alert>
//     </div>
//   );
// }