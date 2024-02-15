import { Route, RouterProvider, createBrowserRouter, createRoutesFromElements } from "react-router-dom";
import Home from "./pages/Home";
import ChatBar from "./pages/RightBar";
import { useEffect, useState } from "react";
import { Autorized, StartManager } from "../wailsjs/go/core/Core";
import Login from "./pages/Login";



function App() {


  const [isAuthorized, setIsAuthorized] = useState(false);
  useEffect(() => {
      setTimeout(() => {
        StartManager()
      }, 1000)

      setTimeout(async () => {
        setIsAuthorized(await Autorized())
        
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

export default App
