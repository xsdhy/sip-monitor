import dayjs from "dayjs";

export type CallRecordListDTO = {
    page?:number
    page_size?:number

    sip_call_id?: string
    node_ip?: string
    ua?: string

    begin_time?: string
    end_time?: string

    from_user?: string
    src_host?: string
    to_user?: string
    dst_host?: string

    date_picker?: [dayjs.Dayjs, dayjs.Dayjs]
}


export type CleanSipRecordDTO = {
    begin_time?: string
    end_time?: string
    method?: string
}