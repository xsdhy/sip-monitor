import { UserInfo,ResponseData, CallDetailsVO, CallRecordRaw, SIPRecordCall } from "@/@types/entity";
import AppAxios from "@/utils/request";


// 用户相关接口
export const userApi = {
  /**
   * 获取当前用户信息
   */
  getCurrentUser(): Promise<ResponseData<UserInfo>> {
    return AppAxios.get<UserInfo>("/user/current");
  },
  
  /**
   * 用户登录
   */
  login(data: { username: string; password: string }): Promise<ResponseData<{ token: string }>> {
    return AppAxios.post<{ token: string }>("/user/login", data);
  },
  

};

// 呼叫相关接口
export const callApi = {
  /**
   * 获取呼叫列表
   */
  getCallList(params?: { page: number; size: number }): Promise<ResponseData<SIPRecordCall[]>> {
    return AppAxios.get<SIPRecordCall[]>("/call/list", params);
  },
  
  /**
   * 获取呼叫详情
   */
  getCallDetail(id: string): Promise<ResponseData<CallDetailsVO>> {
    return AppAxios.get<CallDetailsVO>(`/call/detail/${id}`);
  },


  getCallRecordRaw(id: string): Promise<ResponseData<CallRecordRaw>> {
    return AppAxios.get<CallRecordRaw>(`/call/record/raw/${id}`);
  }
}



// 导出所有API
export default {
  user: userApi,
  call: callApi
  // 可继续添加其他模块
};
