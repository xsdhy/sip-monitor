import {useEffect, useRef, useState} from 'react'
import {Tag, Spin, Modal, Empty, Tabs} from 'antd'


import mermaid from 'mermaid'
import {createSeqHtml, getProtocolName} from './util'
import {CallRecordDetailsVO, CallRecordEntity, CallRecordRaw, CallRecordRawVO} from '../../@types/entity'
import AppAxios from "../../utils/request";
import dayjs from "dayjs";
import {FormatToDateTime} from "../../utils/tools";


interface Prop {
    callID?: string
}


export default function SequenceDiagram(p: Prop) {
    const [loading, setLoading] = useState(true)

    const recordsRef = useRef<HTMLDivElement>(null)
    const relevantsRef = useRef<HTMLDivElement>(null)

    const [records, setRecords] = useState<CallRecordEntity[]>([])
    const [relevants, setRelevants] = useState<CallRecordEntity[]>([])
    const [activeTabKey, setActiveTabKey] = useState<string>("records")


    //消息详情弹窗
    const [recordItem, setRecordItem] = useState<CallRecordRaw>()
    const [recordItemModelShow, setRecordItemModelShow] = useState(false)

    const ShowEmpty=()=>{
        if (!loading && records.length<=0 && relevants.length<=0){
            return <Empty/>
        }else {
            return <></>
        }
    }

    useEffect(() => {
        if (records.length > 0 && recordsRef.current && activeTabKey === "records") {
            renderMermaidDiagram(recordsRef.current, records);
        }

        return () => {
            if (recordsRef.current) {
                recordsRef.current.innerHTML = ''
            }
        }
    }, [records, activeTabKey])
    
    useEffect(() => {
        if (relevants.length > 0 && relevantsRef.current && activeTabKey === "relevants") {
            renderMermaidDiagram(relevantsRef.current, relevants);
        }

        return () => {
            if (relevantsRef.current) {
                relevantsRef.current.innerHTML = ''
            }
        }
    }, [relevants, activeTabKey])

    const handleTabChange = (key: string) => {
        setActiveTabKey(key);
        
        // Re-render diagrams when tab changes
        setTimeout(() => {
            if (key === "records" && records.length > 0 && recordsRef.current) {
                renderMermaidDiagram(recordsRef.current, records);
            } else if (key === "relevants" && relevants.length > 0 && relevantsRef.current) {
                renderMermaidDiagram(relevantsRef.current, relevants);
            }
        }, 100); // Small delay to ensure DOM is ready
    };

    const renderMermaidDiagram = (container: HTMLDivElement, data: CallRecordEntity[]) => {
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
        container.innerHTML = ''
        
        // 创建新的图表容器
        const chartContainer = document.createElement('div')
        chartContainer.className = 'mermaid'
        chartContainer.innerHTML = createSeqHtml(data)
        container.appendChild(chartContainer)
        
        // 渲染图表
        mermaid.initialize({
            theme: 'base',
            sequence: { showSequenceNumbers: true }
        })
        mermaid.run({querySelector: '.mermaid'})

        // 添加点击事件
        chartContainer.addEventListener('click', (e) => {
            const target = e.target as HTMLElement
            // 查找最近的 text 元素
            const textElement = target.closest('text');
            if (textElement) {
                const messageText = textElement.textContent || '';
                // 在序列中查找对应的消息
                const messageIndex = data.findIndex((item, index) => {
                    const expectedText = `${item.method} `;
                    if (item.method === 'INVITE') {
                        return messageText.includes(`${item.method} ${item.from_user} -> ${item.to_user}`);
                    }
                    return messageText.includes(expectedText);
                });

                // 根据item.id 获取record_raw
               AppAxios.get<CallRecordRawVO>(`/record/raw/` + data[messageIndex].id).then(res=>{
                if (res.data.code === 200) {
                    setRecordItem(res.data.data);
                    setRecordItemModelShow(true);
                }
                })
                
         
            }
        });
    }

    useEffect(() => {
        const searchParams = new URLSearchParams(window.location.search);
        const sipCallId = p.callID || searchParams.get('sip_call_id') || "";

        AppAxios.get<CallRecordDetailsVO>(`/record/details?sip_call_id=` + sipCallId).then(res => {
            if (res.data.code === 200) {
                if (res.data.data.records.length > 0) {
                    setRecords(res.data.data.records)
                }
                if (res.data.data.relevants.length > 0 && res.data.data.records.length !== res.data.data.relevants.length) {
                    setRelevants(res.data.data.relevants)
                }
                setLoading(false)
            }
        })
    }, [])

    return (
        <div>
            <Spin tip="Loading..." size="large" spinning={loading}>
                <ShowEmpty/>

                <Tabs defaultActiveKey="records" activeKey={activeTabKey} onChange={handleTabChange}>
                    <Tabs.TabPane tab="当前会话" key="records">
                        <div ref={recordsRef}></div>
                    </Tabs.TabPane>
                    <Tabs.TabPane tab="相关会话" key="relevants">
                        <div ref={relevantsRef}></div>
                    </Tabs.TabPane>
                </Tabs>

                <Modal
                    centered
                    width="80%"
                    open={recordItemModelShow}
                    onCancel={() => {
                        setRecordItemModelShow(false)
                    }}
                    onOk={() => {
                        setRecordItemModelShow(false)
                    }}
                    key={recordItem?.id}
                    title={`信令详情`}>
                    <div>

                        <pre style={{overflowX: 'scroll'}}>{recordItem?.raw}</pre>
                    </div>
                </Modal>
            </Spin>
        </div>
    )
}
