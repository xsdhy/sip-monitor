import {CallRecordEntity} from '../../@types/entity'


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


//https://sequence.davidje13.com/library.htm
export function createSeqHtml(seq: CallRecordEntity[]):string {
    const res: string[] = [`autolabel "[<inc>] <label>"`]

    seq.forEach((item, index) => {
        let dis = 0
        if (index !== 0) {
            dis = seq[index].timestamp_micro/1000000 -seq[index - 1].timestamp_micro/1000000
        }
        //箭头:若是请求则->，否则-->
        const arrowhead = `-${isRequest(item.sip_method) ? '' : '-'}>`
        //内容
        const labelContent = `**${item.sip_method}** ${item.response_desc} ${dis.toFixed(1)}s ${stringToColor(item.cseq_number + '')}`

        res.push(`${item.src_addr}${arrowhead}${item.dst_addr}: ${labelContent}`)
    })

    res.push('terminators box')
    return res.join('\n')
}
