import React, { useEffect, useRef, useState } from "react";
import { manifest } from "../../wailsjs/go/models";
import { DeleteManifest } from "../../wailsjs/go/src/DbManager";

// ContextMenu component

// ContextMenu component
const ContextMenu: React.FC<{ onDelete: () => void; Close: () => void; xPos: number; yPos: number }> = ({ onDelete, Close, xPos, yPos }) => {
  const menuRef = useRef<HTMLDivElement>(null);
  const [menuPosition, setMenuPosition] = useState({ x: xPos, y: yPos });

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(event.target as Node)) {
        Close()
        
      }
    };

    document.addEventListener("click", handleClickOutside);

    return () => {
      document.removeEventListener("click", handleClickOutside);
    };
  }, [Close]);

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
    </div>
  );
};

// Buble component
const Buble: React.FC<{ contact: manifest.Manifest }> = ({ contact }) => {
  const [contextMenuPosition, setContextMenuPosition] = useState({ x: 0, y: 0 });

  const handleBtnClick = () => {
    console.log(contact.name);
  };

  const handleContextMenu = (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
    setContextMenuPosition({ x: event.clientX, y: event.clientY });
  };

  const handleClose = () => {
    setContextMenuPosition({ x: 0, y: 0 });
  }

  const handleDelete = () => {
    // Add logic for deleting the contact here
    console.log(`Deleting ${contact.name}`);
    DeleteManifest(contact)
    window.location.reload();
    // Close the context menu after handling delete
    setContextMenuPosition({ x: 0, y: 0 });
  };

  return (
    <div className="relative-container">
      <button
        onClick={handleBtnClick}
        onContextMenu={handleContextMenu}
        className="sidebar-btn"
      >
        <img
          src="https://avatarfiles.alphacoders.com/762/thumb-76262.png"
          alt={contact.name}
          className="img-icon"
          title={contact.name}
        />
      </button>

      {contextMenuPosition.x !== 0 && contextMenuPosition.y !== 0 && (
        <ContextMenu onDelete={handleDelete} Close={handleClose} xPos={contextMenuPosition.x} yPos={contextMenuPosition.y} />
      )}
    </div>
  );
};

export default Buble;
