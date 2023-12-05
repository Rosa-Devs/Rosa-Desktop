// Chat.tsx

import React, { useState, useEffect } from 'react';
import { GetMessages, NewMessage } from '../../wailsjs/go/src/DbManager';
import { src, manifest } from '../../wailsjs/go/models';




const Chat: React.FC<{ manifest: manifest.Manifest }> = ({ manifest }) => {
  const [messages, setMessages] = useState<src.Message[]>([]);
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

  //console.log("My manifet: " + manifest.name)

//   const fetchMsg = async () => {
//   try {
//     // Make your asynchronous API call or fetch data here
//     const data = await GetMessages(manifest);
//     console.log(data)   
//     return data;
//   } catch (error) {
//     console.error('Error fetching data:', error);
//     return null;
//   }
// };

const fetchMsg = async () => {
  try {
    // Make your asynchronous API call or fetch data here
    const data = await GetMessages(manifest);
    console.log(data);

    data.sort((a, b) => {
      const timei = new Date(a.time).getTime();
      const timej = new Date(b.time).getTime();

      // Compare timestamps to sort from latest to oldest
      return timej - timei;
    });
    
    return data;
  } catch (error) {
    console.error('Error fetching data:', error);
    return null;
  }
};



useEffect(() => {
  const fetchDataInterval = setInterval(async () => {
    const result = await fetchMsg();
    if (result !== null) {
      setMessages(result);
    }
  }, 1000); // Fetch data every 1 second

  // Cleanup the interval when the component unmounts
  return () => clearInterval(fetchDataInterval);
}, [manifest]); // Include 'manifest' in the dependency array

  
  const onEnterPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if(e.keyCode == 13 && e.shiftKey == false) {
      e.preventDefault();
      handleSend()
    }
  }

  
  

  return (
    <div className="flex flex-col h-screen">
      <div className="flex-1 overflow-y-auto px-4 py-8 flex flex-col-reverse">
        {messages.slice().reverse().map((message, index) => (
          <div
            key={index}
            className={`mb-4 ${
              message.sender === 'user' ? 'self-end' : 'self-start'
            } w-3/4`}
          >
            <div
              className={`p-4 max-w-xs mx-2 rounded-lg ${
                message.sender === 'user' ? 'bg-purple-600 text-white ml-auto' : 'bg-gray-200 text-black mr-auto'
              }`}
            >
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
