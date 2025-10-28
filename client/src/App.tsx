import { Routes, Route, BrowserRouter as Router } from "react-router-dom";
import { Toaster } from "react-hot-toast";
import Body from "./pages/Body";
import Login from "./pages/Login";
import Register from "./pages/Register";

export default function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Body />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
      </Routes>
      <Toaster />
    </Router>      
  );
}
