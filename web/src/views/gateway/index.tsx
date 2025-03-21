import { useEffect, useState } from 'react'
import { Table, Button, Modal, Form, Input, message, Popconfirm, Space } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { Gateway } from '@/apis/gateway'
import { gatewayApi } from '@/apis/api'
import dayjs from 'dayjs'
const GatewayList = () => {
  const [data, setData] = useState<Gateway[]>([])
  const [loading, setLoading] = useState(false)
  const [modalVisible, setModalVisible] = useState(false)
  const [form] = Form.useForm()
  const [editingId, setEditingId] = useState<number | null>(null)

  const fetchData = async () => {
    setLoading(true)
    try {
      const res = await gatewayApi.getGatewayList()
      setData(res.data)
    } catch (error) {
      message.error('获取网关列表失败')
    }
    setLoading(false)
  }

  useEffect(() => {
    fetchData()
  }, [])

  const handleAdd = () => {
    setEditingId(null)
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (record: Gateway) => {
    setEditingId(record.id)
    form.setFieldsValue(record)
    setModalVisible(true)
  }

  const handleDelete = async (id: number) => {
    try {
      await gatewayApi.deleteGateway(id.toString())
      message.success('删除成功')
      fetchData()
    } catch (error) {
      message.error('删除失败')
    }
  }

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields()
      if (editingId) {
        await gatewayApi.updateGateway(editingId.toString(), values)
        message.success('更新成功')
      } else {
        await gatewayApi.createGateway(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchData()
    } catch (error) {
      message.error(editingId ? '更新失败' : '创建失败')
    }
  }

  const columns: ColumnsType<Gateway> = [
    {
      title: 'ID',
      dataIndex: 'id',
      width: 80,
    },
    {
      title: '名称',
      dataIndex: 'name',
    },
    {
      title: '地址',
      dataIndex: 'addr',
    },
    {
      title: '备注',
      dataIndex: 'remark',
    },
    {
      title: '创建时间',
      dataIndex: 'create_at',
      width: 180,
      render: (_, record) => {
        return dayjs(record.create_at).format('YYYY-MM-DD HH:mm:ss')
      },
    },
    {
      title: '操作',
      key: 'action',
      width: 200,
      render: (_, record) => (
        <Space>
          <Button type="link" onClick={() => handleEdit(record)}>
            编辑
          </Button>
          <Popconfirm
            title="确定要删除吗？"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div className="p-6">
      <div className="mb-4">
        <Button type="primary" onClick={handleAdd}>
          添加网关
        </Button>
      </div>
      <Table
        columns={columns}
        dataSource={data}
        rowKey="id"
        loading={loading}
        pagination={false}
      />
      <Modal
        title={editingId ? '编辑网关' : '添加网关'}
        open={modalVisible}
        onOk={handleModalOk}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          <Form.Item
            name="name"
            label="名称"
            rules={[{ required: true, message: '请输入名称' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item
            name="addr"
            label="地址"
            rules={[{ required: true, message: '请输入地址' }]}
          >
            <Input />
          </Form.Item>
          <Form.Item name="remark" label="备注">
            <Input.TextArea />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default GatewayList 