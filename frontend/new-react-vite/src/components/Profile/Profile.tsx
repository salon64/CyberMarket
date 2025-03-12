import { useEffect, useState } from "react";
import "../cyberpunk-css-main/cyberpunk.css";
import { globalAddr } from "../../header";
import "./Profile.css";

interface TransactionLog {
  TransID: number;
  Price: number;
  Date: string;
  ItemID: number;
  Buyer: number;
  Seller: number;
}

interface updateUserInfoDataName {
  new_name: string | null;
  // new_pswd: string | null;
}

interface updateUserInfoDataPswd {
  // new_name: string | null;
  new_pswd: string | null;
}


function Profile() {
  const [transactionLog, setTransactionLog] = useState<TransactionLog[]>([]);

  useEffect(() => {
    var fetchString =
      `http://${globalAddr}/displayTransactionslog/` +
      localStorage.getItem("uid");
    fetch(fetchString, { method: "GET" })
      .then((response) => response.json())
      .then((transactionLog) => {
        setTransactionLog(transactionLog);
      })
      .catch((error) => console.log("Error: ", error));
  }, []);

  const uid = localStorage.getItem("uid");

  function changeUsername(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    let newUserID: string = (
      document.getElementById("newUsername") as HTMLInputElement
    ).value;
    let tmp: updateUserInfoDataName = {
      new_name: newUserID,
    }
    var fetchString =
      `http://${globalAddr}/users/` +
      localStorage.getItem("uid");
    fetch(fetchString, {method: "POST",
      headers: new Headers({
        "Authorization": "Bearer " + localStorage.getItem("token")
      }),
      body: JSON.stringify(tmp)
    })
    .then((response) => console.log(response))
    .catch((error) => {
      console.log(error);
    });
  }

  function changePassword(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()

    let newPswd: string = (
      document.getElementById("newUsername") as HTMLInputElement
    ).value;
    let tmp: updateUserInfoDataPswd = {
      new_pswd: newPswd
    }
    var fetchString =
      `http://${globalAddr}/users/` +
      localStorage.getItem("uid");
    fetch(fetchString, {method: "POST",
      headers: new Headers({
        "Authorization": "Bearer " + localStorage.getItem("token")
      }),
      body: JSON.stringify(tmp)
    })
    .then((response) => console.log(response))
    .catch((error) => {
      console.log(error);
    });
  }

  return (
      <div className="left-right-container-profile">
        <div className=""></div>
        <div className="left-profile">
          <h2>You are logged in as: {String(uid)}</h2>
          <h1>Change Username</h1>
          <form method="post" onSubmit={changeUsername}>
            <label>
              New Username:{" "}
              <input name="newUsername" type="text" id="newUsername" />
            </label>
            <br></br>
            <br />
            <button type="submit">Submit Username</button>
            <br></br>
            <br />
            <hr />
          </form>
          <h1>Change Password</h1>
          <form method="post" onSubmit={changePassword}>
            <label>
              New Password:{" "}
              <input name="newPassword" type="text" id="newPassword" />
            </label>
            <br></br>
            <br />
            <button type="submit">Submit Password</button>
            <br></br>
            <br />
            <hr />
          </form>
        </div>

        <div className="right-profile">
          <table className="cyber-table store-table ac-custom">
            <thead>
              <tr className="thead">
                <th>TransID</th>
                <th>Price</th>
                <th>Date</th>
                <th>ItemID</th>
                <th>Buyer</th>
                <th>Seller</th>
              </tr>
            </thead>
            <tbody>
              {transactionLog && transactionLog.length > 0 ? (
                transactionLog.map((transaction) => (
                  <tr key={transaction.TransID}>
                    <td>{transaction.TransID}</td>
                    <td>{transaction.Price}</td>
                    <td>{transaction.Date}</td>
                    <td>{transaction.ItemID}</td>
                    <td>{transaction.Buyer}</td>
                    <td>{transaction.Seller}</td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={6}>No transactions available</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
  );
}

export default Profile;
