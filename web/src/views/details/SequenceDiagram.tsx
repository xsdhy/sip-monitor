import {useEffect, useRef, useState} from 'react'
import {Tag, Spin, Modal, Empty} from 'antd'
import './Sequence.css'


import * as ssd from 'svg-sequence-diagram'
import {createSeqHtml, getProtocolName} from './util'
import {CallRecordDetailsVO, CallRecordEntity} from '../../@types/entity'
import AppAxios from "../../utils/request";
import dayjs from "dayjs";
import {FormatToDateTime} from "../../utils/tools";


interface Prop {
    callID?: string
}


export default function SequenceDiagram(p: Prop) {
    const [loading, setLoading] = useState(true)

    const ssdRef = useRef<HTMLDivElement>(null)

    const [seq, setSeq] = useState<CallRecordEntity[]>([])

    //消息详情弹窗
    const [seqMessageItem, setSeqMessageItem] = useState<CallRecordEntity>()
    const [seqMessageItemModelShow, setSeqMessageItemModelShow] = useState(false)

    const ShowEmpty=()=>{
        if (!loading && seq.length<=0){
            return <Empty/>
        }else {
            return <></>
        }
    }

    useEffect(() => {
        const diagram = new ssd.SequenceDiagram()
        const dom = createSeqHtml(seq)

        diagram.set(dom)
        diagram.addEventListener('click', (e: any) => {
            if (e.type === 'connect') {
                diagram.setHighlight(e.ln)
                setSeqMessageItem(seq[e.ln-1])
                setSeqMessageItemModelShow(true)
            }
        })

        if (ssdRef.current !== null) {
            ssdRef.current.appendChild(diagram.dom())
        }

        //组件卸载时执行的清理逻辑
        return () => {
            diagram.removeAllEventListeners()
            if (ssdRef.current) ssdRef.current.innerHTML = ''
        }

    }, [seq])

    useEffect(() => {
        const searchParams = new URLSearchParams(window.location.search);
        const sipCallId = p.callID || searchParams.get('sip_call_id') || "";

        AppAxios.get<CallRecordDetailsVO>(`/record/details?sip_call_id=` + sipCallId).then(res => {
            setSeq(res.data.data)
            setLoading(false)
        })
    }, [])

    return (
        <div>
            <Spin tip="Loading..." size="large" spinning={loading}>
                <ShowEmpty/>

                <div ref={ssdRef}></div>

                <Modal
                    centered
                    width="80%"
                    open={seqMessageItemModelShow}
                    onCancel={() => {
                        setSeqMessageItemModelShow(false)
                    }}
                    onOk={() => {
                        setSeqMessageItemModelShow(false)
                    }}
                    key={seqMessageItem?.id}
                    title={`${seqMessageItem?.sip_method} ${FormatToDateTime(seqMessageItem?.create_time)}`}>
                    <p>
                        <Tag color="blue">{dayjs(seqMessageItem?.create_time).format('YYYY-MM-DD HH:mm:ss')}</Tag>
                        <Tag color="cyan">{getProtocolName(seqMessageItem?.sip_protocol ? seqMessageItem?.sip_protocol : 0)}</Tag>
                        <Tag color="magenta">length: {seqMessageItem?.raw_msg.length}B</Tag>
                    </p>
                    <div>
                        <pre style={{overflowX: 'scroll'}}>{seqMessageItem?.raw_msg}</pre>
                    </div>
                </Modal>
            </Spin>
        </div>
    )
}
