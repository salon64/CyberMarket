import { useEffect, useState } from "react";
import { X, Trash2 } from 'lucide-react';
import { globalAddr } from "../header.tsx";

interface PopUpCommentsProps {
  onClose: () => void;
  itemId: number | null; // Add itemId as a prop
}

interface ItemTypeReturn {
  Name:      string;
  ImgURL:    string;
  ShortDesc: string;
  DescURL:   string;
  Comments:  PubComment[];
}

interface PubComment {
  CommentID:number;
  UserName: string;
  UserID:   number;
  Grade:    number;
  Comment:  string;
  PostedOn: string;
}

interface MakeComment {
  UserID:   number;
  Grade:    number;
  Comment:  string;
}

const PopUpComments: React.FC<PopUpCommentsProps> = ({ onClose, itemId }) => {
  const [itemTypeReturn, setItemTypeReturn] = useState<ItemTypeReturn | null>(null);

  useEffect(() => {
    const handleEscapeKey = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        onClose();
      }
    };
    document.addEventListener('keydown', handleEscapeKey);
    return () => {
      document.removeEventListener('keydown', handleEscapeKey);
    };
  }, [onClose]);

  const listComments = () => {
    if (itemId === null) return; // Handle case when itemId is null
    const fetchString = `http://${globalAddr}/ItemType/${itemId}`;
    fetch(fetchString, { method: "GET" })
      .then((response) => response.json())
      .then((itemTypeReturn) => {
        setItemTypeReturn(itemTypeReturn);
      })
      .catch((error) => console.error("Error: ", error));
  };

  // Call listComments when component mounts to fetch item details
  useEffect(() => {
    listComments();
  }, [itemId]);

  if (!itemTypeReturn) {
    return <div>Loading...</div>;
  }

  function MakeReview(e: any) {
      e.preventDefault();
  
      let rating: number = (
        document.getElementById("rating") as HTMLInputElement
      ).valueAsNumber;
      let comment: string = (
        document.getElementById("comment") as HTMLInputElement
      ).value;
      
      let tmp: MakeComment = {UserID: Number(localStorage.getItem("uid")), Grade: rating, Comment: comment };
      console.log(tmp);
  
      const fetchString = `http://${globalAddr}/ItemType/${itemId}`;
      fetch(fetchString, {
        method: "POST",
        body: JSON.stringify(tmp),
      })
        .then((response) => response.json())
        .then((data) => {
          console.log(data);
        })
        .catch((error) => { console.log(error) }); // kastar error när det funkar?????????????
      window.location.reload();
  }

  function DeleteComment(num: number) {
    // e.preventDefault();
    const fetchString = `http://${globalAddr}/comment/deletecomment/${num}`;
    console.log(fetchString)
    fetch(fetchString, {method: 'GET'})
      .then((response) => console.log(response))
      .catch((error) => console.error("Error: ", error))
  };

  return (
    <div>
      <X style={{ cursor: 'pointer' }} onClick={onClose} size={30} />
      <h2>Reviews</h2>
      <h3>ItemType rating...</h3>
      <div>
        {itemTypeReturn.Comments === null ? (
          <p>No comments on this itemType yet</p>
        ) : (
          itemTypeReturn.Comments.map((comment) => (
            <div key={comment.CommentID}>
              {comment.UserID === Number(localStorage.getItem("uid")) ? (
                <>
                  <div className="display: inline-block;">
                    <strong>{comment.UserName}</strong>
                    <Trash2 style={{ cursor: 'pointer' }} onClick={() => DeleteComment(comment.CommentID)} size={30} />
                  </div>
                </>
              ) : (
                <>
                <strong>{comment.UserName}</strong>
                </>
              )}
              <p>{comment.Comment}</p>
              <span>{comment.Grade} stars</span>
              <p>{new Date(comment.PostedOn).toLocaleString()}</p>
              <hr />
            </div>
          ))
        )}
      </div>

      <h2>Review ItemType</h2>
      <form method="post" onSubmit={MakeReview}>
        <label>
          rating <input name="rating" type="number" id="rating" />
        </label>
        <br></br>
        <label>
          Comment: <input name="comment" type="text" id="comment" />
        </label>
        <br></br>
        <br />
        <button type="submit">Submit Comment</button>
        <br></br>
        <br />
        <hr />
      </form>
    </div>
  );
};

export default PopUpComments;










// import { useEffect, useState } from "react";
// import { X } from 'lucide-react';
// import { globalAddr } from "../header.tsx"

// interface PopUpCommentsProps {
//   onClose: () => void;
//   itemId: number | null; // Add itemId as a prop
// }

// interface ItemTypeReturn {
//   Name:      string;
//   ImgURL:    string;
//   ShortDesc: string;
//   DescURL:   string;
//   Comments:  PubComment[];
// }

// interface PubComment {
//   UserName: string;
//   UserID:   number;
//   Grade:    number;
//   Comment:  string;
//   PostedOn: string;
// }

// const PopUpComments: React.FC<PopUpCommentsProps> = ({ onClose, itemId }) => {
//   const [itemTypeReturn, setitemTypeReturn] = useState<ItemTypeReturn | null>(null);


//   useEffect(() => {
//     const handleEscapeKey = (event: KeyboardEvent) => {
//       if (event.key === 'Escape') {
//         onClose();
//       }
//     };
//     document.addEventListener('keydown', handleEscapeKey);
//     return () => {
//       document.removeEventListener('keydown', handleEscapeKey);
//     };
//   }, [onClose]);

//   const getComments = () => {
//       var fetchString = `http://`+globalAddr+`/ItemType/` + itemId
//       fetch(fetchString, { method: "GET" }) 
//         .then((response) => response.json())
//         .then((itemTypeReturn) => {setitemTypeReturn(itemTypeReturn)})
//         .catch((error) => console.error("Error: ", error));
//   };

//   return (
//     itemTypeReturn.Comments.map
//     // <div>
//     //   <X style={{ cursor: 'pointer' }} onClick={onClose} size={30}></X>
//     //   <br />

//     // </div>
//   );
// };

// export default PopUpComments;