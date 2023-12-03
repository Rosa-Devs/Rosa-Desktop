import React, { useState, useEffect } from 'react';
import { ManifestList } from "../../wailsjs/go/src/DbManager";
import { manifest } from "../../wailsjs/go/models";
import Buble from '../itemkit/ChatBuble'; // Import the Buble component if not already imported
import { AddManifets } from '../../wailsjs/go/src/DbManager';

const RightSidebar: React.FC = () => {
  const [contacts, setContacts] = useState<manifest.Manifest[] | null>(null);
  const [isUploadPopupOpen, setIsUploadPopupOpen] = useState(false);

  useEffect(() => {
    const fetchContacts = async () => {
      try {
        const contactsData = await ManifestList();
        setContacts(contactsData);
      } catch (error) {
        console.error('Error fetching contacts:', error);
      }
    };

    fetchContacts();
  }, []);

  const handleAddManifest = () => {
    setIsUploadPopupOpen(true);
  };

  const handleCloseUploadPopup = () => {
    setIsUploadPopupOpen(false);
  };


  const handleFileUpload = async (file: File) => {
    // Handle file upload logic here
    // You can use this function to upload the selected file
    try {
      
      const fileContents = await readFile(file);
      console.log('File Data:', fileContents);

      // @ts-ignore
      AddManifets(fileContents)
      window.location.reload();
    } catch (error) {
      console.error('Error reading file:', error);
    }
    // After handling the file, close the popup
    handleCloseUploadPopup();
  };

  const readFile = (file: File) => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = (event: any) => {
        resolve(event.target.result);
      };

      reader.onerror = (error) => {
        reject(error);
      };

      reader.readAsText(file);
    });
  };

  return (
    <div className="sidebar">
      {contacts !== null && contacts.length === 0 ? (
        <p>No rooms!</p>
      ) : (
        <div className="content">
          {contacts?.map((contact) => (
            <Buble key={contact.pubsub} contact={contact} />
          ))}
          <div className="add-button-container">
            <button onClick={handleAddManifest} className="add-button">Add</button>
          </div>
        </div>
      )}

      {isUploadPopupOpen && (
        <div className="upload-popup">
          <div className="upload-popup-content">
            <input type="file" onChange={(e) => handleFileUpload(e.target.files![0])} />
            <button onClick={handleCloseUploadPopup}>Close</button>
          </div>
        </div>
      )}
    </div>
  );
};

export default RightSidebar;
