import { useEffect, useState } from "react";

interface UserItemInterface {
    ItemID: number;
    TypeID: number;
    ItemName: string;
    ItemDescription: string;
    ImgURL: string;
  }

const ItemTableComponent = () => {
    const [userItems, setUserItems] = useState<UserItemInterface[]>([]);
    useEffect(() => {
        fetch("http://ronstad.se/inventory/1", { method: "GET" }) // Replace with your actual API URL
            .then((response) => response.json())
            .then((useritems) => setUserItems(useritems))
            .catch((error) => console.error("Error: ", error));
        }, []);

      
    return (
        <>
            {userItems.map((item) => (
                <tr key={item.ItemID} className="">
                    <td className="">{item.ItemID}</td>
                    <td className="">{item.TypeID}</td>
                    <td className="">{item.ItemName}</td>
                </tr>
            ))}
        </>
    );
}

export default ItemTableComponent