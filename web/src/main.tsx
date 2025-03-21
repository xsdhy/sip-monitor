import React from 'react';
import ReactDOM from 'react-dom/client';
import {Routes, Route} from "react-router-dom";
import BackendLayout from "./views/layout";
import Login from "./views/login/login";
import SequenceDiagram from "./views/details/SequenceDiagram";
import CustomRouter from "./components/CustomRouter";
import {customHistory} from './utils/history';

import './index.css';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <CustomRouter history={customHistory}>
      <Routes>
        <Route path="login" element={<Login/>}/>
        <Route path="/call/details" element={<SequenceDiagram/>}/>
        <Route path="/*" element={<BackendLayout/>}/>
      </Routes>
    </CustomRouter>
  </React.StrictMode>
); 