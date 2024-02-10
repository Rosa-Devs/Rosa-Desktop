import { Route, RouterProvider, createBrowserRouter, createRoutesFromElements } from "react-router-dom";
import Home from "./pages/Home";
import ChatBar from "./pages/RightBar";
import { useEffect } from "react";
import { StartManager } from "../wailsjs/go/src/DbManager";
import CreateChat from "./pages/CreateChat";

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path="/" element={<Home />}>
      {/*<Route path="dashboard" element={<Home />} /> */}
      {/* ... etc. */}
      <Route path="/create" element={<CreateChat/>}></Route>
    </Route>
  )
);

function Router() {

    useEffect(() => {
      setTimeout(() => {
        StartManager("db1_test", "Mihalic2040")
      }, 1000)

      console.log("Db manager started");
    }, []);
    return (
    <div className="calm-gradient-background">
      <RouterProvider router={router} /> 
    </div>
    )
}

export default Router
