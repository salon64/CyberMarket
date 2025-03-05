import { BrowserRouter, Routes, Route } from "react-router";
import Layout from "./components/Layout";
import Login from "./components/Login/login";
import Marketplace from "./components/Marketplace/Marketplace";
import Inventory from "./components/Inventory/Inventory";
import NoPage from "./components/NoPage";
import Profile from "./components/Profile/Profile";
import Register from "./components/Register/Register";
import ProtectedRoute from "./components/Protection/ProtectedRoute"
import AdminPage from "./components/Admin/AdminPage";
import Cart from "./components/cart/cart";
import "./App.css";

const App = () => {
  return (
    <div>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route element={<ProtectedRoute children/>}>
              <Route path="Marketplace" element={<Marketplace />} />
              <Route path="Inventory" element={<Inventory />} />
              <Route path="Profile" element={<Profile/>} />
              <Route path="AdminPage" element={<AdminPage/>} />
              <Route path="Cart" element={<Cart/>} />
            </Route>
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
