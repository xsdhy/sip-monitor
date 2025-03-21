import { useEffect, useState } from "react";
import { DatePicker, Card, Table, Button, Row, Col, Spin, message, Statistic } from 'antd';
import { Column } from '@ant-design/charts';
import { callApi } from "@/apis/api";
import { CallStatVO } from "@/@types/entity";
import dayjs, { Dayjs } from 'dayjs';
import { RangeValueType } from "@/@types/base.t";
import { ValueDate } from "@/@types/base.t";

const { RangePicker } = DatePicker;

const TimePresets: ValueDate<RangeValueType<Dayjs>>[] = [
    {label: '今天', value: [dayjs().startOf('day'), dayjs()]},
    {label: '昨天', value: [dayjs().subtract(1, 'day').startOf('day'), dayjs().startOf('day')]},
    {label: '上周', value: [dayjs().subtract(1, 'week').startOf('week').startOf('day'), dayjs().subtract(1, 'week').endOf('week').endOf('day')]},
    {label: '本月', value: [dayjs().startOf('month'), dayjs().endOf('month')]},
    {label: '上月', value: [dayjs().subtract(1, 'month').startOf('month'), dayjs().subtract(1, 'month').endOf('month')] },
    {label: '1天内', value: [dayjs().subtract(1, 'day').startOf('day'), dayjs()] },
    {label: '3天内', value: [dayjs().subtract(3, 'day').startOf('day'), dayjs()] },
    {label: '7天内', value: [dayjs().subtract(7, 'day').startOf('day'), dayjs()] },
    {label: '15天内', value: [dayjs().subtract(15, 'day').startOf('day'), dayjs()] },
    {label: '1个月内', value: [dayjs().subtract(1, 'month').startOf('day'), dayjs()] },
    {label: '3个月内', value: [dayjs().subtract(3, 'month').startOf('month'), dayjs()] },
    {label: '6个月内', value: [dayjs().subtract(6, 'month').startOf('month'), dayjs()] },
    {label: '1年内', value: [dayjs().subtract(1, 'year').startOf('day'), dayjs()] },
]

