import { Route, RouterProvider, createBrowserRouter, createRoutesFromElements } from "react-router-dom";
import Home from "./pages/Home";
import ChatBar from "./pages/ChatBar";
import { useEffect } from "react";
import { StartManager } from "../wailsjs/go/src/DbManager";

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path="/" element={<Home />}>
      {/*<Route path="dashboard" element={<Home />} /> */}
      {/* ... etc. */}
    </Route>
  )
);

function Router() {

    useEffect(() => {
      setTimeout(() => {
        StartManager("db1_test")
      }, 1000)

      console.log("Db manager started");
    }, []);
    return (
       <RouterProvider router={router} /> 
    )
}

export default Router
