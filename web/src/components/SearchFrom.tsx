import {Button, Form, Input, DatePicker, Space, Row, Col, Divider} from 'antd'
import {CallRecordListDTO} from '../@types/dto_list'
import React, { useState } from "react";
import {RangeValueType, ValueDate} from "../@types/base.t";
import dayjs, {Dayjs} from "dayjs";
import { SearchOutlined, DownOutlined, UpOutlined } from '@ant-design/icons';

const {RangePicker} = DatePicker;
const FormatDatetime = "YYYY-MM-DD HH:mm:ss"
const FormatTime = "HH:mm:ss"
const TimePresets: ValueDate<RangeValueType<Dayjs>>[] = [
    {label: '今天', value: [dayjs().startOf('day'), dayjs()]},
    {label: '昨天', value: [dayjs().subtract(1, 'day').startOf('day'), dayjs().startOf('day')]},
    {label: '本周', value: [dayjs().startOf('week'), dayjs().endOf('week')]},
    {
        label: '上周',
        value: [dayjs().subtract(1, 'week').startOf('week').startOf('day'), dayjs().subtract(1, 'week').endOf('week').endOf('day')]
    },
    {label: '本月', value: [dayjs().startOf('month'), dayjs().endOf('month')]},
]

const SelectConvOptions=[{ value: 'eq', label: '等于' },{ value: 'neq', label: '不等于' }]

interface Prop {
    search: (ft: CallRecordListDTO) => void
}

export function SearchForm(p: Prop) {
    const [form] = Form.useForm();
    const [expand, setExpand] = useState(false);

    const onFinish = (ft: CallRecordListDTO) => {
        if (ft.date_picker) {
            ft.begin_time = ft.date_picker[0].format('YYYY-MM-DD') + ' ' + ft.date_picker[0].format('HH:mm:ss')
            ft.end_time = ft.date_picker[1].format('YYYY-MM-DD') + ' ' + ft.date_picker[1].format('HH:mm:ss')
        }
        console.log(ft)
        p.search(ft)
    }

    const onReset = () => {
        form.resetFields();
        form.submit();
    };

    const toggleExpand = () => {
        setExpand(!expand);
    };

    return (
            <Form
                form={form}
                name="search_form"
                onFinish={onFinish}
                autoComplete="off"
                layout="horizontal"
            >
                <Row gutter={[16, 0]}>
                    <Col span={6}>
                        <Form.Item<CallRecordListDTO> name="sip_call_id" style={{ marginBottom: 8 }}>
                            <Input placeholder="CALL_ID" allowClear style={{ borderRadius: '4px' }}/>
                        </Form.Item>
                    </Col>
                    <Col span={6}>
                        <Form.Item<CallRecordListDTO> name="session_id" style={{ marginBottom: 8 }}>
                            <Input placeholder="SessionID" allowClear style={{ borderRadius: '4px' }}/>
                        </Form.Item>
                    </Col>
                    <Col span={6}>
                        <Form.Item<CallRecordListDTO> name="date_picker" style={{ marginBottom: 8 }}>
                            <RangePicker
                                showTime={{format: FormatTime}}
                                presets={TimePresets}
                                placeholder={["开始时间", "结束时间"]}
                                format={FormatDatetime}
                                style={{ width: '100%', borderRadius: '4px' }}
                            />
                        </Form.Item>
                    </Col>
                    <Col span={6} style={{ textAlign: 'right' }}>
                        <Space>
                            <Button 
                                type="primary" 
                                htmlType="submit" 
                                icon={<SearchOutlined />}
                                style={{ borderRadius: '4px' }}
                            >
                                搜索
                            </Button>
                            <Button 
                                onClick={onReset}
                                style={{ borderRadius: '4px' }}
                            >
                                重置
                            </Button>
                            <Button 
                                type="link" 
                                onClick={toggleExpand}
                                icon={expand ? <UpOutlined /> : <DownOutlined />}
                            >
                                {expand ? '收起' : '展开'}
                            </Button>
                        </Space>
                    </Col>
                </Row>
                
                {expand && (
                    <>
                        <Divider style={{ margin: '12px 0' }} />
                        <Row gutter={[16, 0]}>
                            <Col span={6}>
                                <Form.Item<CallRecordListDTO> name="from_user" style={{ marginBottom: 8 }}>
                                    <Input placeholder="主叫号码" allowClear style={{ borderRadius: '4px' }}/>
                                </Form.Item>
                            </Col>
                            <Col span={6}>
                                <Form.Item<CallRecordListDTO> name="to_user" style={{ marginBottom: 8 }}>
                                    <Input placeholder="被叫号码" allowClear style={{ borderRadius: '4px' }}/>
                                </Form.Item>
                            </Col>
                            <Col span={6}>
                                <Form.Item<CallRecordListDTO> name="src_host" style={{ marginBottom: 8 }}>
                                    <Input placeholder="来源IP" allowClear style={{ borderRadius: '4px' }}/>
                                </Form.Item>
                            </Col>
                            <Col span={6}>
                                <Form.Item<CallRecordListDTO> name="dst_host" style={{ marginBottom: 8 }}>
                                    <Input placeholder="目标IP" allowClear style={{ borderRadius: '4px' }}/>
                                </Form.Item>
                            </Col>
                        </Row>
                        <Row gutter={[16, 0]}>
                            <Col span={6}>
                                <Form.Item<CallRecordListDTO> name="hangup_code" style={{ marginBottom: 8 }}>
                                    <Input placeholder="挂断代码" allowClear style={{ borderRadius: '4px' }}/>
                                </Form.Item>
                            </Col>
                        </Row>
                    </>
                )}
            </Form>

    )
}
