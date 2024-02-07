import {ConfigProvider, Layout, Menu, theme} from 'antd';
import React, {useState} from 'react';
import {Route, Routes, useLocation, useNavigate} from "react-router-dom";

//设置本地语言
import zhCN from 'antd/es/locale/zh_CN';


import Home from "./home/home";
import RecordAll from "./record/recordAll";
import RecordRegister from "./record/register";
import RecordCall from "./record/call";

import SystemStats from "./system/db";


import dayjs from 'dayjs';
import dayjsLocal from 'dayjs/locale/zh-cn';
dayjs.locale(dayjsLocal)


const {Header, Content, Footer, Sider} = Layout;

const items = [
    {label: '工作台', key: 'home'},
    {label: '注册', key: 'record/register'},
    {label: '呼叫', key: 'record/call'},
    {label: '消息', key: 'record/all'},
    {label: '数据库', key: 'system/stats'},
];


const BackendLayout = () => {
    const navigate = useNavigate();

    const location = useLocation();
    const [current, setCurrent] = useState(location.pathname.substring(1));

    const [collapsed, setCollapsed] = useState(false);
    const {token: {colorBgContainer}} = theme.useToken();



    const onClick = (e:any) => {
        setCurrent(e.key);
        navigate(`/${e.key}`);
    };

    return (
        <ConfigProvider
            theme={{
                token: {
                    "wireframe": false,
                    "borderRadius": 3
                },
            }}
            locale={zhCN}>
            <Layout style={{minHeight:"100vh"}}>
                <Sider width={120} collapsible collapsed={collapsed} onCollapse={(value) => setCollapsed(value)}>
                    <div
                        style={{
                            height: 32,
                            margin: 16,
                            background: 'rgba(255, 255, 255, 0)',
                        }}
                    />
                    <Menu
                        theme="dark"
                        mode="inline"
                        selectedKeys={[current]}
                        onClick={onClick}
                        defaultSelectedKeys={['4']}
                        items={items}
                    />
                </Sider>
                <Layout className="site-layout">
                    <Header style={{padding: 0, marginBottom: 10, height: 40, background: colorBgContainer}}/>
                    <Content style={{margin: '0 10px',minHeight:"100vh"}}>
                        <div style={{padding: 10, minHeight:"100vh", background: colorBgContainer}}>
                            <Routes>
                                <Route path="/" element={<Home/>}/>
                                <Route path="/record/all" element={<RecordAll/>}/>
                                <Route path="/record/register" element={<RecordRegister/>}/>
                                <Route path="/record/call" element={<RecordCall/>}/>
                                <Route path="/system/stats" element={<SystemStats/>}/>
                            </Routes>
                        </div>
                    </Content>
                    <Footer style={{textAlign: 'center'}}>
                        SIP监控平台-SipMonitor © 2024
                    </Footer>
                </Layout>
            </Layout>
        </ConfigProvider>
    )
}

export default BackendLayout;