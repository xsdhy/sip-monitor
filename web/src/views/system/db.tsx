import {Button, DatePicker, Flex, Form, Modal, notification, Select, Table} from 'antd'
import {CallRecordListDTO, CleanSipRecordDTO} from '../../@types/dto_list'
import {useEffect, useState} from 'react'
import {SystemDBStatsVO} from '../../@types/entity'

import type {ColumnsType} from "antd/es/table";
import {ExclamationCircleFilled} from '@ant-design/icons';
import dayjs from "dayjs";
import AppAxios from "../../utils/request";
import {SizeConversion} from "../../utils/tools";
import {RangePickerProps} from "antd/es/date-picker";


function SystemStats() {
    const {confirm} = Modal;

    const [isModalOpen, setIsModalOpen] = useState(false);

    const [calls, setCalls] = useState<SystemDBStatsVO[]>([])

    const [loading, setLoading] = useState(true)


    // eslint-disable-next-line arrow-body-style
    const disabledDate: RangePickerProps['disabledDate'] = (current) => {
        return current && current > dayjs().endOf('day');
    };


    const columns: ColumnsType<SystemDBStatsVO> = [
        {
            title: '名称',
            dataIndex: 'name',
            key: 'name',
        },
        {
            title: '文档数',
            dataIndex: 'count',
            key: 'count',
        },
        {
            title: '平均对象大小',
            dataIndex: 'avg_obj_size',
            render: (_, record) => {
                return <span>{SizeConversion(record.avg_obj_size)}</span>
            },
        },

        {
            title: '内存中总大小',
            dataIndex: 'size',
            render: (_, record) => {
                return <span>{SizeConversion(record.size)}</span>
            },
        },
        {
            title: '存储大小',
            dataIndex: 'storage_size',
            render: (_, record) => {
                return <span>{SizeConversion(record.storage_size)}</span>
            },
        },
        {
            title: '索引数量',
            dataIndex: 'index_count',
            key: 'index_count',
        },
        {
            title: '索引总大小',
            dataIndex: 'total_index_size',
            render: (_, record) => {
                return <span>{SizeConversion(record.total_index_size)}</span>
            },
        },
        {
            title: '空闲大小',
            dataIndex: 'free_storage_size',
            render: (_, record) => {
                return <span>{SizeConversion(record.free_storage_size)}</span>
            },
        },
        {
            title: '总大小',
            dataIndex: 'total_size',
            render: (_, record) => {
                return <span>{SizeConversion(record.total_size)}</span>
            },
        },
        {
            title: '动作',
            key: 'action',
            width: 80,
            render: (_, record) => (
                <Button type="primary"
                        danger
                        disabled={record.name !== "call_records"}
                        onClick={() => {
                            setIsModalOpen(true)
                        }}>清理</Button>
            ),
        },
    ]


    useEffect(() => {
        if (loading) {
            AppAxios.get('/system/db/stats', {})
                .then(res => {
                    setCalls(res.data.data)
                    setLoading(false)
                })
                .catch()
        }
    }, [loading])

    return (<div>
            <h3>数据库状态</h3>

            <Table scroll={{x: 1300}}
                   style={{marginTop: "10px"}}
                   columns={columns}
                   dataSource={calls}
                   loading={loading}
                   pagination={false}
                   rowKey="id"/>

            <Modal open={isModalOpen}
                   title="数据库清理"
                   onCancel={() => {
                       setIsModalOpen(false)
                   }}
                   footer={null}>
                <Form
                    labelCol={{span: 8}}
                    wrapperCol={{span: 16}}
                    initialValues={{end_time: dayjs().subtract(1, 'day').startOf('day')}}
                    style={{maxWidth: 600}}
                    onFinish={(req: CleanSipRecordDTO) => {
                        req.begin_time = req.begin_time ? dayjs(req.begin_time).format("YYYY-MM-DD HH:mm:ss") : ""
                        req.end_time = req.end_time ? dayjs(req.end_time).format("YYYY-MM-DD HH:mm:ss") : ""

                        confirm({
                            title: '确认要进行清理操作吗？',
                            icon: <ExclamationCircleFilled/>,
                            content: '此操作不可逆',
                            onOk() {
                                AppAxios.get('/system/db/clean_sip_record', {params: req})
                                    .then(res => {
                                        if (res.data.code !== 200) {
                                            notification["error"]({
                                                message: "清理失败",
                                                description: res.data.msg,
                                            });
                                        } else {
                                            setIsModalOpen(false)
                                            setLoading(true)
                                            notification["success"]({
                                                message: "操作成功",
                                                description: res.data.msg,
                                            });
                                        }
                                    })
                                    .catch()
                            },
                            onCancel() {
                                console.log('Cancel');
                            },
                        });


                    }}
                    autoComplete="off"
                >
                    <Flex wrap="wrap" gap="small">
                        <Form.Item<CleanSipRecordDTO>
                            label="事件方法"
                            name="method"
                        >
                            <Select
                                allowClear
                                defaultValue="REGISTER"
                                style={{width: 300}}
                                options={[
                                    {value: '', label: '全部'},
                                    {value: 'INVITE', label: 'INVITE'},
                                    {value: 'REGISTER', label: 'REGISTER'},
                                    {value: 'OPTION', label: 'OPTION'},
                                ]}
                            />
                        </Form.Item>


                        <Form.Item<CallRecordListDTO>
                            label="开始时间"
                            name="begin_time">
                            <DatePicker
                                style={{width: 300}}
                                allowClear
                                disabledDate={disabledDate}
                                showTime={{format: 'HH:mm:ss'}}
                                format="YYYY-MM-DD HH:mm:ss"
                                presets={[
                                    {label: '1天前', value: dayjs().subtract(1, 'day').startOf('day')},
                                    {label: '2天前', value: dayjs().subtract(2, 'day').startOf('day')},
                                    {label: '3天前', value: dayjs().subtract(3, 'day').startOf('day')},
                                    {label: '1周前', value: dayjs().subtract(1, 'week').startOf('day')},
                                    {label: '2周前', value: dayjs().subtract(2, 'week').startOf('day')},
                                    {label: '1月前', value: dayjs().subtract(1, 'month').startOf('day')},
                                    {label: '2月前', value: dayjs().subtract(2, 'month').startOf('day')},
                                ]}
                            />
                        </Form.Item>
                        <Form.Item<CallRecordListDTO>
                            label="截止时间"
                            rules={[{required: true, message: '清选择截止时间'}]}
                            name="end_time">
                            <DatePicker
                                allowClear
                                style={{width: 300}}
                                showTime={{format: 'HH:mm:ss'}}
                                disabledDate={disabledDate}
                                format="YYYY-MM-DD HH:mm:ss"
                                presets={[
                                    {label: '1天前', value: dayjs().subtract(1, 'day').startOf('day')},
                                    {label: '2天前', value: dayjs().subtract(2, 'day').startOf('day')},
                                    {label: '3天前', value: dayjs().subtract(3, 'day').startOf('day')},
                                    {label: '1周前', value: dayjs().subtract(1, 'week').startOf('day')},
                                    {label: '2周前', value: dayjs().subtract(2, 'week').startOf('day')},
                                    {label: '1月前', value: dayjs().subtract(1, 'month').startOf('day')},
                                    {label: '2月前', value: dayjs().subtract(2, 'month').startOf('day')},
                                ]}
                            />
                        </Form.Item>

                        <Form.Item wrapperCol={{offset: 8, span: 16}}>
                            <Button style={{width: 300}} type="primary" danger htmlType="submit">
                                清理
                            </Button>
                        </Form.Item>
                    </Flex>
                </Form>
            </Modal>
        </div>
    )
}

export default SystemStats
