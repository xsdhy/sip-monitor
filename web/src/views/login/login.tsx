import { Button, Form, Input, message, Card, Typography, Spin, Row, Col } from 'antd';
import AppAxios from "../../utils/request";
import { customHistory } from "../../utils/history";
import { useState } from 'react';
import { UserOutlined, LockOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;

const Login = () => {
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);

    const onFinish = (values: any) => {
        setLoading(true);
        AppAxios.post("/login", values)
            .then(res => {
                if (2000 === res.data.code) {
                    window.localStorage.setItem("token", res.data.data);
                    message.success("登录成功");
                    customHistory.push("/");
                } else {
                    message.error(res.data.msg || "登录失败");
                }
            })
            .catch(err => {
                console.error("Login error:", err);
                message.error("登录失败，请稍后重试");
            })
            .finally(() => {
                setLoading(false);
            });
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
            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%' }}>

                    <Card
                        style={{
                            width: '100%',
                            borderRadius: '8px',
                            boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)'
                        }}
                    >
                        <div style={{ textAlign: "center", marginBottom: 24 }}>
                            <Title level={2} style={{ color: '#1677ff', marginBottom: 8 }}>SIP监控系统</Title>
                            <Text type="secondary">登录以访问系统</Text>
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
                                    rules={[{ required: true, message: "请输入用户名" }]}
                                >
                                    <Input 
                                        prefix={<UserOutlined style={{ color: 'rgba(0,0,0,.25)' }} />} 
                                        placeholder="用户名" 
                                    />
                                </Form.Item>
                                <Form.Item
                                    name="password"
                                    rules={[{ required: true, message: "请输入密码" }]}
                                >
                                    <Input.Password 
                                        prefix={<LockOutlined style={{ color: 'rgba(0,0,0,.25)' }} />} 
                                        placeholder="密码" 
                                    />
                                </Form.Item>
                                <Form.Item>
                                    <Button 
                                        type="primary" 
                                        htmlType="submit" 
                                        loading={loading} 
                                        style={{ width: '100%', height: '40px' }}
                                    >
                                        登录
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