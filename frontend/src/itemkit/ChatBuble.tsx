import React, { useEffect, useRef, useState } from "react";
import Optional from "../models/Optional";
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import { ExportManifest } from "../../wailsjs/go/app/App";
import { DeleteManifest } from "../api/api";
import { models } from "../models/manifest";

// ContextMenu component

interface ContextMenuProps {
  onDelete: () => void;
  onClose: () => void;
  onExport: () => void;
  xPos: number;
  yPos: number;
}
// ContextMenu component
const ContextMenu: React.FC<ContextMenuProps> = ({ onDelete, onClose, onExport, xPos, yPos }) => {
  const menuRef = useRef<HTMLDivElement>(null);
  const [menuPosition, setMenuPosition] = useState({ x: xPos, y: yPos });

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        onClose()
        
      }
    };

    document.addEventListener("click", handleClickOutside);

    return () => {
      document.removeEventListener("click", handleClickOutside);
    };
  }, [onClose]);

  useEffect(() => {
    // Adjust the menu position if it goes beyond the right side of the screen
    const screenWidth = window.innerWidth;
    const menuWidth = menuRef.current?.offsetWidth || 0;

    if (xPos + menuWidth > screenWidth) {
      setMenuPosition({ x: xPos - menuWidth, y: yPos });
    }
  }, [xPos, yPos]);

  const menuStyle: React.CSSProperties = {
    position: "absolute",
    top: menuPosition.y,
    left: menuPosition.x,
  };

  return (
    <div className="context-menu" style={menuStyle} ref={menuRef}>
      <button onClick={onDelete}>Delete</button>
      <br />
      <button onClick={onExport}>Share</button>
    </div>
  );
};

// Buble component
const Buble: React.FC<{ contact: models.Manifest; setManifest: React.Dispatch<React.SetStateAction<any>> }> = ({ contact, setManifest }) => {
  const [contextMenuPosition, setContextMenuPosition] = useState({ x: 0, y: 0 });

  const handleBtnClick = () => {
    setManifest(contact)
  };

  const handleContextMenu = (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
    setContextMenuPosition({ x: event.clientX, y: event.clientY });
  };

  const handleClose = () => {
    setContextMenuPosition({ x: 0, y: 0 });
  }

  const handleExportManifest = () => {
    ExportManifest(contact)
  }

  const handleDelete = () => {
    // Add logic for deleting the contact here
    console.log(`Deleting ${contact.name}`);
    DeleteManifest(contact)
    window.location.reload();
    // Close the context menu after handling delete
    setContextMenuPosition({ x: 0, y: 0 });
  };

  let optional_obj: Optional = { Image: "https://www.shutterstock.com/image-vector/default-avatar-profile-icon-social-600nw-1677509740.jpg" };  // Assuming an appropriate default value
  if (contact.optional !== null && contact.optional !== undefined && contact.optional !== "") {
      // Parse the JSON string into an object
      optional_obj = JSON.parse(contact.optional);

      // Now you can use optional_obj as needed
      
  } else {
      // Handle the case when contact.optional is null or undefined
      console.error("contact.optional is null or undefined");
  }

  return (
    <div className="relative-container">
      <button
        onClick={handleBtnClick}
        onContextMenu={handleContextMenu}
        className="sidebar-btn"
      >
        <img
          src={optional_obj.Image}
          alt={contact.name}
          className="img-icon"
          title={contact.name}
        />
      </button>

      {contextMenuPosition.x !== 0 && contextMenuPosition.y !== 0 && (
        <ContextMenu onDelete={handleDelete} onExport={handleExportManifest} onClose={handleClose} xPos={contextMenuPosition.x} yPos={contextMenuPosition.y} />
      )}
    </div>
  );
};

export default Buble;

