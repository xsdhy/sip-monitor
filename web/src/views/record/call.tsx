import {Button, Space, Table, Tag} from 'antd'
import {SearchForm} from '../../components/SearchFrom'
import {CallRecordListDTO} from '../../@types/dto_list'
import {useEffect, useState} from 'react'
import { SIPRecordCall} from '../../@types/entity'

import type {ColumnsType} from "antd/es/table";

import dayjs from "dayjs";
import CommonPagination from "../../components/Pagination";
import AppAxios from "../../utils/request";
import {OpenSeqModel} from "../../utils/tools";

function RecordCall() {
    const [calls, setCalls] = useState<SIPRecordCall[]>([])

    const [loading, setLoading] = useState(false)
    const [listTotal, setListTotal] = useState(0)
    const [listPage, setListPage] = useState(1)
    const [listPageSize, setListPageSize] = useState(10)


    const [searchDTO, setSearchDTO] = useState<CallRecordListDTO>({})

    const columns: ColumnsType<SIPRecordCall> = [
        {
            title: '主叫',
            dataIndex: 'from_user',
            key: 'from_user',
            fixed: 'left',
        },
        {
            title: '被叫',
            dataIndex: 'to_user',
            key: 'to_user',
            fixed: 'left',
        },
        {
            title: '来源',
            dataIndex: 'src_host',
            width: 180,
            render: (_, record) => {
                return <div>{record.src_addr}<br/>{record.dst_addr}</div>
            },
        },
        {
            title: '创建结束',
            key: 'create_time',
            width: 240,
            render: (_, record) => {
                return <div>
                    {record.create_time ? "创建:"+dayjs(record.create_time).format('YYYY-MM-DD HH:mm:ss') : ""}
                    <br/>
                    {record.end_time ? "结束:"+dayjs(record.end_time).format('YYYY-MM-DD HH:mm:ss') : ""}
                    </div>
            },
        },
        {
            title: '振铃应答',
            key: 'ringing_time',
            width: 240,
            render: (_, record) => {
                return <div>
                    {record.ringing_time ? "振铃:"+dayjs(record.ringing_time).format('YYYY-MM-DD HH:mm:ss') : ""}
                    <br/>
                    {record.answer_time ? "应答:"+dayjs(record.answer_time).format('YYYY-MM-DD HH:mm:ss') : ""}
                    </div>
            },
        },
        {
            title: '通话时长',
            key: 'duration',
            width: 240,
            render: (_, record) => {
                return <div><Tag color="blue">总时长:{record.call_duration}s</Tag>  <Tag color="default">振铃:{record.ringing_duration}s</Tag>  <Tag color="green">通话:{record.talk_duration}s</Tag></div>
            },
        },
        {
            title: 'Code',
            dataIndex: 'hangup_code',
            key: 'hangup_code',
            render: (_, record) => {
                return <div>{record.hangup_code}({record.hangup_cause})</div>
            },
        },
        {
            title: 'Action',
            key: 'action',
            width: 80,
            fixed: 'right',
            render: (_, record) => (
                <Button type="link" onClick={() => {
                    OpenSeqModel(record.sip_call_id, record.session_id)
                }}>信令</Button>
            ),
        },
    ]

    function Search(ft: CallRecordListDTO) {
        setSearchDTO(ft)
        setLoading(true)
    }


    useEffect(() => {
        searchDTO.page_size = listPageSize
        searchDTO.page = listPage

        setLoading(true)
        AppAxios.get('/record/call', {params: searchDTO})
            .then(res => {
                // @ts-ignore
                setCalls(res.data.data)
                // @ts-ignore
                setListTotal(res.data.meta.total)
                // @ts-ignore
                setListPageSize(res.data.meta.page_size)
                setLoading(false)
            })
            .catch()

    }, [searchDTO, listPage, listPageSize])

    return (<div>
            <h3>呼叫管理</h3>
            <div style={{lineHeight: 3, display: "flex", justifyContent: "space-between",}}>
            <SearchForm search={Search}/>
            </div>

            <Table scroll={{x: 1300}}
                   style={{marginTop:"10px"}}
                   columns={columns}
                   dataSource={calls}
                   loading={loading}
                   pagination={false}
                   bordered
                   footer={() => (<CommonPagination
                       onChange={(page: number) => {
                           setListPage(page);
                           setLoading(true);
                       }}
                       onShowSizeChange={(current: number, size: number) => {
                           setListPage(1);
                           setListPageSize(size);
                           setLoading(true);
                       }}
                       current={listPage}
                       total={listTotal}
                       size={listPageSize}
                   />)}
                   rowKey="id"/>
        </div>
    )
}

export default RecordCall
