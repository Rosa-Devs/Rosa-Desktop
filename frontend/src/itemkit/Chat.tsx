// Chat.tsx

import React, { useState, useEffect } from 'react';
import { ChangeListeningDb, GetMessages, GetProfile, NewMessage } from '../../wailsjs/go/core/Core';
import { models, manifest } from '../../wailsjs/go/models';
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'




const Chat: React.FC<{ manifest: manifest.Manifest }> = ({ manifest }) => {
  const [messages, setMessages] = useState<models.Message[]>([]);
  const [newMessage, setNewMessage] = useState<string>('');

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNewMessage(e.target.value);
  };

  const handleSend = () => {
    if (newMessage.trim() !== '') {
      //setMessages([...messages, { text: newMessage, sender: 'user' }]);
      setNewMessage('');

      NewMessage(manifest, newMessage)

      
      // You can handle sending the message to the server or any other logic here
    }
  };

  

  const fetchMsg = async () => {
    try {
      // Make your asynchronous API call or fetch data here
      const data = await GetMessages(manifest);
      //console.log(data);

      //console.log(data[0].time)
      //console.log(new Date(data[0].time).getTime())
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

  useEffect(() => {
    EventsOff("update")
    EventsOn("update", () => {
      // Fetch data on update event
      fetchDataInterval()
      console.log("Db update event")
    })
  })
  

  return (
    <div className="flex flex-col h-screen">
      <div className="flex-1 overflow-y-auto px-4 py-8 flex flex-col-reverse">
        {messages.slice().reverse().map((message, index) => (
          <div
            key={index}
            className={`mb-4 ${
              message.sender === nickname ? 'self-end' : 'self-start'
            } w-3/4`}
          >
            <div className="">
              {
                message.sender === nickname ? (
                  <div className="avatar_msg">
                    <img
                      src={"https://www.shutterstock.com/image-vector/default-avatar-profile-icon-social-600nw-1677509740.jpg"}
                      alt={`${message.sender}'s Avatar`}
                      className="w-12 h-12 rounded-full self-end"
                    />
                  </div>
                ) :(
                  <div className="avatar_user">
                    <img
                      src={"https://www.shutterstock.com/image-vector/default-avatar-profile-icon-social-600nw-1677509740.jpg"}
                      alt={`${message.sender}'s Avatar`}
                      className="w-12 h-12 rounded-full self-start"
                    />
                  </div>
                  
                )
              }
            </div>
            
            <div
              className={`p-4 max-w-xs mx-2 rounded-lg ${
                message.sender === nickname ? 'bg-gray-200 text-black ml-auto' : 'bg-gray-200 text-black mr-auto'
              }`}
            >
              <div className="text-gray-500 text-xs mb-1">
                {message.sender}
              </div>
                {message.data}
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

export default Chat;
