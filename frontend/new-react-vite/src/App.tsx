import { BrowserRouter, Routes, Route } from "react-router";
import Layout from "./components/Layout";
import Login from "./components/Login/login";
import Marketplace from "./components/Marketplace/Marketplace";
import Inventory from "./components/Inventory/Inventory";
import NoPage from "./components/NoPage";
import Profile from "./components/Profile/Profile";
import Register from "./components/Register/Register";
import ProtectedRoute from "./components/Protection/ProtectedRoute"
//import reactLogo from './assets/react.svg'
//import viteLogo from '/vite.svg'
import "./App.css";
import { useState } from "react";
type AuthUser = {id: string; name: string};

const App = () => {
  const [user, setUser] = useState<AuthUser | null>(null);
  return (
    <div>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Layout />}>
               {/* ignore error :) */}
             <Route path="Marketplace" element={<Marketplace />} />
            
            <Route path="Inventory" element={<Inventory />} />
            <Route path="Profile" element={<Profile/>} />
            
            <Route path="*" element={<NoPage />} />
          

            </Route>
          <Route index element={<Login />} />
          <Route path="Register" element={<Register />} />

          
        </Routes>
      </BrowserRouter>
    </div>
  );
};

export default App;
