import {useEffect, useRef, useState} from 'react'
import {Tag, Spin, Modal, Empty} from 'antd'


import mermaid from 'mermaid'
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
        if (seq.length > 0 && ssdRef.current) {
            // 初始化 mermaid
            mermaid.initialize({
                startOnLoad: true,
                theme: 'default',
                sequence: {
                    diagramMarginX: 50,
                    diagramMarginY: 10,
                    actorMargin: 50,
                    width: 150,
                    height: 65,
                    boxMargin: 10,
                    boxTextMargin: 5,
                    noteMargin: 10,
                    messageMargin: 35
                }
            })

            // 清除之前的内容
            ssdRef.current.innerHTML = ''
            
            // 创建新的图表容器
            const container = document.createElement('div')
            container.className = 'mermaid'
            container.innerHTML = createSeqHtml(seq)
            ssdRef.current.appendChild(container)
            
            // 渲染图表
            // mermaid.contentLoaded = false
            mermaid.init({ sequence: { showSequenceNumbers: true } }, '.mermaid')

            // 添加点击事件
            container.addEventListener('click', (e) => {
                const target = e.target as HTMLElement
                if (target.closest('.messageText')) {
                    const index = parseInt(target.getAttribute('data-index') || '0')
                    setSeqMessageItem(seq[index])
                    setSeqMessageItemModelShow(true)
                }
            })
        }

        return () => {
            if (ssdRef.current) {
                ssdRef.current.innerHTML = ''
            }
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
