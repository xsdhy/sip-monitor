import SequenceDiagram from "../views/details/SequenceDiagram";
import {Button, message} from "antd";
import { useRef } from "react";

interface Prop {
    callID: string
    sessionID: string
}

export function SequenceModel(p: Prop) {
    const containerRef = useRef<HTMLDivElement>(null);

    const copyToClipboardAlternative = async () => {
        try {
            const mermaidContainer = document.querySelector('.mermaid');
            if (!mermaidContainer) {
                message.error('未找到序列图内容');
                return;
            }

            // 获取实际显示的 SVG 元素
            const svgElement = mermaidContainer.querySelector('svg');
            if (!svgElement) {
                message.error('未找到 SVG 元素');
                return;
            }

            // 获取实际尺寸
            const rect = svgElement.getBoundingClientRect();
            const width = rect.width;
            const height = rect.height;

            // 创建一个新的 svg 数据
            const serializer = new XMLSerializer();
            const svgData = serializer.serializeToString(svgElement);
            
            // 创建一个带有白色背景的完整 SVG 数据
            const fullSvgData = `
                <svg xmlns="http://www.w3.org/2000/svg" width="${width}" height="${height}">
                    <rect width="100%" height="100%" fill="white"/>
                    ${svgData}
                </svg>
            `;
            
            // 转换为 blob
            const blob = new Blob([fullSvgData], { type: 'image/svg+xml' });
            const url = URL.createObjectURL(blob);
            
            // 创建 Canvas
            const canvas = document.createElement('canvas');
            const ctx = canvas.getContext('2d');
            if (!ctx) {
                message.error('浏览器不支持 Canvas');
                URL.revokeObjectURL(url);
                return;
            }
            
            // 设置 canvas 尺寸
            const pixelRatio = window.devicePixelRatio || 1;
            canvas.width = width * pixelRatio;
            canvas.height = height * pixelRatio;
            
            try {
                // 在 Canvas 上绘制 SVG
                const img = new Image();
                await new Promise((resolve, reject) => {
                    img.onload = () => {
                        try {
                            // 填充白色背景
                            ctx.fillStyle = 'white';
                            ctx.fillRect(0, 0, canvas.width, canvas.height);
                            
                            // 考虑设备像素比
                            ctx.scale(pixelRatio, pixelRatio);
                            
                            // 绘制图像
                            ctx.drawImage(img, 0, 0, width, height);
                            
                            // 转换为图片并复制到剪贴板
                            canvas.toBlob(async (blob) => {
                                if (blob) {
                                    try {
                                        const clipboardItem = new ClipboardItem({ 'image/png': blob });
                                        await navigator.clipboard.write([clipboardItem]);
                                        message.success('已复制到剪贴板');
                                        resolve(true);
                                    } catch (err) {
                                        console.error('剪贴板写入失败:', err);
                                        message.error('复制到剪贴板失败');
                                        reject(err);
                                    }
                                } else {
                                    message.error('生成图片失败');
                                    reject(new Error('Failed to create blob'));
                                }
                            }, 'image/png', 1.0);
                        } catch (err) {
                            console.error('画布绘制失败:', err);
                            message.error('图片处理失败');
                            reject(err);
                        }
                    };
                    
                    img.onerror = (e) => {
                        console.error('图片加载失败:', e);
                        message.error('图片加载失败');
                        reject(e);
                    };
                    
                    img.src = url;
                });
            } finally {
                URL.revokeObjectURL(url);
            }
        } catch (error) {
            console.error('复制失败:', error);
            message.error('复制到剪贴板失败');
        }
    };

    return (
        <div ref={containerRef}>
              <div style={{ marginBottom: '10px' }}>
              <strong>CallID:</strong><span>{p.callID}</span><br/>
              <strong>SessionID:</strong><span>{p.sessionID}</span>
              </div>
            <SequenceDiagram callID={p.callID}/>
            <div style={{ marginTop: '10px' }}>
                <Button type="link">
                    <a target="_blank" href={'/call/details?sip_call_id=' + p.callID}>新页面中打开</a>
                </Button>
                <Button type="link" onClick={copyToClipboardAlternative}>
                    保存图片到剪切板
                </Button>
            </div>
        </div>
    );
}