import RightSidebar from "./RightBar";

import { useState } from "react";
import Chat from "../itemkit/Chat";
import TitleBar from "../itemkit/TitleBar";
import { models } from "../models/manifest";




const Home = () => {

  const initialManifest = models.Manifest.createFrom({
    name: "Initial Name",
    pubsub: "Initial Pubsub",
    optional: "{\"Image\": \"http://localhost\"}"
  });

  const [chatManifest, setManifest] = useState(initialManifest);
  //console.log("Manifet: " + chatManifest.name)
  
  return (
    <div className="home-container">
      <div className="app flex">
        <div className="flex-1 justify-center">
          <Chat manifest={chatManifest}/> 
        </div>
        <RightSidebar setManifest={setManifest} />
      </div>
    </div>
  );
}


export default Home