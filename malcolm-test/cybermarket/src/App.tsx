import CRTScreen from "./components/CRTScreen";
import "./App.css";

function App() {
  return (
    <>
      <CRTScreen> 
        <h1>CYBERPUNK</h1>
      </CRTScreen>
    </>
  );
}


export default App;




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