import { Route, RouterProvider, createBrowserRouter, createRoutesFromElements } from "react-router-dom";
import Home from "./pages/Home";
import ChatBar from "./pages/ChatBar";

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path="/" element={<Home />}>
      {/*<Route path="dashboard" element={<Home />} /> */}
      {/* ... etc. */}
    </Route>
  )
);

function Router() {
    return (
       <RouterProvider router={router} /> 
    )
}

export default Router
