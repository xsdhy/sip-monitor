import React, { useState, useEffect } from 'react';
import { Card, Form, Input, Button, Tabs, message, Spin } from 'antd';
import { userApi } from '@/apis/api';
import { UserInfo } from '@/@types/entity';

const { TabPane } = Tabs;

const Profile: React.FC = () => {
  const [profileForm] = Form.useForm();
  const [passwordForm] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [_, setUserData] = useState<UserInfo>();
  const [loadingUserData, setLoadingUserData] = useState(true);

  // 获取用户信息
  useEffect(() => {
    userApi.getCurrentUser()
      .then(res => {
        setUserData(res.data);
          profileForm.setFieldsValue({
            nickname: res.data.nickname,
            username: res.data.username,
          });
      })
      .catch(err => {
        console.error('Failed to get user info:', err);
        message.error('获取用户信息失败');
      })
      .finally(() => {
        setLoadingUserData(false);
      });
  }, [profileForm]);

  // 更新个人信息
  const handleUpdateProfile = (values: any) => {
    setLoading(true);
    
      userApi.updateUser({
        id: 0,
        username: "",
        status: "",
        create_at: "",
        update_at: "",
        nickname: values.nickname,
      })
      .then(() => {
        message.success('个人信息更新成功');
          // 刷新用户数据
          return userApi.getCurrentUser();
      })
      .catch(err => {
        console.error('Failed to update profile:', err);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  // 更新密码
  const handleUpdatePassword = (values: any) => {
    setLoading(true);
    userApi.updatePassword({
      old_password: values.oldPassword,
      new_password: values.newPassword,
    })
      .then(() => {
        message.success('密码更新成功');
          passwordForm.resetFields();
      })
      .catch(err => {
        console.error('Failed to update password:', err);
      })
      .finally(() => {
        setLoading(false);
      });
  };

  if (loadingUserData) {
    return (
      <div style={{ display: 'flex', justifyContent: 'center', padding: '100px' }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <Card title="个人中心" bordered={false}>
      <Tabs defaultActiveKey="profile">
        <TabPane tab="个人信息" key="profile">
          <Form
            form={profileForm}
            layout="vertical"
            onFinish={handleUpdateProfile}
            style={{ maxWidth: 500 }}
          >
            <Form.Item
              label="用户名"
              name="username"
            >
              <Input disabled />
            </Form.Item>
            <Form.Item
              label="昵称"
              name="nickname"
              rules={[{ required: true, message: '请输入昵称!' }]}
            >
              <Input />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" loading={loading}>
                更新信息
              </Button>
            </Form.Item>
          </Form>
        </TabPane>
        <TabPane tab="修改密码" key="password">
          <Form
            form={passwordForm}
            layout="vertical"
            onFinish={handleUpdatePassword}
            style={{ maxWidth: 500 }}
          >
            <Form.Item
              label="原密码"
              name="oldPassword"
              rules={[{ required: true, message: '请输入原密码!' }]}
            >
              <Input.Password />
            </Form.Item>
            <Form.Item
              label="新密码"
              name="newPassword"
              rules={[{ required: true, message: '请输入新密码!' }]}
            >
              <Input.Password />
            </Form.Item>
            <Form.Item
              label="确认新密码"
              name="confirmPassword"
              dependencies={['newPassword']}
              rules={[
                { required: true, message: '请确认新密码!' },
                ({ getFieldValue }) => ({
                  validator(_, value) {
                    if (!value || getFieldValue('newPassword') === value) {
                      return Promise.resolve();
                    }
                    return Promise.reject(new Error('两次输入的密码不一致!'));
                  },
                }),
              ]}
            >
              <Input.Password />
            </Form.Item>
            <Form.Item>
              <Button type="primary" htmlType="submit" loading={loading}>
                更新密码
              </Button>
            </Form.Item>
          </Form>
        </TabPane>
      </Tabs>
    </Card>
  );
};

export default Profile; 