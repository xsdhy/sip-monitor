import {Modal} from "antd";
import {SequenceModel} from "../components/SequenceModel";
import dayjs from "dayjs";


export function SizeConversion(limit: number): string {
    var size = "";
    if (limit < 0.1 * 1024) {                            //小于0.1KB，则转化成B
        size = limit.toFixed(2) + "B"
    } else if (limit < 0.1 * 1024 * 1024) {            //小于0.1MB，则转化成KB
        size = (limit / 1024).toFixed(2) + "KB"
    } else if (limit < 0.1 * 1024 * 1024 * 1024) {        //小于0.1GB，则转化成MB
        size = (limit / (1024 * 1024)).toFixed(2) + "MB"
    } else {                                            //其他转化成GB
        size = (limit / (1024 * 1024 * 1024)).toFixed(2) + "GB"
    }

    var sizeStr = size + "";                        //转成字符串
    var index = sizeStr.indexOf(".");                    //获取小数点处的索引
    var dou = sizeStr.substr(index + 1, 2)            //获取小数点后两位的值
    if (dou === "00") {                                //判断后两位是否为00，如果是则删除00
        return sizeStr.substring(0, index) + sizeStr.substr(index + 3, 2)
    }
    return size;
}


export function ShowIPText(country: string | undefined, city: string | undefined): string {
    country = country || "";
    city = city || "";

    if (country === "局域网" || city === "局域网") {
        return "局域网";
    } else if (country === "中国" && city === "") {
        return country;
    } else if (country === "中国" || city !== "") {
        return city;
    } else if (country !== "中国" && city !== "") {
        return `${country}-${city}`;
    } else {
        return country !== "中国" ? country : "其他";
    }
}

export function FormatToDateTime(t: string|undefined):string {
    if (t===undefined){
        return ""
    }
    return dayjs(t).format('YYYY-MM-DD HH:mm:ss')
}

export function OpenSeqModel(callId: string, sessionId: string) {
    Modal.info({
        title: "呼叫信令",
        icon: null,
        content: (<SequenceModel callID={callId} sessionID={sessionId}/>),
        width: 800,
        closable: true,
        footer: null,
    })
}
