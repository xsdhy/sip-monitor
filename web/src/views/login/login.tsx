import { Button, Form, Input, message, Card, Typography, Spin, } from 'antd';
import { customHistory } from "../../utils/history";
import { useState } from 'react';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import api from '@/apis/api';

const { Title, Text } = Typography;

const Login = () => {
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);

    const onFinish = (values: any) => {
        setLoading(true);

        api.user.login(values).then(res => {
            window.localStorage.setItem("token", res.data);
            message.success("Login successful");
            customHistory.push("/");
        }).catch(err => {
            console.error("Login error:", err);
            message.error("Login failed, please try again later");
        })
    };

    return (
        <div style={{
            height: '100vh',
            width: '100vw',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            background: 'linear-gradient(to right, #1677ff, #06b6d4)',
            overflow: 'hidden',
            position: 'relative'
        }}>
            <div 
                style={{
                    position: 'absolute',
                    background: 'rgba(255, 255, 255, 0.1)',
                    borderRadius: '50%',
                    width: '500px',
                    height: '500px',
                    top: '-100px',
                    left: '-100px'
                }}
            />
            <div 
                style={{
                    position: 'absolute',
                    background: 'rgba(255, 255, 255, 0.1)',
                    borderRadius: '50%',
                    width: '300px',
                    height: '300px',
                    bottom: '50px',
                    right: '100px'
                }}
            />
            <div style={{ display: 'flex',width: '80%',maxWidth: '400px', justifyContent: 'center', alignItems: 'center', height: '100%' }}>
                    <Card
                        style={{
                            width: '100%',
                            borderRadius: '8px',
                            boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)'
                        }}
                    >
                        <div style={{ textAlign: "center", marginBottom: 24 }}>
                            <Title level={2} style={{ color: '#1677ff', marginBottom: 8 }}>SIP Monitor</Title>
                            <Text type="secondary">Login to access the system</Text>
                        </div>

                        <Spin spinning={loading}>
                            <Form
                                form={form}
                                onFinish={onFinish}
                                size="large"
                                layout="vertical"
                            >
                                <Form.Item
                                    name="username"
                                    rules={[{ required: true, message: "Please enter your username" }]}
                                >
                                    <Input 
                                        prefix={<UserOutlined style={{ color: 'rgba(0,0,0,.25)' }} />} 
                                        placeholder="Username" 
                                    />
                                </Form.Item>
                                <Form.Item
                                    name="password"
                                    rules={[{ required: true, message: "Please enter your password" }]}
                                >
                                    <Input.Password 
                                        prefix={<LockOutlined style={{ color: 'rgba(0,0,0,.25)' }} />} 
                                        placeholder="Password" 
                                    />
                                </Form.Item>
                                <Form.Item>
                                    <Button 
                                        type="primary" 
                                        htmlType="submit" 
                                        loading={loading} 
                                        style={{ width: '100%', height: '40px' }}
                                    >
                                        Login
                                    </Button>
                                </Form.Item>
                            </Form>
                        </Spin>
                    </Card>

            </div>
        </div>
    );
}

export default Login;