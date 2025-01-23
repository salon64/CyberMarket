import Alert from "./components/Alert";
import ListGroup from "./components/ListGroup";

// function App() {
//   let devs = ["Malcolm", "shaya", "Olle"];

//   const handleSelectItem = (item: string) => {
//     console.log(item)
//   }

//   return <div><ListGroup items={devs} heading="Devs" onSelectItem={handleSelectItem}/></div>;
// }
function App() {
  return (
    <div>
      <Alert> 
        Hello <span>World</span>
      </Alert>
    </div>
  );
}

export default App;