function Home() {
    const [loading, setLoading] = useState(false);
    const [callStats, setCallStats] = useState<CallStatVO[]>([]);
    const [dateRange, setDateRange] = useState<[dayjs.Dayjs, dayjs.Dayjs]>([
        dayjs().startOf('day'),
        dayjs()
    ]);

    const fetchCallStats = async () => {
        try {
            setLoading(true);
            const params = {
                begin_time: dateRange[0].format('YYYY-MM-DDTHH:mm:ssZ'),
                end_time: dateRange[1].format('YYYY-MM-DDTHH:mm:ssZ')
            };
            
            const response = await callApi.getCallStat(params);
            setCallStats(response.data || []);
        } catch (error) {
            message.error('获取通话统计数据失败');
            console.error(error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchCallStats();
    }, []);

    const handleDateChange = (dates: any) => {
        if (dates && dates.length === 2) {
            setDateRange([dates[0], dates[1]]);
        }
    };

    const handleSearch = () => {
        fetchCallStats();
    };

    const statisticCard = () => {
        const getValueStyle = (answered: number, total: number) => {
            const rate = (answered / total) * 100;
            let color = '#cf1322'; // 默认红色 (<40%)
            
            if (rate >= 70) {
                color = '#3f8600'; // 绿色 (>=70%)
            } else if (rate >= 40) {
                color = '#096dd9'; // 蓝色 (40%-70%)
            }
            
            return {
                fontSize: 30,
                color,
            };
        };

        const statisticData = callStats.map(stat => (
            <Col key={stat.ip}>
                <Statistic 
                    title={stat.gateway ? `${stat.gateway} (${stat.ip})` : stat.ip}
                    value={(stat.answered / stat.total) * 100}
                    precision={2}
                    suffix="%"
                    valueStyle={getValueStyle(stat.answered, stat.total)}
                />
            </Col>
        ));

        return (
            <Card title="通话接通率" style={{ marginBottom: 16 }}>
                <Row gutter={16}>
                    {statisticData}
                </Row>
            </Card>
        );
    };

    // 转换数据为图表格式
    const transformDataForCharts = () => {
        const chartData: any[] = [];
        
        callStats.forEach(stat => {
            chartData.push({ ip: stat.ip, type: '总通话量', value: stat.total });
            chartData.push({ ip: stat.ip, type: '已接通', value: stat.answered });
        });
        
        return chartData;
    };
    
    // 柱状图配置
    const config = {
        data: transformDataForCharts(),
        isGroup: true,
        xField: 'ip',
        yField: 'value',
        seriesField: 'type',
        
        // 删除可能导致问题的label配置
        // 设置图例位置和其他样式
        legend: {
            position: 'top'
        },
        
        // 调整柱子样式
        columnStyle: {
            radius: [4, 4, 0, 0],
        },
        
        // 设置颜色
        color: ['#4377FE', '#0BA25F', '#F5B60D', '#E96A59', '#9861E5', '#4DCCCC', '#F2637B', '#B4BDFF'],
        
        // X轴设置
        xAxis: {
            label: {
                autoHide: true,
                autoRotate: false,
            },
        },
    };

    // 表格列定义
    const columns = [
        {
            title: 'IP',
            dataIndex: 'ip',
            key: 'ip',
            render: (text: string, record: CallStatVO) => {
                return `${record.gateway ? `${record.gateway} (${text})` : text}`;
            },
        },
        {
            title: '总通话量',
            dataIndex: 'total',
            key: 'total',
            sorter: (a: CallStatVO, b: CallStatVO) => a.total - b.total,
        },
        {
            title: '已接通',
            dataIndex: 'answered',
            key: 'answered',
            sorter: (a: CallStatVO, b: CallStatVO) => a.answered - b.answered,
        },
        {  
            title:'接通率',
            key: 'answered_rate',
            render: (record: CallStatVO) => {
                return `${((record.answered / record.total) * 100).toFixed(2)}%`;
            },
        },
        {
            title: '0xx状态码',
            dataIndex: 'hangup_code_0_count',
            key: 'hangup_code_0_count',
        },
        // {
        //     title: '1xx状态码',
        //     dataIndex: 'hangup_code_1xx_count',
        //     key: 'hangup_code_1xx_count',
        // },
        {
            title: '2xx状态码',
            dataIndex: 'hangup_code_2xx_count',
            key: 'hangup_code_2xx_count',
        },
        // {
        //     title: '3xx状态码',
        //     dataIndex: 'hangup_code_3xx_count',
        //     key: 'hangup_code_3xx_count',
        // },
        {
            title: '4xx状态码',
            dataIndex: 'hangup_code_4xx_count',
            key: 'hangup_code_4xx_count',
        },
        {
            title: '5xx状态码',
            dataIndex: 'hangup_code_5xx_count',
            key: 'hangup_code_5xx_count',
        },
    ];

    return (
        <div>
            <h3>工作台</h3>
            <Card style={{ marginBottom: 16 }}>
                <Row gutter={16} align="middle">
                    <Col>
                        <RangePicker
                            showTime
                            value={dateRange}
                            onChange={handleDateChange}
                            style={{ marginRight: 16 }}
                            presets={TimePresets}
                        />
                    </Col>
                    <Col>
                        <Button type="primary" onClick={handleSearch}>
                            查询
                        </Button>
                    </Col>
                </Row>
            </Card>

            <Spin spinning={loading}>
                {statisticCard()}

                <Card title="通话状态码统计" style={{ marginBottom: 16 }}>
                    <div style={{ height: 400 }}>
                        <Column {...config} />
                    </div>
                </Card>
                

                <Card title="通话数据表格">
                    <Table
                        dataSource={callStats}
                        columns={columns}
                        rowKey="ip"
                        pagination={false}
                        scroll={{ x: 'max-content' }}
                    />
                </Card>
            </Spin>
        </div>
    );
}

export default Home;