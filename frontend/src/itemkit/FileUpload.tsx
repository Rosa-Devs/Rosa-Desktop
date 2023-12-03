// FileUploadPopup.jsx
import React, { useState } from "react";
import { AiOutlineClose } from 'react-icons/ai';


// @ts-ignore
const FileUploadPopup = ({ isOpen, onClose, onFileUpload }) => {
  const [selectedFile, setSelectedFile] = useState(null);

  const handleFileChange = (event: any) => {
    const file = event.target.files[0];
    setSelectedFile(file);
  };

  const handleUpload = () => {
    if (selectedFile) {
      onFileUpload(selectedFile);
      onClose();
    }
  };

  return (
  
    <div className={`file-upload-popup ${isOpen ? "block" : "hidden"}`}>
      <div className="overlay" onClick={onClose}></div>
      <div className="popup">
        <div className="header">
          <h2>Upload File</h2>
          <button onClick={onClose}>
            <AiOutlineClose />
          </button>
        </div>
        <input type="file" onChange={handleFileChange} />
        <button onClick={handleUpload}>Upload</button>
      </div>
    </div>
  );
};

export default FileUploadPopup;
