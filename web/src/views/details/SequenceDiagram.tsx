import { useEffect, useRef, useState } from "react";
import { Spin, Modal, Empty, Tabs } from "antd";
import * as echarts from 'echarts';

import mermaid from "mermaid";
import { createSeqHtml } from "./util";
import {
  CallRecordEntity,
  CallRecordRaw,
  RtcpReportRaw
} from "../../@types/entity";
import { callApi } from "@/apis/api";

interface Prop {
  callID?: string;
}

export default function SequenceDiagram(p: Prop) {
  const [loading, setLoading] = useState(true);

  const recordsRef = useRef<HTMLDivElement>(null);
  const relevantsRef = useRef<HTMLDivElement>(null);
  const rtcpRef = useRef<HTMLDivElement>(null);

  const [records, setRecords] = useState<CallRecordEntity[]>([]);
  const [relevants, setRelevants] = useState<CallRecordEntity[]>([]);
  const [rtcpPackets, setRtcpPackets] = useState<RtcpReportRaw[]>([]);
  const [activeTabKey, setActiveTabKey] = useState<string>("records");

  //消息详情弹窗
  const [recordItem, setRecordItem] = useState<CallRecordRaw>();
  const [recordItemModelShow, setRecordItemModelShow] = useState(false);

  //RTCP详情弹窗
  const [rtcpItem, setRtcpItem] = useState<RtcpReportRaw>();
  const [rtcpItemModalShow, setRtcpItemModalShow] = useState(false);

  const ShowEmpty = () => {
    if (!loading && records.length <= 0 && relevants.length <= 0 && rtcpPackets.length <= 0) {
      return <Empty />;
    } else {
      return <></>;
    }
  };

  useEffect(() => {
    if (
      records.length > 0 &&
      recordsRef.current &&
      activeTabKey === "records"
    ) {
      renderMermaidDiagram(recordsRef.current, records);
    }

    return () => {
      if (recordsRef.current) {
        recordsRef.current.innerHTML = "";
      }
    };
  }, [records, activeTabKey]);

  useEffect(() => {
    if (
      relevants.length > 0 &&
      relevantsRef.current &&
      activeTabKey === "relevants"
    ) {
      renderMermaidDiagram(relevantsRef.current, relevants);
    }

    return () => {
      if (relevantsRef.current) {
        relevantsRef.current.innerHTML = "";
      }
    };
  }, [relevants, activeTabKey]);

  useEffect(() => {
    if (
      rtcpPackets.length > 0 &&
      rtcpRef.current &&
      activeTabKey === "rtcp"
    ) {
      renderRtcpChart(rtcpRef.current, rtcpPackets);
    }

    return () => {
      if (rtcpRef.current) {
        rtcpRef.current.innerHTML = "";
      }
    };
  }, [rtcpPackets, activeTabKey]);

  const handleTabChange = (key: string) => {
    setActiveTabKey(key);

    // Re-render diagrams when tab changes
    setTimeout(() => {
      if (key === "records" && records.length > 0 && recordsRef.current) {
        renderMermaidDiagram(recordsRef.current, records);
      } else if (
        key === "relevants" &&
        relevants.length > 0 &&
        relevantsRef.current
      ) {
        renderMermaidDiagram(relevantsRef.current, relevants);
      } else if (
        key === "rtcp" &&
        rtcpPackets.length > 0 &&
        rtcpRef.current
      ) {
        renderRtcpChart(rtcpRef.current, rtcpPackets);
      }
    }, 100); // Small delay to ensure DOM is ready
  };

  const renderMermaidDiagram = (
    container: HTMLDivElement,
    data: CallRecordEntity[]
  ) => {
    // 初始化 mermaid
    mermaid.initialize({
      startOnLoad: true,
      theme: "default",
      sequence: {
        diagramMarginX: 50,
        diagramMarginY: 10,
        actorMargin: 50,
        width: 150,
        height: 65,
        boxMargin: 10,
        boxTextMargin: 5,
        noteMargin: 10,
        messageMargin: 35,
      },
    });

    // 清除之前的内容
    container.innerHTML = "";

    // 创建新的图表容器
    const chartContainer = document.createElement("div");
    chartContainer.className = "mermaid";
    chartContainer.innerHTML = createSeqHtml(data);
    container.appendChild(chartContainer);

    // 渲染图表
    mermaid.initialize({
      theme: "base",
      sequence: { showSequenceNumbers: true },
    });
    mermaid.run({ querySelector: ".mermaid" });

    // 添加点击事件
    chartContainer.addEventListener("click", (e) => {
      const target = e.target as HTMLElement;
      // 查找最近的 text 元素
      const textElement = target.closest("text");
      if (textElement) {
        const messageText = textElement.textContent || "";
        // 在序列中查找对应的消息
        const messageIndex = data.findIndex((item) => {
          const expectedText = `${item.method} `;
          if (item.method === "INVITE") {
            return messageText.includes(
              `${item.method} ${item.from_user} -> ${item.to_user}`
            );
          }
          return messageText.includes(expectedText);
        });

        // 根据item.id 获取record_raw
        callApi.getCallRecordRaw(data[messageIndex].id.toString()).then((res) => {
          setRecordItem(res.data);
          setRecordItemModelShow(true);
        });
      }
    });
  };

  const renderRtcpChart = (
    container: HTMLDivElement,
    data: RtcpReportRaw[]
  ) => {
    // 清除之前的内容
    container.innerHTML = "";

    // 创建RTCP图表容器
    const chartContainer = document.createElement("div");
    chartContainer.style.width = "100%";
    chartContainer.style.height = "400px";
    container.appendChild(chartContainer);

    // 解析RTCP数据
    const parsedData = data.map(item => {
      try {
        const raw = JSON.parse(item.raw);
        return {
          ...item,
          parsedData: raw,
          time: new Date(item.create_time),
          source: `${item.src_addr} → ${item.dst_addr}`
        };
      } catch (e) {
        console.error("RTCP数据解析错误", e);
        return null;
      }
    }).filter(Boolean);

    // 按源地址和目标地址组织数据
    const connectionGroups: Record<string, any[]> = {};
    parsedData.forEach(item => {
      if (!item) return;
      const key = item.source;
      if (!connectionGroups[key]) {
        connectionGroups[key] = [];
      }
      connectionGroups[key].push(item);
    });

    // 创建选项卡容器
    const tabContainer = document.createElement("div");
    tabContainer.style.marginBottom = "20px";
    container.insertBefore(tabContainer, chartContainer);

    // 创建连接列表
    let activeConnection = '';
    Object.keys(connectionGroups).forEach((connection, index) => {
      const btn = document.createElement("button");
      btn.innerText = connection;
      btn.style.margin = "5px";
      btn.style.padding = "5px 10px";
      btn.style.border = "1px solid #ccc";
      btn.style.borderRadius = "4px";
      btn.style.cursor = "pointer";
      
      if (index === 0) {
        activeConnection = connection;
        btn.style.backgroundColor = "#1890ff";
        btn.style.color = "white";
      }
      
      btn.onclick = () => {
        // 重置所有按钮样式
        tabContainer.querySelectorAll("button").forEach(b => {
          b.style.backgroundColor = "";
          b.style.color = "";
        });
        
        // 设置当前按钮样式
        btn.style.backgroundColor = "#1890ff";
        btn.style.color = "white";
        
        // 更新图表
        updateChart(connectionGroups[connection]);
      };
      
      tabContainer.appendChild(btn);
    });

    // 初始化echarts实例
    const chart = echarts.init(chartContainer);
    
    // 设置点击事件
    chart.on('click', (params: any) => {
      const index = params.dataIndex;
      const connectionData = connectionGroups[activeConnection];
      if (connectionData && connectionData[index]) {
        setRtcpItem(connectionData[index]);
        setRtcpItemModalShow(true);
      }
    });

    // 更新图表函数
    const updateChart = (connectionData: any[]) => {
      // 处理时间数据
      const times = connectionData.map(item => {
        const date = new Date(item.create_time);
        return `${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}`;
      });

      // 准备数据
      const packetData = connectionData.map(item => item.parsedData.sender_information.packets || 0);
      const octetData = connectionData.map(item => item.parsedData.sender_information.octets || 0);
      const jitterData = connectionData.map(item => {
        const reportBlocks = item.parsedData.report_blocks;
        return reportBlocks && reportBlocks.length > 0 ? reportBlocks[0].ia_jitter || 0 : 0;
      });
      const lostData = connectionData.map(item => {
        const reportBlocks = item.parsedData.report_blocks;
        return reportBlocks && reportBlocks.length > 0 ? reportBlocks[0].fraction_lost || 0 : 0;
      });

      // 设置图表选项
      const option = {
        title: {
          text: 'RTCP报告',
          left: 'center'
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'cross',
            label: {
              backgroundColor: '#6a7985'
            }
          }
        },
        legend: {
          data: ['packets', 'octets', 'ia_jitter', 'fraction_lost'],
          top: 30
        },
        toolbox: {
          feature: {
            saveAsImage: {}
          }
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: [
          {
            type: 'category',
            boundaryGap: false,
            data: times
          }
        ],
        yAxis: [
          {
            type: 'value',
            name: 'packets/octets'
          },
          {
            type: 'value',
            name: 'jitter/lost',
            axisLabel: {
              formatter: '{value}'
            }
          }
        ],
        series: [
          {
            name: 'packets',
            type: 'bar',
            data: packetData,
            itemStyle: {
              color: '#FF9A9E'
            }
          },
          {
            name: 'octets',
            type: 'line',
            data: octetData,
            itemStyle: {
              color: '#67C23A'
            }
          },
          {
            name: 'ia_jitter',
            type: 'line',
            yAxisIndex: 1,
            data: jitterData,
            itemStyle: {
              color: '#409EFF'
            }
          },
          {
            name: 'fraction_lost',
            type: 'line',
            yAxisIndex: 1,
            data: lostData,
            itemStyle: {
              color: '#E6A23C'
            }
          }
        ]
      };

      // 设置图表
      chart.setOption(option);
    };

    // 初始显示第一个连接的数据
    if (activeConnection && connectionGroups[activeConnection]) {
      updateChart(connectionGroups[activeConnection]);
    }

    // 窗口大小变化时调整图表大小
    window.addEventListener('resize', () => {
      chart.resize();
    });
  };

  useEffect(() => {
    const searchParams = new URLSearchParams(window.location.search);
    const sipCallId = p.callID || searchParams.get("sip_call_id") || "";

    callApi.getCallDetail(sipCallId).then((res) => {
      if (res.data.records.length > 0) {
        setRecords(res.data.records);
      }
      if (
        res.data.relevants.length > 0 &&
        res.data.records.length !== res.data.relevants.length
      ) {
        setRelevants(res.data.relevants);
      }
      if (res.data.rtcp_packets && res.data.rtcp_packets.length > 0) {
        setRtcpPackets(res.data.rtcp_packets);
      }
      setLoading(false);
    });
  }, []);

  return (
    <div>
      <Spin tip="Loading..." size="large" spinning={loading}>
        <ShowEmpty />

        <Tabs
          defaultActiveKey="records"
          activeKey={activeTabKey}
          onChange={handleTabChange}
        >
          <Tabs.TabPane tab="当前会话" key="records">
            <div ref={recordsRef}></div>
          </Tabs.TabPane>
          <Tabs.TabPane tab="相关会话" key="relevants">
            <div ref={relevantsRef}></div>
          </Tabs.TabPane>
          <Tabs.TabPane tab="RTCP" key="rtcp">
            <div ref={rtcpRef}></div>
          </Tabs.TabPane>
        </Tabs>

        <Modal
          centered
          width="80%"
          open={recordItemModelShow}
          onCancel={() => {
            setRecordItemModelShow(false);
          }}
          onOk={() => {
            setRecordItemModelShow(false);
          }}
          key={recordItem?.id}
          title={`信令详情`}
        >
          <div>
            <pre style={{ overflowX: "scroll" }}>{recordItem?.raw}</pre>
          </div>
        </Modal>

        <Modal
          centered
          width="80%"
          open={rtcpItemModalShow}
          onCancel={() => {
            setRtcpItemModalShow(false);
          }}
          onOk={() => {
            setRtcpItemModalShow(false);
          }}
          key={rtcpItem?.id}
          title={`RTCP详情 (${rtcpItem?.src_addr} → ${rtcpItem?.dst_addr})`}
        >
          <div>
            <pre style={{ overflowX: "scroll" }}>{
              rtcpItem ? JSON.stringify(JSON.parse(rtcpItem.raw), null, 2) : ''
            }</pre>
          </div>
        </Modal>
      </Spin>
    </div>
  );
}
