import { MouseEventHandler } from "react";
import Contact from "./Contact";







const Buble: React.FC<{ contact: Contact }> = ({ contact }) => {

    const handle_btn = () => {
        console.log(contact.name)
    }

    return (
        <button onClick={handle_btn} className="sidebar-btn">
        <img src={contact.imageUrl} alt={contact.name} className="img-icon" />
        </button>
    );
};


export default Buble