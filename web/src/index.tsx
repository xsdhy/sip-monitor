import ReactDOM from 'react-dom/client'


import React from 'react';

import {Routes, Route} from "react-router-dom"
import BackendLayout from "./views/layout";
import Login from "./views/login/login";

import SequenceDiagram from "./views/details/SequenceDiagram";

import CustomRouter from "./components/CustomRouter";
import {customHistory} from './utils/history'


const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);
root.render(
    <CustomRouter history={customHistory}>
        <Routes>
            <Route path="login" element={<Login/>}/>

            <Route path="/call/details" element={<SequenceDiagram/>}/>
            <Route path="/*" element={<BackendLayout/>}/>
        </Routes>
    </CustomRouter>
);
