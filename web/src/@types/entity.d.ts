export interface CallRecordDetailsVO {
    code: number
    data: CallRecordEntity[]
    meta: MetaVO
    msg: string
    time: string
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

export interface SipCallFlowDiagramData {
    columns: string[]
    messages: CallRecordEntity[]
}


export interface CallRecordEntity {
    id: number
    uuid: string
    call_id: string
    create_time: string

    method: string
    src: string
    dst: string

    timestamp: number
    body: string
}


export interface SIPRecordCall {
    id: string;

    node_ip: string;
    call_id: string;

    to_user: string;
    from_user: string;

    user_agent: string;

    src_host: string;
    src_port: number;
    src_addr: string;
    src_country_name: string;
    src_city_name: string;

    dst_host: string;
    dst_port: number;
    dst_addr: string;
    dst_country_name: string;
    dst_city_name: string;

    create_time: string; // Assuming time.Time is serialized as a string
    ringing_time: string; // Assuming time.Time is serialized as a string
    answer_time: string; // Assuming time.Time is serialized as a string
    end_time: string; // Assuming time.Time is serialized as a string

    call_duration: number;
    ringing_duration: number;
    talk_duration: number;
}


export interface SIPRecordRegister {
    id: string;

    node_ip: string;

    create_time: string; // Assuming time.Time is serialized as a string

    call_id: string;

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