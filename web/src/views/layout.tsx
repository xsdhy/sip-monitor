import {ConfigProvider, Layout, Menu, theme, Dropdown, Button, message, Avatar, Badge} from 'antd';
import React, {useState, useEffect} from 'react';
import {Route, Routes, useLocation, useNavigate} from "react-router-dom";
import { UserOutlined, LogoutOutlined, AppstoreOutlined, PhoneOutlined, DatabaseOutlined, TeamOutlined, HomeOutlined } from '@ant-design/icons';

//设置本地语言
import zhCN from 'antd/es/locale/zh_CN';
import AppAxios from '../utils/request';
import { customHistory } from '../utils/history';

import Home from "./home/home";

import RecordRegister from "./record/register";
import RecordCall from "./record/call";

import SystemStats from "./system/db";
import Profile from "./profile/profile";
import Users from "./system/users";

import dayjs from 'dayjs';
import dayjsLocal from 'dayjs/locale/zh-cn';
dayjs.locale(dayjsLocal)


const {Header, Content, Footer, Sider} = Layout;

const items = [
    {label: '首页', key: 'home', icon: <HomeOutlined /> },
    {label: '注册', key: 'record/register', icon: <AppstoreOutlined /> },
    {label: '呼叫', key: 'record/call', icon: <PhoneOutlined /> },
    {label: '数据', key: 'system/stats', icon: <DatabaseOutlined /> },
    {label: '用户', key: 'system/users', icon: <TeamOutlined /> },
];


const BackendLayout = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const [current, setCurrent] = useState(location.pathname.substring(1) || 'home');
    const [collapsed, setCollapsed] = useState(false);
    const {token: {colorBgContainer, colorPrimary, borderRadius}} = theme.useToken();
    const [userData, setUserData] = useState<any>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        // Check if user is authenticated
        const token = localStorage.getItem('token');
        if (!token) {
            customHistory.push('/login');
            return;
        }

        // Fetch user info
        AppAxios.get('/user/current')
            .then(res => {
                if (res.data.code === 2000) {
                    setUserData(res.data.data);
                } else {
                    message.error('获取用户信息失败，请重新登录');
                    localStorage.removeItem('token');
                    customHistory.push('/login');
                }
            })
            .catch(err => {
                console.error('Failed to get user info:', err);
                localStorage.removeItem('token');
                customHistory.push('/login');
            })
            .finally(() => {
                setLoading(false);
            });
    }, []);

    const onClick = (e:any) => {
        setCurrent(e.key);
        navigate(`/${e.key}`);
    };

    const handleLogout = () => {
        localStorage.removeItem('token');
        customHistory.push('/login');
        message.success('已退出登录');
    };

    const handleProfileClick = () => {
        navigate('/profile');
    };

    const userMenu = [
        {
            key: 'profile',
            label: '个人信息',
            icon: <UserOutlined />,
            onClick: handleProfileClick,
        },
        {
            key: 'logout',
            label: '退出登录',
            icon: <LogoutOutlined />,
            onClick: handleLogout,
        },
    ];

    if (loading) {
        return <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
            加载中...
        </div>;
    }

    return (
        <ConfigProvider
            theme={{
                algorithm: [ theme.compactAlgorithm],
                token: {
                    "wireframe": false,
                    "borderRadius": 4,
                    "colorPrimary": "#1677ff"
                },
            }}
            locale={zhCN}>
            <Layout style={{minHeight:"100vh"}}>
                <Sider 
                    width={100} 
                    collapsible 
                    collapsed={collapsed} 
                    onCollapse={(value) => setCollapsed(value)}
                    style={{
                        boxShadow: '0 2px 8px rgba(0, 0, 0, 0.15)',
                        zIndex: 10
                    }}
                >
                    <div
                        style={{
                            height: 64,
                            margin: 16,
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            fontSize: collapsed ? '14px' : '20px',
                            fontWeight: 'bold',
                            color: '#fff'
                        }}
                    >
                        {collapsed ? 'SIP' : 'SIP'}
                    </div>
                    <Menu
                        theme="dark"
                        mode="inline"
                        selectedKeys={[current]}
                        onClick={onClick}
                        defaultSelectedKeys={['home']}
                        items={items}
                        style={{
                            borderRight: 0
                        }}
                    />
                </Sider>
                <Layout className="site-layout">
                    <Header style={{
                        padding: '0 24px', 
                        background: colorBgContainer,
                        display: 'flex',
                        justifyContent: 'flex-end',
                        alignItems: 'center',
                        height: 64,
                        boxShadow: '0 1px 4px rgba(0, 0, 0, 0.1)',
                        position: 'sticky',
                        top: 0,
                        zIndex: 1,
                    }}>
                        {userData && (
                            <Dropdown menu={{ items: userMenu }} placement="bottomRight">
                                <Button type="link" style={{ height: 64, display: 'flex', alignItems: 'center' }}>
                                <Avatar 
                                            style={{ backgroundColor: colorPrimary, marginRight: 8 }} 
                                            icon={<UserOutlined />} 
                                        />
                                    <span style={{ marginLeft: 8 }}>{userData.nickname || userData.username}</span>
                                </Button>
                            </Dropdown>
                        )}
                    </Header>
                    <Content style={{ margin: '16px', overflow: 'initial' }}>
                        <div style={{
                            padding: 24, 
                            minHeight: "calc(100vh - 64px - 69px - 32px)", 
                            background: colorBgContainer,
                            borderRadius: borderRadius,
                            boxShadow: '0 1px 4px rgba(0, 0, 0, 0.1)'
                        }}>
                            <Routes>
                                <Route path="home" element={<Home/>}/>
                                <Route path="record/register" element={<RecordRegister/>}/>
                                <Route path="record/call" element={<RecordCall/>}/>
                                <Route path="system/stats" element={<SystemStats/>}/>
                                <Route path="system/users" element={<Users/>}/>
                                <Route path="profile" element={<Profile/>}/>
                                <Route path="/" element={<Home/>}/>
                            </Routes>
                        </div>
                    </Content>
                    <Footer style={{
                        textAlign: 'center',
                        padding: '16px 50px',
                        color: 'rgba(0, 0, 0, 0.45)',
                        fontSize: '14px'
                    }}>
                        SIP监控平台 © {new Date().getFullYear()} 
                    </Footer>
                </Layout>
            </Layout>
        </ConfigProvider>
    )
}

export default BackendLayout;