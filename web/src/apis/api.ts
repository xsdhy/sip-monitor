import { UserInfo,ResponseData, CallDetailsVO, CallRecordRaw, SIPRecordCall, CallStatVO } from "@/@types/entity";
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
  login(data: { username: string; password: string }): Promise<ResponseData<string>> {
    return AppAxios.post<string>("/login", data);
  },

  /**
   * 获取用户列表
   */
  getUserList(params?: { page: number; page_size: number }): Promise<ResponseData<UserInfo[]>> {
    return AppAxios.get<UserInfo[]>("/user/list", params);
  },

  /**
   * 获取用户详情
   */
  getUserDetail(id: string): Promise<ResponseData<UserInfo>> {
    return AppAxios.get<UserInfo>(`/user/detail/${id}`);
  },

  /**
   * 创建用户
   */
  createUser(data: UserInfo): Promise<ResponseData<UserInfo>> {
    return AppAxios.post<UserInfo>("/user/create", data);
  },

  /**
   * 更新用户
   */
  updateUser(data: UserInfo): Promise<ResponseData<UserInfo>> {
    return AppAxios.put<UserInfo>("/user/update", data);
  },

  /**
   * 删除用户
   */
  deleteUser(id: string): Promise<ResponseData<UserInfo>> {
    return AppAxios.delete<UserInfo>(`/user/delete/${id}`);
  },

};

// 呼叫相关接口
export const callApi = {
  /**
   * 获取呼叫列表
   */
  getCallList(params?: { page: number; page_size: number }): Promise<ResponseData<SIPRecordCall[]>> {
    return AppAxios.get<SIPRecordCall[]>("/record/call", params);
  },
  
  /**
   * 获取呼叫详情
   */
  getCallDetail(id: string): Promise<ResponseData<CallDetailsVO>> {
    return AppAxios.get<CallDetailsVO>(`/record/detail/${id}`);
  },


  getCallRecordRaw(id: string): Promise<ResponseData<CallRecordRaw>> {
    return AppAxios.get<CallRecordRaw>(`/record/raw/${id}`);
  },

  /**
   * 获取通话统计数据
   */
  getCallStat(params: { begin_time?: string; end_time?: string }): Promise<ResponseData<CallStatVO[]>> {
    return AppAxios.post<CallStatVO[]>("/stat/call", params);
  }
}



// 导出所有API
export default {
  user: userApi,
  call: callApi
  // 可继续添加其他模块
};
