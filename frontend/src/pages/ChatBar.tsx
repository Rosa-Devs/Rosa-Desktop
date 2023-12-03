import { useEffect, useState } from "react";
import { AiOutlineHome, AiOutlineSetting, AiOutlineInfoCircle } from 'react-icons/ai';
import Buble from "../itemkit/ChatBuble";
import Contact from "../itemkit/Contact";
import { GetContacts} from "../../wailsjs/go/main/App"


// const contacts: Contact[] = [
//   { id: 1, name: 'John Doe', imageUrl: 'https://cdn3.iconfinder.com/data/icons/shinysocialball/512/Technorati_512x512.png' },
//   { id: 2, name: 'Jane Smith', imageUrl: 'https://ccia.ugr.es/cvg/CG/images/base/5.gif' },
//   { id: 3, name: 'Alice Johnson', imageUrl: 'https://upload.wikimedia.org/wikipedia/commons/c/cc/Icon_Pinguin_1_512x512.png' },
//   // Add more contacts as needed
// ];

const RightSidebar: React.FC = () => {

  const [contacts, setContacts] = useState<Contact[]>([]);

  useEffect(() => {
    const fetchContacts = async () => {
      try {
        const contactsData = await GetContacts();
        setContacts(contactsData);
      } catch (error) {
        console.error('Error fetching contacts:', error);
      }
    };

    fetchContacts();
  }, []);

  return (
    <div className="sidebar">
      {contacts.map((contact) => (
          <Buble key={contact.id} contact={contact} />
        ))}
    </div>
  );
};

export default RightSidebar;