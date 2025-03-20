import {CallRecordEntity} from '../../@types/entity'
import dayjs from "dayjs";

export function getProtocolName(num: number):string {
    if (num === 6) {return 'TCP'}
    if (num === 17) {return 'UDP'}
    if (num === 22) {return 'TLS'}
    if (num === 50) {return 'ESP'}
    return 'Unknown'
}

export function stringToColor(str: string) {
    let hash = 0
    for (let i = 0; i < str.length; i++) {
        hash = str.charCodeAt(i) + ((hash << 5) - hash)
    }

    let color = '#'

    for (let i = 0; i < 3; i++) {
        const value = (hash >> (i * 8)) & 0xff
        const v16 = '00' + value.toString(16)
        color += v16.substring(v16.length - 2)
    }
    return color
}

const regEx = /\d+/
export function isRequest(method: string) {
    return !regEx.test(method)
}



export function createSeqHtml(seq: CallRecordEntity[]): string {
    const res: string[] = [
        'sequenceDiagram'
    ]

    seq.forEach((item, index) => {
        let dis = 0
        if (index !== 0) {
            dis = seq[index].timestamp_micro/1000000 - seq[index - 1].timestamp_micro/1000000
        }
        item.dst_addr=item.dst_addr.replace(":","_")
        item.src_addr=item.src_addr.replace(":","_")
        const arrow = isRequest(item.method) ? '->>' : '-->>'
        let messageText = `${item.method} `
        
        if (item.method === 'INVITE') {
            messageText += `${item.from_user} -> ${item.to_user} `
        }
        messageText += `${item.response_desc} ${dis.toFixed(2)}s`

        res.push(`${item.src_addr}${arrow}${item.dst_addr}: ${messageText}`)
        
        // 添加时间戳注释
        if (index === 0 || dis > 1) {
            res.push(`Note over ${item.src_addr}: ${dayjs(item.timestamp_micro/1000).format('YYYY-MM-DD HH:mm:ss')}`)
        }
    })

    return res.join('\n')
}
