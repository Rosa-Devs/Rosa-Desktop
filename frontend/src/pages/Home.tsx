import RightSidebar from "./RightBar";
import { StartManager, AddManifets, ManifestList } from "../../wailsjs/go/src/DbManager"
import { useState } from "react";
import Chat from "../itemkit/Chat";
import TitleBar from "../itemkit/TitleBar";
import { manifest } from "../../wailsjs/go/models";



const Home = () => {

  const initialManifest = manifest.Manifest.createFrom({
    name: "Initial Name",
    pubsub: "Initial Pubsub",
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