// Chat.tsx

import React, { useState, useEffect, useRef } from 'react';
import { EventsOn, EventsOff, LogPrint } from '../../wailsjs/runtime/runtime'
import Valid from "../asset/valid.png"
import Err from "../asset/error.svg"
import { toast } from 'react-toastify';
import { models } from "../models/manifest"
import { ChangeListeningDb, MessagesList, NewMessage, Trust } from '../api/api';
import { GetAddres, GetProfile } from '../../wailsjs/go/app/App';

const Chat: React.FC<{ manifest: models.Manifest }> = ({ manifest }) => {
  const [messages, setMessages] = useState<models.Message[]>([]);
  const [newMessage, setNewMessage] = useState<string>('');
  const [contextMenuMessage, setContextMenuMessage] = useState<models.Message | null>(null);


  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNewMessage(e.target.value);
  };

  const handleSend = () => {
    if (newMessage.trim() !== '') {
      //setMessages([...messages, { text: newMessage, sender: 'user' }]);
      setNewMessage('');

      // NewMessage(manifest, newMessage)
      NewMessage({msg: newMessage, manifest})
      
      // You can handle sending the message to the server or any other logic here
    }
  };

  

  const fetchMsg = async () => {
    try {
      // Make your asynchronous API call or fetch data here
      const data = await MessagesList(manifest);
      //console.log(data);

      //console.log(data[0].time)
      //console.log(new Date(data[0].time).getTime())
      if (data === null) {
        return null
      }
      data.sort((a, b) => new Date(a.time).getTime() - new Date(b.time).getTime());
      //console.log(data)
      return data;
    } catch (error) {
      console.error('Error fetching data:', error);
      return null;
    }
  };

  const fetchDataInterval = async () => {
      const result = await fetchMsg();


      if (result !== null) {
        
        setMessages(result);
      }
  } // Fetch msg data

  useEffect(() => {
    // Fetch data on first run
    fetchDataInterval()
  }, [manifest]);

  
  const onEnterPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if(e.keyCode == 13 && e.shiftKey == false) {
      e.preventDefault();
      handleSend()
    }
  }


  const [nickname, setNickname] = useState('User');
  useEffect(() => {
    const getProfile = async () => {
      try {
        // Call the GetProfileFunc function to fetch the profile data
        const data = await GetProfile();
        
        // Set the profile data in the component state
        setNickname(data);
        console.log(nickname, data)

        if (data === "") {
          setTimeout(() => {
            getProfile();
          }, 2000); // 5000 milliseconds = 5 seconds
        }
      } catch (error) {
        console.error("Error getting username:", error);

        // Retry after 5 seconds if data is an empty string
        
      }
  };

    // Call the getProfile function when the component mounts
    getProfile();
  }, []);

  useEffect(() => {
    ChangeListeningDb(manifest) 
  }, [manifest])



  const [eventSource, setEventSource] = useState<EventSource | null>(null);

  useEffect(() => {
    const setupEventSource = async () => {
      console.log('MOUNTED');

      // Close the existing EventSource connection if it exists
      if (eventSource !== null) {
        eventSource.close();
      }

      // Get the address dynamically
      const address = await GetAddres()
      // const address = "localhost:8080";
      const newEventSource = new EventSource(`http://${address}/event/listen`);

      // Set the event listeners for the new EventSource object
      newEventSource.addEventListener("message", (event) => {
        console.log("Received event:", event.data);
        // Add your event handling logic here
        fetchDataInterval()
      });

      newEventSource.onerror = function(error) {
        console.error('EventSource failed:', error);
      };

      // Set the new EventSource object to state
      setEventSource(newEventSource);
    };

    // Call setupEventSource when the component mounts or when the manifest changes
    setupEventSource();

    // Clean up the EventSource connection when the component unmounts or when the manifest changes
    return () => {
      if (eventSource !== null) {
        eventSource.close();
      }
    };
  }, [manifest]); // Dependency array includes 'manifest'

  


  // useEffect(() => {
  //   EventsOff("update")
  //   EventsOn("update", () => {
  //     // Fetch data on update event
  //     fetchDataInterval()
  //     console.log("Db update event")
  //   })
  // })
  
  

  const [contextMenuPosition, setContextMenuPosition] = useState({ x: 0, y: 0 });
  const [showContextMenu, setShowContextMenu] = useState(false);

  const handleContextMenu = (event: React.MouseEvent<HTMLDivElement>, msg: models.Message) => {
    event.preventDefault();
    setContextMenuPosition({ x: event.clientX, y: event.clientY });
    setShowContextMenu(true)
    setContextMenuMessage(msg)
  };

  const handleContextMenuClose = () => {
    setShowContextMenu(false);
  };

  return (
    <div className="flex flex-col h-screen">
      <div className="flex-1 overflow-y-auto px-4 py-8 flex flex-col-reverse">
        {messages.slice().reverse().map((message, index) => (
          <div
            onContextMenu={(e) => handleContextMenu(e, message)}
            key={index}
            className={`mb-4 ${
              message.sender.name === nickname ? 'self-end' : 'self-start'
            } w-3/4`}
          >
            <div className="">
              {
                message.sender.name === nickname ? (
                  <div className="avatar_msg">
                    <img
                      src={message.sender.avatar !== "" ? message.sender.avatar : "https://www.shutterstock.com/image-vector/default-avatar-profile-icon-social-600nw-1677509740.jpg"}
                      alt={`${message.sender}'s Avatar`}
                      className="w-12 h-12 rounded-full self-end"
                    />
                  </div>
                ) :(
                  <div className="avatar_user">
                    <img
                      src={message.sender.avatar !== "" ? message.sender.avatar : "https://www.shutterstock.com/image-vector/default-avatar-profile-icon-social-600nw-1677509740.jpg"}
                      alt={`${message.sender}'s Avatar`}
                      className="w-12 h-12 rounded-full self-start"
                    />
                  </div>
                  
                )
              }
            </div>
            
            <div
              className={`p-4 max-w-xs mx-2 rounded-lg ${
                message.sender.name === nickname ? 'bg-gray-200 text-black ml-auto' : 'bg-gray-200 text-black mr-auto'
              }`}
            >
              <div className="text-gray-500 text-xs mb-1">
                <div className="avatar-user">
                  {message.sender.name}
                  <img src={message.valid ? Valid : Err} className='icon' alt={message.valid ? "Valid" : "Error"}/>
                </div>
              </div>
                {message.data}
                {showContextMenu && (
                  <ContextMenu
                    onClose={handleContextMenuClose}
                    xPos={contextMenuPosition.x}
                    yPos={contextMenuPosition.y}
                    msg={contextMenuMessage}
                  />
                )}
            </div>
          </div> 
        ))}
      </div>
      <div className="flex items-center p-4">
        <input
          type="text"
          value={newMessage}
          onKeyDown={onEnterPress}
          // onKeyDown={handleKeyDown}
          onChange={handleInputChange}
          placeholder="SELECT FROM * ..."
          className="flex-1 p-2 border rounded-md mr-2"
        />
        <button onClick={handleSend} className="px-4 py-2 bg-purple-700 text-white rounded-md">
          Send
        </button>
      </div>
    </div>
  );
};


interface ContextMenuProps {
  onClose: () => void;
  xPos: number;
  yPos: number;
  msg: models.Message | null;
}

const ContextMenu: React.FC<ContextMenuProps> = ({onClose, xPos, yPos, msg}) => {
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

  const handleTrust = async () => {
    // toast("Wait")
    if (msg !== null) {
      Trust(msg)
      console.log(msg)
    }
    
    
    toast("Account trusted!")
  }

  return (
    <div className="context-menu" style={menuStyle} ref={menuRef}>
      <button onClick={handleTrust}>Trust this Aunthor</button>
      <br />
    </div>
  );
};

export default Chat;
