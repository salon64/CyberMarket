import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router";
import Layout from "./pages/Layout";
import Login from "./pages/login";
import Marketplace from "./pages/Marketplace";
import Contact from "./pages/Contact";
import NoPage from "./pages/NoPage";
//import reactLogo from './assets/react.svg'
//import viteLogo from '/vite.svg'
import './App.css'
const App = () => {
  return (
    
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Layout />}>
                <Route index element={<Login />} />
                <Route path="Marketplace" element={<Marketplace />} />
                <Route path="contact" element={<Contact />} />
                <Route path="*" element={<NoPage />} />
                </Route>
            </Routes>
        </BrowserRouter>
  );
};

export default App;
