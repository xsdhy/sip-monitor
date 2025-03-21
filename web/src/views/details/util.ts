import {CallRecordEntity} from '../../@types/entity'
import dayjs from "dayjs";



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
