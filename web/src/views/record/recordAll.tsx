import {Button,  Space, Table} from 'antd'
import {SearchForm} from '../../components/SearchFrom'
import {CallRecordListDTO} from '../../@types/dto_list'
import {useEffect, useState} from 'react'
import {CallRecordEntity} from '../../@types/entity'

import type {ColumnsType} from "antd/es/table";

import dayjs from "dayjs";
import CommonPagination from "../../components/Pagination";
import AppAxios from "../../utils/request";
import  {OpenSeqModel,ShowIPText} from "../../utils/tools";

function RecordAll() {
    const [calls, setCalls] = useState<CallRecordEntity[]>([])

    const [loading, setLoading] = useState(false)
    const [listTotal, setListTotal] = useState(0)
    const [listPage, setListPage] = useState(0)
    const [listPageSize, setListPageSize] = useState(10)


    const [searchDTO, setSearchDTO] = useState<CallRecordListDTO>({})

    const columns: ColumnsType<CallRecordEntity> = [
        {
            title: '主叫',
            dataIndex: 'from_user',
            key: 'from_user',
        },
        {
            title: '被叫',
            dataIndex: 'to_user',
            key: 'to_user',
        },

        {
            title: '来源',
            dataIndex: 'src_host',
            render: (_, record) => {
                return <div>{record.src_host}({ShowIPText(record.src_country_name,record.src_city_name)})</div>
            },
        },



        {
            title: '目标',
            dataIndex: 'dst_host',
            render: (_, record) => {
                return <div>{record.dst_host}({ShowIPText(record.dst_country_name,record.dst_city_name)})</div>
            },
        },
        {
            title: 'UA',
            dataIndex: 'user_agent',
            key: 'user_agent',
            ellipsis: true,
        },

        {
            title: '时间',
            key: 'create_time',
            width: 180,
            render: (_, record) => dayjs(record.create_time).format('YYYY-MM-DD HH:mm:ss'),
        },
        {
            title: 'CallID',
            dataIndex: 'sip_call_id',
            key: 'sip_call_id',
            ellipsis: true,
        },
        {
            title: 'Action',
            key: 'action',
            width: 80,
            render: (_, record) => (
                <Button type="link" onClick={()=>{OpenSeqModel(record.sip_call_id)}}>信令</Button>
            ),
        },
    ]

    function Search(ft: CallRecordListDTO) {
        setSearchDTO(ft)
        setLoading(true)
    }


    useEffect(() => {
        searchDTO.page_size=listPageSize
        searchDTO.page=listPage

        setLoading(true)
        AppAxios.get('/record/all', {params: searchDTO})
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

    }, [searchDTO,listPage,listPageSize])

    return (<div>
            <h3>SIP消息管理</h3>
            <div style={{
                lineHeight: 3, display: "flex", justifyContent: "space-between",
            }}>
                <Space>
                    <SearchForm search={Search}/>
                </Space>
            </div>

            <Table scroll={{x: 1300}}
                   style={{marginTop:"10px"}}
                   columns={columns}
                   dataSource={calls}
                   loading={loading}
                   pagination={false}
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

export default RecordAll
