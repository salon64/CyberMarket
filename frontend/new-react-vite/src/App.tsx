import { HashRouter, Routes, Route } from "react-router-dom";
import Layout from "./components/Layout";
import Login from "./components/Login/login";
import Marketplace from "./components/Marketplace/Marketplace";
import Inventory from "./components/Inventory/Inventory";
import NoPage from "./components/NoPage";
import Profile from "./components/Profile/Profile";
import Register from "./components/Register/Register";
import ProtectedRoute from "./components/Protection/ProtectedRoute";
import AdminPage from "./components/Admin/AdminPage";
import Cart from "./components/cart/cart";
import "./App.css";

const App = () => {
  return (
    <HashRouter>  {/* Corrected Router */}
      <Routes>
        {/* Public routes */}
        <Route path="/" element={<Layout />}>
          <Route index element={<Login />} /> {/* Fixed placement */}
          <Route path="Register" element={<Register />} />
          
          {/* Protected routes */}
          <Route element={<ProtectedRoute children/>}>
            <Route path="Marketplace" element={<Marketplace />} />
            <Route path="Inventory" element={<Inventory />} />
            <Route path="Profile" element={<Profile />} />
            <Route path="AdminPage" element={<AdminPage />} />
            <Route path="Cart" element={<Cart />} />
          </Route>

          {/* Catch-all for unknown routes */}
          <Route path="*" element={<NoPage />} />
        </Route>
      </Routes>
    </HashRouter>
  );
};

export default App;
