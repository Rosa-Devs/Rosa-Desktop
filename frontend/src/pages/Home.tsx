import RightSidebar from "./ChatBar";
import ChatBar from "./ChatBar"



const Home = () => {
    return (
    <div className="app flex">
      {/* Your main content goes here */}
      <div className="flex-1">
        <p>Main Content</p>
      </div>

      {/* Include the RightSidebar component */}
      <RightSidebar />
    </div>
  );
}


export default Home