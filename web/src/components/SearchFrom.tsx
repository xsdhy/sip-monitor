import {Button, Form, Input, DatePicker, Space, Flex, Row, Select} from 'antd'
import {CallRecordListDTO} from '../@types/dto_list'
import React from "react";
import {RangeValueType, ValueDate} from "../@types/base.t";
import dayjs, {Dayjs} from "dayjs";

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

    const formStyle: React.CSSProperties = {
        maxWidth: 'none',
        paddingLeft: 10,
    };
    const formItemStyle: React.CSSProperties = {
        marginRight: 10,
    };

    const formSelectItemStyle: React.CSSProperties = {
        width: 100,
    };


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



    return (
        <Form
            form={form}
            name="basic"
            style={formStyle}
            labelAlign="left"
            onFinish={onFinish}
            autoComplete="off"
        >
            <Row gutter={24}>
                <Form.Item<CallRecordListDTO> style={formItemStyle} label="CALL_ID" name="sip_call_id"><Input allowClear/></Form.Item>
                {/* <Form.Item<CallRecordListDTO> style={formItemStyle} label="节点IP" name="node_ip"><Input allowClear/></Form.Item> */}

                {/* <Form.Item<CallRecordListDTO> style={formItemStyle} label="SIPMethod">
                    <Space.Compact>
                        <Form.Item name="sip_method_opr"><Select style={formSelectItemStyle} defaultValue="eq" options={SelectConvOptions}/></Form.Item>
                        <Form.Item name="sip_method"><Input style={{width:80}} allowClear/></Form.Item>
                    </Space.Compact>
                </Form.Item> */}


                {/* <Form.Item<CallRecordListDTO> style={formItemStyle} label="ResponseCode">
                    <Space.Compact>
                        <Form.Item name="response_code_opr"><Select style={formSelectItemStyle} defaultValue="eq" options={SelectConvOptions}/></Form.Item>
                        <Form.Item name="response_code"><Input style={{width:80}} allowClear/></Form.Item>
                    </Space.Compact>
                </Form.Item> */}


                {/* <Form.Item<CallRecordListDTO> style={formItemStyle} label="UA">
                    <Space.Compact>
                        <Form.Item name="ua_opr"><Select style={formSelectItemStyle} defaultValue="eq" options={SelectConvOptions}/></Form.Item>
                        <Form.Item name="ua"><Input allowClear/></Form.Item>
                    </Space.Compact>
                </Form.Item> */}


                <Form.Item<CallRecordListDTO> style={formItemStyle} label="主叫">
                    <Space.Compact>
                        <Form.Item name="from_user_opr"><Select style={formSelectItemStyle} defaultValue="eq" options={SelectConvOptions}/></Form.Item>
                        <Form.Item name="from_user"><Input allowClear/></Form.Item>
                    </Space.Compact>
                </Form.Item>
                <Form.Item<CallRecordListDTO> style={formItemStyle} label="来源IP">
                    <Space.Compact>
                        <Form.Item name="src_host_opr"><Select style={formSelectItemStyle} defaultValue="eq" options={SelectConvOptions}/></Form.Item>
                        <Form.Item name="src_host"><Input allowClear/></Form.Item>
                    </Space.Compact>
                </Form.Item>


                <Form.Item<CallRecordListDTO> style={formItemStyle} label="被叫">
                    <Space.Compact>
                        <Form.Item name="to_user_opr"><Select style={formSelectItemStyle} defaultValue="eq" options={SelectConvOptions}/></Form.Item>
                        <Form.Item name="to_user"><Input allowClear/></Form.Item>
                    </Space.Compact>
                </Form.Item>
                <Form.Item<CallRecordListDTO> style={formItemStyle} label="目标IP">
                    <Space.Compact>
                        <Form.Item name="dst_host_opr"><Select style={formSelectItemStyle} defaultValue="eq" options={SelectConvOptions}/></Form.Item>
                        <Form.Item name="dst_host"><Input allowClear/></Form.Item>
                    </Space.Compact>
                </Form.Item>

               

                <Form.Item<CallRecordListDTO> style={formItemStyle} label="时间" name="date_picker">
                    <RangePicker
                        showTime={{format: FormatTime}}
                        presets={TimePresets}
                        placeholder={["开始时间", "结束时间"]}
                        format={FormatDatetime}
                    />
                </Form.Item>
            </Row>
            <div style={{textAlign: 'right'}}>
                <Space size="small">
                    <Button type="primary" htmlType="submit">搜索</Button>
                    <Button type="default" htmlType="button" onClick={onReset}>重置</Button>
                </Space>
            </div>
        </Form>
    )
}
