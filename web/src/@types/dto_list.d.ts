import dayjs from "dayjs";

export type CallRecordListDTO = {
    page?:number
    page_size?:number

    sip_call_id?: string
    session_id?: string


    begin_time?: string
    end_time?: string

    from_user?: string
    to_user?: string
    src_host?: string
    dst_host?: string

    hangup_code?: string

    date_picker?: [dayjs.Dayjs, dayjs.Dayjs]
}


export type CleanSipRecordDTO = {
    begin_time?: string
    end_time?: string
    method?: string
}