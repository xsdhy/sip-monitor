import { useEffect, useRef, useState } from "react";
import { Spin, Modal, Empty, Tabs } from "antd";
import * as echarts from 'echarts';

import mermaid from "mermaid";
import { createSeqHtml } from "./util";
import {
  CallRecordEntity,
  CallRecordRaw,
  RtcpReport,
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
  const rtcpReportRef = useRef<HTMLDivElement>(null);
  const localRtcpAnalysisRef = useRef<HTMLDivElement>(null);
  const [records, setRecords] = useState<CallRecordEntity[]>([]);
  const [relevants, setRelevants] = useState<CallRecordEntity[]>([]);
  const [rtcpPackets, setRtcpPackets] = useState<RtcpReportRaw[]>([]);
  const [rtcpReport, setRtcpReport] = useState<RtcpReport>();
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

  useEffect(() => {
    if (rtcpReport && rtcpReportRef.current && activeTabKey === "rtcp_report") {
      renderRtcpReport(rtcpReportRef.current, rtcpReport);
    }

    return () => {
      if (rtcpReportRef.current) {
        rtcpReportRef.current.innerHTML = "";
      }
    };
  }, [rtcpReport, activeTabKey]);

  useEffect(() => {
    if (rtcpPackets.length > 0 && localRtcpAnalysisRef.current && activeTabKey === "local_rtcp_analysis") {
      renderLocalRtcpAnalysis(localRtcpAnalysisRef.current, rtcpPackets);
    }

    return () => {
      if (localRtcpAnalysisRef.current) {
        localRtcpAnalysisRef.current.innerHTML = "";
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
      } else if (
        key === "rtcp_report" &&
        rtcpReport &&
        rtcpReportRef.current
      ) {
        renderRtcpReport(rtcpReportRef.current, rtcpReport);
      } else if (
        key === "local_rtcp_analysis" &&
        rtcpPackets.length > 0 &&
        localRtcpAnalysisRef.current
      ) {
        renderLocalRtcpAnalysis(localRtcpAnalysisRef.current, rtcpPackets);
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

  const renderRtcpReport = (container: HTMLDivElement |  null, data: RtcpReport) => {
    if (!container) return;

    // 清除之前的内容
    container.innerHTML = "";

    // 创建RTCP报告表格容器
    const tableContainer = document.createElement("div");
    tableContainer.style.width = "100%";
    tableContainer.style.padding = "20px";
    container.appendChild(tableContainer);

    // 创建表格
    const table = document.createElement("table");
    table.style.width = "100%";
    table.style.borderCollapse = "collapse";
    table.style.marginTop = "20px";
    table.style.boxShadow = "0 0 10px rgba(0, 0, 0, 0.1)";
    tableContainer.appendChild(table);

    // 创建表头
    const thead = document.createElement("thead");
    table.appendChild(thead);
    
    const headerRow = document.createElement("tr");
    thead.appendChild(headerRow);
    
    const headers = ["指标", "Aleg", "Bleg"];
    headers.forEach(header => {
      const th = document.createElement("th");
      th.textContent = header;
      th.style.padding = "12px 15px";
      th.style.backgroundColor = "#f8f9fa";
      th.style.borderBottom = "2px solid #ddd";
      th.style.textAlign = "left";
      headerRow.appendChild(th);
    });

    // 创建表格内容
    const tbody = document.createElement("tbody");
    table.appendChild(tbody);
    
    // 定义要显示的指标
    const metrics = [
      { name: "MOS", aleg: data.aleg_mos, bleg: data.bleg_mos },
      { name: "丢包数", aleg: data.aleg_packet_lost, bleg: data.bleg_packet_lost },
      { name: "数据包总数", aleg: data.aleg_packet_count, bleg: data.bleg_packet_count },
      { name: "丢包率", aleg: data.aleg_packet_lost_rate, bleg: data.bleg_packet_lost_rate },
      { name: "平均抖动", aleg: data.aleg_jitter_avg, bleg: data.bleg_jitter_avg },
      { name: "最大抖动", aleg: data.aleg_jitter_max, bleg: data.bleg_jitter_max },
      { name: "平均延迟", aleg: data.aleg_delay_avg, bleg: data.bleg_delay_avg },
      { name: "最大延迟", aleg: data.aleg_delay_max, bleg: data.bleg_delay_max }
    ];
    
    // 添加行
    metrics.forEach((metric, index) => {
      const row = document.createElement("tr");
      row.style.backgroundColor = index % 2 === 0 ? "#fff" : "#f8f9fa";
      
      const metricCell = document.createElement("td");
      metricCell.textContent = metric.name;
      metricCell.style.padding = "12px 15px";
      metricCell.style.borderBottom = "1px solid #ddd";
      metricCell.style.fontWeight = "bold";
      row.appendChild(metricCell);
      
      const alegCell = document.createElement("td");
      alegCell.textContent = metric.aleg?.toString() || "N/A";
      alegCell.style.padding = "12px 15px";
      alegCell.style.borderBottom = "1px solid #ddd";
      row.appendChild(alegCell);
      
      const blegCell = document.createElement("td");
      blegCell.textContent = metric.bleg?.toString() || "N/A";
      blegCell.style.padding = "12px 15px";
      blegCell.style.borderBottom = "1px solid #ddd";
      row.appendChild(blegCell);
      
      tbody.appendChild(row);
    });
  };

  const renderLocalRtcpAnalysis = (container: HTMLDivElement, data: RtcpReportRaw[]) => {
    // 清除之前的内容
    container.innerHTML = "";

    // 创建分析容器
    const analysisContainer = document.createElement("div");
    analysisContainer.style.width = "100%";
    analysisContainer.style.padding = "20px";
    container.appendChild(analysisContainer);

    // 添加说明文本
    const infoText = document.createElement("div");
    infoText.innerHTML = `
      <div style="margin-bottom: 20px; padding: 10px; background-color: #f9f9fa; border-radius: 4px; border-left: 4px solid #1890ff;">
        <p style="margin: 0; padding: 0;"><strong>关于RTCP分析（仅供参考）：</strong></p>
        <p style="margin: 5px 0 0; padding: 0;">在SIP语音通话场景下，基于HEP捕获的RTCP数据包</p>
        <p style="margin: 5px 0 0; padding: 0;">MOS分数范围为1-4.5，其中4.3-4.5为优秀，4.0-4.3为良好，3.5-4.0为一般，小于3.5为较差。</p>
      </div>
    `;
    analysisContainer.appendChild(infoText);

    // 解析RTCP数据
    const parsedData = data.map(item => {
      try {
        const raw = JSON.parse(item.raw);
        return {
          ...item,
          parsedData: raw,
          time: new Date(item.create_time),
          source: `${item.src_addr} → ${item.dst_addr}`,
          key: `${item.src_addr}_${item.dst_addr}`
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
      const key = item.key;
      if (!connectionGroups[key]) {
        connectionGroups[key] = [];
      }
      connectionGroups[key].push(item);
    });

    // 为每个连接计算质量指标
    Object.keys(connectionGroups).forEach(connectionKey => {
      const connectionData = connectionGroups[connectionKey];
      const sourceDestPair = connectionData[0].source;
      
      // 创建连接卡片
      const connectionCard = document.createElement("div");
      connectionCard.style.marginBottom = "30px";
      connectionCard.style.padding = "20px";
      connectionCard.style.borderRadius = "8px";
      connectionCard.style.boxShadow = "0 4px 12px rgba(0, 0, 0, 0.1)";
      connectionCard.style.backgroundColor = "#fff";
      analysisContainer.appendChild(connectionCard);
      
      // 添加标题
      const title = document.createElement("h3");
      title.textContent = `连接: ${sourceDestPair}`;
      title.style.marginBottom = "15px";
      title.style.borderBottom = "1px solid #eee";
      title.style.paddingBottom = "10px";
      connectionCard.appendChild(title);
      
      // 计算质量指标
      const qualityMetrics = calculateQualityMetrics(connectionData);
      
      // 创建指标网格
      const metricsGrid = document.createElement("div");
      metricsGrid.style.display = "grid";
      metricsGrid.style.gridTemplateColumns = "repeat(2, 1fr)";
      metricsGrid.style.gap = "15px";
      connectionCard.appendChild(metricsGrid);
      
      // 添加指标项
      addMetricItem(metricsGrid, "平均丢包率", `${(qualityMetrics.avgPacketLoss * 100).toFixed(2)}%`, qualityMetrics.avgPacketLoss > 0.03);
      addMetricItem(metricsGrid, "最大丢包率", `${(qualityMetrics.maxPacketLoss * 100).toFixed(2)}%`, qualityMetrics.maxPacketLoss > 0.05);
      addMetricItem(metricsGrid, "平均抖动", `${qualityMetrics.avgJitter.toFixed(2)} ms`, qualityMetrics.avgJitter > 30);
      addMetricItem(metricsGrid, "最大抖动", `${qualityMetrics.maxJitter.toFixed(2)} ms`, qualityMetrics.maxJitter > 50);
      
      if (qualityMetrics.avgRoundTripDelay > 0) {
        addMetricItem(metricsGrid, "平均往返延迟", `${qualityMetrics.avgRoundTripDelay.toFixed(2)} ms`, qualityMetrics.avgRoundTripDelay > 300);
        addMetricItem(metricsGrid, "最大往返延迟", `${qualityMetrics.maxRoundTripDelay.toFixed(2)} ms`, qualityMetrics.maxRoundTripDelay > 500);
      }
      
      // 计算MOS值
      const mosValue = calculateMOS(qualityMetrics.avgJitter, qualityMetrics.avgPacketLoss, qualityMetrics.avgRoundTripDelay);
      addMetricItem(metricsGrid, "MOS评分 (1-5)", mosValue.toFixed(2), mosValue < 3.5);
      
      // 添加质量评估
      const qualityAssessment = document.createElement("div");
      qualityAssessment.style.marginTop = "20px";
      qualityAssessment.style.padding = "10px";
      qualityAssessment.style.borderRadius = "4px";
      
      let assessmentText = "";
      let bgColor = "";
      
      if (mosValue >= 4.3) {
        assessmentText = "通话质量: 优秀";
        bgColor = "#f0f9eb";
        qualityAssessment.style.color = "#67c23a";
        qualityAssessment.style.border = "1px solid #e1f3d8";
      } else if (mosValue >= 3.5) {
        assessmentText = "通话质量: 良好";
        bgColor = "#f4f4f5";
        qualityAssessment.style.color = "#909399";
        qualityAssessment.style.border = "1px solid #e9e9eb";
      } else if (mosValue >= 3.0) {
        assessmentText = "通话质量: 一般";
        bgColor = "#fdf6ec";
        qualityAssessment.style.color = "#e6a23c";
        qualityAssessment.style.border = "1px solid #faecd8";
      } else {
        assessmentText = "通话质量: 较差";
        bgColor = "#fef0f0";
        qualityAssessment.style.color = "#f56c6c";
        qualityAssessment.style.border = "1px solid #fde2e2";
      }
      
      qualityAssessment.textContent = assessmentText;
      qualityAssessment.style.backgroundColor = bgColor;
      qualityAssessment.style.fontWeight = "bold";
      qualityAssessment.style.textAlign = "center";
      connectionCard.appendChild(qualityAssessment);
    });

    // 如果没有数据，显示提示
    if (Object.keys(connectionGroups).length === 0) {
      const noDataMsg = document.createElement("div");
      noDataMsg.textContent = "没有足够的RTCP数据进行分析";
      noDataMsg.style.padding = "20px";
      noDataMsg.style.textAlign = "center";
      noDataMsg.style.color = "#999";
      analysisContainer.appendChild(noDataMsg);
    }

    // 增加原始数据的显示
    if (data.length > 0) {
      const rawDataContainer = document.createElement("div");
      rawDataContainer.style.marginTop = "30px";
      rawDataContainer.style.padding = "20px";
      rawDataContainer.style.backgroundColor = "#f8f9fa";
      rawDataContainer.style.borderRadius = "8px";
      rawDataContainer.style.border = "1px solid #e9e9eb";
      analysisContainer.appendChild(rawDataContainer);
      
      // 添加标题
      const rawDataTitle = document.createElement("div");
      rawDataTitle.textContent = "原始RTCP数据";
      rawDataTitle.style.fontWeight = "bold";
      rawDataTitle.style.marginBottom = "10px";
      rawDataTitle.style.fontSize = "16px";
      rawDataTitle.style.display = "flex";
      rawDataTitle.style.justifyContent = "space-between";
      rawDataTitle.style.alignItems = "center";
      rawDataContainer.appendChild(rawDataTitle);
      

      
      // 添加数据内容区域
      const rawDataContent = document.createElement("pre");
      rawDataContent.style.marginTop = "10px";
      rawDataContent.style.overflowX = "auto";
      rawDataContent.style.padding = "10px";
      rawDataContent.style.backgroundColor = "#f1f1f1";
      rawDataContent.style.borderRadius = "4px";
      rawDataContent.style.fontSize = "12px";
      rawDataContent.style.lineHeight = "1.4";
      rawDataContent.style.maxHeight = "500px";


      if (data.length > 0) {
        //遍历data，取其中的RAW，并拼接成一个字符串
        const rawData = data.map(item => item.raw).join("\n");
        rawDataContent.textContent = rawData;
        rawDataContainer.appendChild(rawDataContent);
      }
      

    }
  };

  // 辅助函数：添加指标项
  const addMetricItem = (container: HTMLElement, label: string, value: string, isWarning: boolean) => {
    const metricItem = document.createElement("div");
    metricItem.style.padding = "10px";
    metricItem.style.backgroundColor = "#f9f9f9";
    metricItem.style.borderRadius = "4px";
    metricItem.style.display = "flex";
    metricItem.style.justifyContent = "space-between";
    metricItem.style.alignItems = "center";
    
    const labelElem = document.createElement("span");
    labelElem.textContent = label;
    labelElem.style.fontWeight = "500";
    
    const valueElem = document.createElement("span");
    valueElem.textContent = value;
    valueElem.style.fontWeight = "bold";
    valueElem.style.color = isWarning ? "#f56c6c" : "#67c23a";
    
    metricItem.appendChild(labelElem);
    metricItem.appendChild(valueElem);
    container.appendChild(metricItem);
  };

  // 计算质量指标
  const calculateQualityMetrics = (data: any[]) => {
    let totalJitter = 0;
    let maxJitter = 0;
    let totalPacketLoss = 0;
    let maxPacketLoss = 0;
    let totalRoundTripDelay = 0;
    let maxRoundTripDelay = 0;
    let delayCount = 0;
    
    data.forEach(item => {
      const reportBlocks = item.parsedData.report_blocks || [];
      if (reportBlocks.length > 0) {
        // 抖动 - 从RTP时间戳单位转换为毫秒
        // 大多数SIP通话使用G.711编解码器，采样率为8000Hz
        // 抖动(ms) = ia_jitter / (采样率 / 1000)
        const rawJitter = reportBlocks[0].ia_jitter || 0;
        const sampleRate = 8000; // G.711默认采样率
        const jitter = rawJitter / (sampleRate / 1000); // 转换为毫秒
        
        totalJitter += jitter;
        maxJitter = Math.max(maxJitter, jitter);
        
        // 丢包率 (0-255, 转换为0-1)
        const fractionLost = (reportBlocks[0].fraction_lost || 0) / 255;
        totalPacketLoss += fractionLost;
        maxPacketLoss = Math.max(maxPacketLoss, fractionLost);
        
        // 往返延迟计算
        // 如果同时有LSR和DLSR值，可以计算往返延迟
        if (reportBlocks[0].lsr && reportBlocks[0].dlsr) {
          // DLSR是以1/65536秒为单位的
          const dlsr = (reportBlocks[0].dlsr || 0) / 65536 * 1000; // 转换为毫秒
          totalRoundTripDelay += dlsr;
          maxRoundTripDelay = Math.max(maxRoundTripDelay, dlsr);
          delayCount++;
        }
      }
    });
    
    const count = data.length;
    return {
      avgJitter: count > 0 ? totalJitter / count : 0,
      maxJitter,
      avgPacketLoss: count > 0 ? totalPacketLoss / count : 0,
      maxPacketLoss,
      avgRoundTripDelay: delayCount > 0 ? totalRoundTripDelay / delayCount : 0,
      maxRoundTripDelay
    };
  };

  // 计算MOS值 (使用E-model简化版)
  const calculateMOS = (jitter: number, packetLoss: number, delay: number) => {
    // 基础R值
    let R = 93.2;
    
    // 延迟影响 (IDT)
    // 单向延迟约为往返延迟的一半
    const oneWayDelay = delay / 2;
    const delayFactor = oneWayDelay > 160 ? 0.024 * oneWayDelay + 0.11 * (oneWayDelay - 120) * (oneWayDelay > 120 ? 1 : 0) : oneWayDelay / 40;
    R -= delayFactor;
    
    // 丢包影响 (默认使用G.711编解码器的丢包特性)
    // Ie-eff = Ie + (95 - Ie) * Ppl / (Ppl + Bpl)
    // 对于G.711，Ie通常为0，Bpl通常为10
    const Ie = 0; // G.711的设备损伤因子
    const Bpl = 10; // G.711的丢包鲁棒性因子
    const packetLossPercent = packetLoss * 100;
    const lossImpairment = packetLossPercent > 0 ? Ie + (95 - Ie) * packetLossPercent / (packetLossPercent + Bpl) : 0;
    R -= lossImpairment;
    
    // 抖动影响 (基于jitter buffer分析的简化模型)
    // 抖动会导致额外的丢包或延迟
    const jitterImpairment = jitter > 30 ? 1 + (jitter - 30) / 10 : 0;
    R -= jitterImpairment;
    
    // 转换为MOS (1-5)
    let MOS = 1;
    if (R > 0) {
      if (R < 100) {
        MOS = 1 + 0.035 * R + 0.000007 * R * (R - 60) * (100 - R);
      } else {
        MOS = 4.5;
      }
    }
    
    // 限制范围
    MOS = Math.max(1, Math.min(4.5, MOS));
    
    return MOS;
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
      if (res.data.rtcp_report) {
        setRtcpReport(res.data.rtcp_report);
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
          <Tabs.TabPane tab="RTCP报告" key="rtcp_report">
            <div ref={rtcpReportRef}></div>
          </Tabs.TabPane>
          <Tabs.TabPane tab="本地RTCP分析" key="local_rtcp_analysis">
            <div ref={localRtcpAnalysisRef}></div>
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
