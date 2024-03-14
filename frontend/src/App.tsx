import { Route, RouterProvider, createBrowserRouter, createRoutesFromElements } from "react-router-dom";
import Home from "./pages/Home";
import ChatBar from "./pages/RightBar";
import { useEffect, useState } from "react";
import Login from "./pages/Login";
import {StartManager} from "../wailsjs/go/app/App"
import { Authorized } from "./api/api"



function App() {


  const [isAuthorized, setIsAuthorized] = useState(false);
  useEffect(() => {
      setTimeout(() => {
        StartManager()
      }, 1000)

      setTimeout(async () => {
        setIsAuthorized(await Authorized())
        
      }, 3000);

      console.log("Db manager started");
    }, []);


  return (
    <div className="calm-gradient-background">
      {isAuthorized ? (
        <Home/>
      ) : (
        <div className="">
          <div className="appAside" />
           <div className="appForm">
            <Login/>
           </div>
          
        </div>
        
      )}
    </div>
    )
}

export default App;

