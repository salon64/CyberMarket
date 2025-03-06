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
        {/* Public Routes */}
          <Route path="/" element={<Layout />}>
          {/* Protected Routes */}
            <Route element={<ProtectedRoute children roles="0"/>}>
              <Route path="Marketplace" element={<Marketplace />} />
              <Route path="Inventory" element={<Inventory />} />
              <Route path="Profile" element={<Profile/>} />

              <Route path="Cart" element={<Cart/>} />
            </Route>
              <Route element={<ProtectedRoute children roles="1"/>}>
                <Route path="AdminPage" element={<AdminPage/>} />
              </Route>
            <Route path="*" element={<NoPage />} />
          

            </Route>
          <Route index element={<Login />} />
          <Route path="Register" element={<Register />} />

          
        </Routes>
    </HashRouter>
  );
};

export default App;
