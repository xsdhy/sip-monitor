import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Space, Popconfirm, message, Card, Typography, Tag } from 'antd';

import dayjs from 'dayjs';
import { UserAddOutlined, EditOutlined, DeleteOutlined } from '@ant-design/icons';
import { userApi } from '@/apis/api';
import { UserInfo } from '@/@types/entity';
const { Title } = Typography;



const Users: React.FC = () => {
  const [users, setUsers] = useState<UserInfo[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [visible, setVisible] = useState<boolean>(false);
  const [confirmLoading, setConfirmLoading] = useState<boolean>(false);
  const [currentUser, setCurrentUser] = useState<UserInfo | null>(null);
  const [form] = Form.useForm();

  // 获取用户列表
  const fetchUsers = () => {
    setLoading(true);
   userApi.getUserList()
      .then(res => {
        setUsers(res.data);
      })
      .catch(error => {
        console.error('获取用户列表失败:', error);
        message.error('获取用户列表失败');
      })
      .finally(() => {
        setLoading(false);
      });
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  // 添加/编辑用户
  const showModal = (user?: UserInfo) => {
    form.resetFields();
    if (user) {
      setCurrentUser(user);
      form.setFieldsValue({
        username: user.username,
        nickname: user.nickname,
        password: '',
        id: user.id,
      });
    } else {
      setCurrentUser(null);
      form.setFieldsValue({
        username: '',
        nickname: '',
        password: '',
      });
    }
    setVisible(true);
  };

  // 处理表单提交
  const handleSubmit = () => {
    form.validateFields()
      .then(values => {
        setConfirmLoading(true);
        if (currentUser) {
          values.id = currentUser.id
          // 更新用户
          userApi.updateUser(values)
            .then(() => {
              message.success('更新用户成功');
              setVisible(false);
              fetchUsers();
            })
            .catch(error => {
              console.error('更新用户失败:', error);
            })
            .finally(() => {
              setConfirmLoading(false);
            });
        } else {
          // 创建新用户
          userApi.createUser(values)
            .then(() => {
              message.success('创建用户成功');
              setVisible(false);
              fetchUsers();
            })
            .catch(error => {
              console.error('创建用户失败:', error);

            })
            .finally(() => {
              setConfirmLoading(false);
            });
        }
      })
      .catch(info => {
        console.log('表单验证失败:', info);
      });
  };

  // 删除用户
  const handleDelete = (id: number) => {
    userApi.deleteUser(id.toString())
      .then(() => {
        message.success('删除用户成功');
        fetchUsers();
      })
      .catch(error => {
        console.error('删除用户失败:', error);
        message.error('删除用户失败');
      });
  };

  // 表格列定义
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '用户名',
      dataIndex: 'username',
      key: 'username',
      render: (text: string) => <Tag color="blue">{text}</Tag>,
    },
    {
      title: '昵称',
      dataIndex: 'nickname',
      key: 'nickname',
    },
    {
      title: '创建时间',
      dataIndex: 'create_at',
      key: 'create_at',
      render: (text: string) => dayjs(text).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '更新时间',
      dataIndex: 'update_at',
      key: 'update_at',
      render: (text: string) => dayjs(text).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '操作',
      key: 'action',
      width: 160,
      render: (_: any, record: UserInfo) => (
        <Space size="middle">
          <Button 
            type="link" 
            icon={<EditOutlined />} 
            onClick={() => showModal(record)}
            style={{ padding: '0 8px' }}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除此用户吗？"
            description="删除后不可恢复！"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
            okButtonProps={{ danger: true }}
          >
            <Button 
              type="link" 
              danger 
              icon={<DeleteOutlined />}
              style={{ padding: '0 8px' }}
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Card>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
          <Title level={4} style={{ margin: 0 }}>用户管理</Title>
          <Button 
            type="primary" 
            icon={<UserAddOutlined />}
            onClick={() => showModal()}
          >
            添加用户
          </Button>
        </div>
        
        <Table
          columns={columns}
          dataSource={users}
          rowKey="id"
          loading={loading}
          pagination={{ 
            pageSize: 10,
            showSizeChanger: true,
            showTotal: (total) => `共 ${total} 条记录`
          }}
          bordered
        />
      </Card>
      
      <Modal
        title={currentUser ? '编辑用户' : '添加用户'}
        open={visible}
        onOk={handleSubmit}
        confirmLoading={confirmLoading}
        onCancel={() => setVisible(false)}
        maskClosable={false}
      >
        <Form
          form={form}
          layout="vertical"
          name="userForm"
        >
          <Form.Item
            name="username"
            label="用户名"
            rules={[
              { required: true, message: '请输入用户名' },
              { min: 3, message: '用户名至少3个字符' }
            ]}
          >
            <Input />
          </Form.Item>
          
          <Form.Item
            name="nickname"
            label="昵称"
            rules={[{ required: true, message: '请输入昵称' }]}
          >
            <Input />
          </Form.Item>
          
          <Form.Item
            name="password"
            label="密码"
            rules={[
              { 
                required: !currentUser, 
                message: '请输入密码' 
              },
              { 
                min: 6, 
                message: '密码至少6个字符',
                warningOnly: !!currentUser 
              }
            ]}
          >
            <Input.Password placeholder={currentUser ? '不填写则不修改密码' : '请输入密码'} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default Users; 