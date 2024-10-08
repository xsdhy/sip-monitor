import {useEffect, useMemo, useRef, useState} from 'react'
import {Tag, Spin, Modal, Empty} from 'antd'
import './Sequence.css'


import * as ssd from 'svg-sequence-diagram'
import {createSeqHtml, getProtocolName} from './util'
import {CallRecordDetailsVO, CallRecordEntity, SipCallFlowDiagramData} from '../../@types/entity'
import AppAxios from "../../utils/request";
import dayjs from "dayjs";
import {FormatToDateTime} from "../../utils/tools";


interface Prop {
    callID?: string
}


export default function SequenceDiagram(p: Prop) {
    const [loading, setLoading] = useState(true)


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

    const SipCallFlowDiagram = () => {
        const columnWidth = 200;
        const rowHeight = 60;


        const data = useMemo(():SipCallFlowDiagramData => {
            const uniqueAddresses = new Set();
            seq.forEach(item => {
                uniqueAddresses.add(item.src);
                uniqueAddresses.add(item.dst);
            });
            const cols:any[] = Array.from(uniqueAddresses);

            return { columns: cols, messages: seq };
        }, [seq]);

        const svgWidth = columnWidth * data.columns.length;
        const svgHeight = rowHeight * data.messages.length;

        // Function to generate arrow path
        const generateArrowPath = (startX:number, endX:number, y:number) => {
            const arrowSize = 10;
            const direction = startX < endX ? 1 : -1;
            const arrowTip = endX - direction * arrowSize;
            return `M${startX},${y} L${arrowTip},${y} M${arrowTip},${y-arrowSize/2} L${endX},${y} L${arrowTip},${y+arrowSize/2}`;
        };


        return (
            <div className="w-full h-full overflow-auto">
                <svg width={svgWidth} height={svgHeight + 50}>
                    
                    {/* Draw vertical lines for each column */}
                    {data.columns.map((col, index) => (
                        // 每一行是IP:PORT,如果IP相同，则线之间的距离短一些
                        // 计算每个IP的出现次数
                        
                        <line
                            x1={columnWidth * (index + 0.5)}
                            y1={0}
                            x2={columnWidth * (index + 0.5)}
                            y2={svgHeight}
                            stroke="#e0e0e0"
                            strokeDasharray="5,5"
                        />
                    ))}

                    {/* Draw messages and arrows */}
                    {data.messages.map((msg, index) => {
                        const startX = columnWidth * (data.columns.indexOf(msg.src) + 0.5);
                        const endX = columnWidth * (data.columns.indexOf(msg.dst) + 0.5);
                        const y = rowHeight * (index + 0.5) + 50; // Add 50px for column headers

                        return (
                            <g key={index}>
                                {/* Arrow */}
                                <path
                                    d={generateArrowPath(startX, endX, y)}
                                    fill="none"
                                    stroke="#000"
                                    strokeWidth={2}
                                />
                                {/* Message text */}
                                <text
                                    x={(startX + endX) / 2}
                                    y={y - 15}
                                    textAnchor="middle"
                                    fill="#000"
                                    fontSize={12}
                                >
                                    {msg.method}
                                </text>
                                {/* Time text */}
                                <text
                                    x={(startX + endX) / 2}
                                    y={y + 20}
                                    textAnchor="middle"
                                    fill="#666"
                                    fontSize={10}
                                >
                                    {msg.create_time}
                                </text>
                            </g>
                        );
                    })}

                    {/* Column headers */}
                    {data.columns.map((col:string, index:number) => (
                        <text
                            x={columnWidth * (index + 0.5)}
                            y={30}
                            textAnchor="middle"
                            fill="#333"
                            fontSize={14}
                            fontWeight="bold">
                            {col}
                        </text>
                    ))}
                </svg>
            </div>
        );
    };



    useEffect(() => {
        const searchParams = new URLSearchParams(window.location.search);
        const sipCallId = p.callID || searchParams.get('call_id') || "";

        AppAxios.get<CallRecordDetailsVO>(`/record/details?call_id=` + sipCallId).then(res => {
            setSeq(res.data.data)
            setLoading(false)
        })
    }, [])

    return (
        <div>
            <Spin tip="Loading..." size="large" spinning={loading}>
                <ShowEmpty/>

                {/*<div ref={ssdRef}></div>*/}

                <div className="w-full h-screen p-4">
                    <SipCallFlowDiagram />
                </div>

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
                    title={` ${FormatToDateTime(seqMessageItem?.create_time)}`}>
                    <p>
                        <Tag color="blue">{dayjs(seqMessageItem?.create_time).format('YYYY-MM-DD HH:mm:ss')}</Tag>

                        <Tag color="magenta">length: {seqMessageItem?.body.length}B</Tag>
                    </p>
                    <div>
                        <pre style={{overflowX: 'scroll'}}>{seqMessageItem?.body}</pre>
                    </div>
                </Modal>
            </Spin>
        </div>
    )
}
