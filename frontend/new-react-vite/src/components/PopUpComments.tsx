import { useEffect, useState } from "react";
import { X, Trash2 } from "lucide-react";
import { globalAddr } from "../header.tsx";

interface PopUpCommentsProps {
  onClose: () => void;
  itemId: number | null;
}

interface ItemTypeReturn {
  Name: string;
  ImgURL: string;
  ShortDesc: string;
  DescURL: string;
  Comments: PubComment[];
}

interface PubComment {
  CommentID: number;
  UserName: string;
  UserID: number;
  Grade: number;
  Comment: string;
  PostedOn: string;
  ParentCommentID: number;
}

interface MakeComment {
  UserID: number;
  Grade: number;
  Comment: string;
  ParentCommentID: number | null;
}

interface fml {
  theComment: PubComment;
  subComments: fml[];
}

const PopUpComments: React.FC<PopUpCommentsProps> = ({ onClose, itemId }) => {
  const [itemTypeReturn, setItemTypeReturn] = useState<ItemTypeReturn | null>(
    null
  );
  const [num, setNum] = useState<number>(0);

  useEffect(() => {
    const handleEscapeKey = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        onClose();
      }
    };
    document.addEventListener("keydown", handleEscapeKey);
    return () => {
      document.removeEventListener("keydown", handleEscapeKey);
    };
  }, [onClose]);

  const listComments = () => {
    // e.preventDefault()
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
  }, [itemId, num]);

  if (!itemTypeReturn) {
    return <div>Loading...</div>;
  }

  function MakeReview1(parentID: number) {
    // e.preventDefault();

    let rating: number = (document.getElementById(`rating${parentID}`) as HTMLInputElement)
      .valueAsNumber;
    let comment: string = (
      document.getElementById(`comment${parentID}`) as HTMLInputElement
    ).value;

    let tmp: MakeComment = {
      UserID: Number(localStorage.getItem("uid")),
      Grade: rating,
      Comment: comment,
      ParentCommentID: parentID,
    };
    // console.log(tmp);

    const fetchString = `http://${globalAddr}/ItemType/${itemId}`;
    fetch(fetchString, {
      method: "POST",
      body: JSON.stringify(tmp),
    })
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
      })
      .catch((error) => {
        console.log(error);
      }); // kastar error när det funkar?????????????
    setNum(num + 1);
  }

  function MakeReview2(parentID: number | null) {
    // e.preventDefault();

    let rating: number = (document.getElementById("rating2") as HTMLInputElement)
      .valueAsNumber;
    let comment: string = (
      document.getElementById("comment2") as HTMLInputElement
    ).value;

    let tmp: MakeComment = {
      UserID: Number(localStorage.getItem("uid")),
      Grade: rating,
      Comment: comment,
      ParentCommentID: parentID,
    };
    // console.log(tmp);

    const fetchString = `http://${globalAddr}/ItemType/${itemId}`;
    fetch(fetchString, {
      method: "POST",
      body: JSON.stringify(tmp),
    })
      .then((response) => response.json())
      .then((data) => {
        console.log(data);
      })
      .catch((error) => {
        console.log(error);
      }); // kastar error när det funkar?????????????
    setNum(num + 1);
  }

  function DeleteComment(num: number) {
    // e.preventDefault();
    const fetchString = `http://${globalAddr}/comment/deletecomment/${num}`;
    console.log(fetchString)
    fetch(fetchString, {
      method: 'GET',
      headers: new Headers({
        "Authorization": "Bearer " + localStorage.getItem("token")
      }),
    })
      .then((response) => console.log(response))
      .catch((error) => console.error("Error: ", error));
    setNum(num + 1);
  }

  function grow(commentList: PubComment[]) {
    // console.log(itemTypeReturn?.Comments);
    const commentDict = new Map<number, fml>();

    // initiates them all with their key: id, value: itself, emptyList
    commentList.forEach((comment) => {
      commentDict.set(comment.CommentID, {
        theComment: comment,
        subComments: [],
      });
    });

    // initiates empty roots
    const rootComments: fml[] = [];

    commentList.forEach((comment) => {
      if (comment.ParentCommentID === null) {
        // add the struct to the list of rootstructs
        rootComments.push(commentDict.get(comment.CommentID)!);
      }
    });

    commentList.forEach((comment) => {
      if (comment.ParentCommentID !== null) {
        // get parent comment
        const parentComment = commentDict.get(comment.ParentCommentID!);
        // get the child
        const currentComment = commentDict.get(comment.CommentID)!;
        if (parentComment) {
          parentComment.subComments.push(currentComment);
        }
      }
    });
    return rootComments;
  }

  const CommentComponent: React.FC<{ commentStruct: fml; depth: number }> = ({
    commentStruct,
    depth,
  }) => {
    return (
      <div
        className="ml-4 border-l pl-2"
        style={{ marginLeft: `${depth * 40}px` }}
      >
        {commentStruct.theComment.UserID ===
        Number(localStorage.getItem("uid")) ? (
          <>
            <div className="display: inline-block;">
              <p>
                <strong>{commentStruct.theComment.UserName}</strong>:{" "}
                {commentStruct.theComment.Comment}
              </p>
              <p>
                <strong>Rating:</strong>: {commentStruct.theComment.Grade}
              </p>
              <Trash2
                style={{ cursor: "pointer" }}
                onClick={() =>
                  DeleteComment(commentStruct.theComment.CommentID)
                }
                size={30}
              />
            </div>
          </>
        ) : (
          <div className="display: inline-block;">
            <p>
              <strong>{commentStruct.theComment.UserName}</strong>:{" "}
              {commentStruct.theComment.Comment}
            </p>
            <p>
              <strong>Rating:</strong>: {commentStruct.theComment.Grade}
            </p>
          </div>
        )}
        <form method="post" onSubmit={(event) => {MakeReview1(commentStruct.theComment.CommentID); event.preventDefault();} }>
        <label>
          rating <input name="rating" type="number" id={`rating${commentStruct.theComment.CommentID}`}  />
        </label>
        <label>
          Comment: <input name="comment" type="text" id={`comment${commentStruct.theComment.CommentID}`} />
        </label>
        <br />
        <button type="submit">Submit Comment</button>
        <br></br>
        <br />
      </form>
        {commentStruct.subComments && (
          <div>
            {commentStruct.subComments.map((sub) => (
              <CommentComponent
                key={sub.theComment.CommentID}
                commentStruct={sub}
                depth={depth + 1} // Increase depth for nested comments
              />
            ))}
          </div>
        )}
      </div>
    );
  };

  return (
    <div>
      <X style={{ cursor: "pointer" }} onClick={onClose} size={30} />
      <h2>Reviews</h2>
      <div>
        {itemTypeReturn.Comments === null ? (
          <p>No comments on this itemType yet</p>
        ) : (
          // <p>tmp</p>
          <>
            {grow(itemTypeReturn.Comments).map((com) => (
              <CommentComponent
                key={com.theComment.CommentID}
                commentStruct={com}
                depth={0}
              />
            ))}
          </>
        )}
      </div>
      <h2>Review ItemType</h2>
      <form method="post" onSubmit={(event) => {MakeReview2(null); event.preventDefault();} }>
        <label>
          rating <input name="rating" type="number" id="rating2" />
        </label>
        <br></br>
        <label>
          Comment: <input name="comment" type="text" id="comment2" />
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
