export interface CallRecordDetailsVO {
    code: number
    data: CallDetailsVO
    msg: string
    time: string
}

export interface CallRecordRawVO {
    code: number
    data: CallRecordRaw
    msg: string
    time: string
}
export interface CallRecordRaw {
   id: number
   create_time: string
   raw: string
}


export interface CallRecordBaseVO {
    code: number
    data: T
    meta?: MetaVO
    msg: string
    time: string
}

export interface MetaVO {
    page: number
    page_size: number
    total: number
}

export interface CallRecordEntity {
    id: number
    sip_call_id: string
    method: string
    response_desc: string

    to_user: string
    from_user: string

    src_addr: string
    dst_addr: string

    create_time: string
    timestamp_micro: number
}


export interface CallDetailsVO {
    records: CallRecordEntity[]
    relevants: CallRecordEntity[]
}

export interface SIPRecordCall {
    id: string;

    node_ip: string;
    sip_call_id: string;
    session_id: string;

    to_user: string;
    from_user: string;

    user_agent: string;


    src_addr: string;
    dst_addr: string;

    create_time: string; // Assuming time.Time is serialized as a string
    ringing_time: string; // Assuming time.Time is serialized as a string
    answer_time: string; // Assuming time.Time is serialized as a string
    end_time: string; // Assuming time.Time is serialized as a string

    call_duration: number;
    ringing_duration: number;
    talk_duration: number;

    hangup_code: number;
    hangup_cause: string;
}


export interface SIPRecordRegister {
    id: string;

    node_ip: string;

    create_time: string; // Assuming time.Time is serialized as a string

    sip_call_id: string;

    from_user: string;

    user_agent: string;

    register_times: number;
    failures_times: number;
    successes_times: number;

    src_host: string;
    src_port: number;
    src_addr: string;
    src_country_name: string;
    src_city_name: string;
}


export interface SystemDBStatsVO {
    name: string;

    capped: boolean;
    count: number;
    index_count: number;
    avg_obj_size: number;
    free_storage_size: number;
    size: number;
    storage_size: number;
    total_index_size: number;
    total_size: number;


}