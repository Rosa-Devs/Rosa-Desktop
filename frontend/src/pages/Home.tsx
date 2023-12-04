import RightSidebar from "./RightBar";
import { StartManager, AddManifets, ManifestList } from "../../wailsjs/go/src/DbManager"
import { useState } from "react";
import Chat from "../itemkit/Chat";
import TitleBar from "../itemkit/TitleBar";



const Home = () => {
  
  return (
    <div className="home-container">
      <div className="app flex">
        {/* Your main content goes here */}
        <div className="flex-1 justify-center">
          <Chat/> 
        </div>


        {/* Include the RightSidebar component */}

        
        <RightSidebar />
      </div>
    </div>
  );
}


export default Home