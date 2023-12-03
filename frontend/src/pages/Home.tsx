import RightSidebar from "./ChatBar";
import ChatBar from "./ChatBar"
import { StartManager, AddManifets, ManifestList } from "../../wailsjs/go/src/DbManager"
import { useState } from "react";



const Home = () => {
  
  return (
    <div className="app flex">
      {/* Your main content goes here */}
      <div className="flex-1">
        <p>Content</p>
      </div>


      {/* Include the RightSidebar component */}
      <RightSidebar />
    </div>
  );
}


export default Home