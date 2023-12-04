// Chat.tsx

import React, { useState, useEffect } from 'react';

interface Message {
  text: string;
  sender: 'user' | 'other';
}

const Chat: React.FC = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState<string>('');

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNewMessage(e.target.value);
  };

  const handleSend = () => {
    if (newMessage.trim() !== '') {
      setMessages([...messages, { text: newMessage, sender: 'user' }]);
      setNewMessage('');
      // You can handle sending the message to the server or any other logic here
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault(); // Prevents the default behavior of the Enter key (e.g., adding a new line)
      handleSend();
    }
  };

  useEffect(() => {
    // Simulate initial messages
    setMessages([
      { text: 'Hello!', sender: 'other' },
      { text: 'Hi there!', sender: 'user' },
      { text: 'How are you?', sender: 'other' },
    ]);
  }, []);

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
              {message.text}
            </div>
          </div>
        ))}
      </div>
      <div className="flex items-center p-4">
        <input
          type="text"
          value={newMessage}
          onKeyDown={handleKeyDown}
          onChange={handleInputChange}
          placeholder="My msg is?..."
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
