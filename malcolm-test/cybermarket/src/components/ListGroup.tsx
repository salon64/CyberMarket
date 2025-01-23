// import { MouseEvent  } from "react";

import { useState } from "react";

function ListGroup() {
  const devs = ["Malcolm", "shaya", "Olle"];
  // const handleClick =  (event: MouseEvent) => console.log(event);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  return (
    <>
      <h1>List</h1>
      <ul className="list-group">
        {devs.map((item, index) => (
          <li
            className={
              selectedIndex === index
                ? "list-group-item active"
                : "list-group-item"
            }
            key={item}
            onClick={() => {
              setSelectedIndex(index);
            }}
          >
            {item}
          </li>
        ))}
      </ul>
    </>
  );
}

export default ListGroup;
