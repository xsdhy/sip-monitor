import { CallRecordListDTO } from "@/@types/dto_list";
import { UserInfo,ResponseData, CallDetailsVO, CallRecordRaw, SIPRecordCall, CallStatVO } from "@/@types/entity";
import AppAxios from "@/utils/request";
import { Gateway } from "./gateway";


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
    return AppAxios.get<UserInfo[]>("/users", params);
  },

  /**
   * 获取用户详情
   */
  getUserDetail(id: string): Promise<ResponseData<UserInfo>> {
    return AppAxios.get<UserInfo>(`/users/${id}`);
  },

  /**
   * 创建用户
   */
  createUser(data: UserInfo): Promise<ResponseData<UserInfo>> {
    return AppAxios.post<UserInfo>("/users", data);
  },

  /**
   * 更新用户
   */
  updateUser(data: UserInfo): Promise<ResponseData<UserInfo>> {
    return AppAxios.put<UserInfo>(`/users/${data.id}`, data);
  },

  /**
   * 删除用户
   */
  deleteUser(id: string): Promise<ResponseData<UserInfo>> {
    return AppAxios.delete<UserInfo>(`/users/${id}`);
  },

  /**
   * 更新密码
   */
  updatePassword(data: { old_password: string; new_password: string }): Promise<ResponseData<UserInfo>> {
    return AppAxios.put<UserInfo>("/user/password", data);
  }



};

// 呼叫相关接口
export const callApi = {
  /**
   * 获取呼叫列表
   */
  getCallList(params?: CallRecordListDTO): Promise<ResponseData<SIPRecordCall[]>> {
    return AppAxios.get<SIPRecordCall[]>("/record/call", params);
  },
  
  /**
   * 获取呼叫详情
   */
  getCallDetail(id: string): Promise<ResponseData<CallDetailsVO>> {
    return AppAxios.get<CallDetailsVO>(`/record/details?sip_call_id=${id}`);
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

export const gatewayApi = {
  /**
   * 获取网关列表
   */
  getGatewayList(): Promise<ResponseData<Gateway[]>> {
    return AppAxios.get<Gateway[]>("/gateways");
  },

  /**
   * 获取网关详情
   */
  getGatewayDetail(id: string): Promise<ResponseData<Gateway>> {
    return AppAxios.get<Gateway>(`/gateways/${id}`);
  },

  /**
   * 创建网关
   */
  createGateway(data: Gateway): Promise<ResponseData<Gateway>> {
    return AppAxios.post<Gateway>("/gateways", data);
  },

  /**
   * 更新网关
   */
  updateGateway(id: string, data: Gateway): Promise<ResponseData<Gateway>> {
    return AppAxios.put<Gateway>(`/gateways/${id}`, data);
  },

  /**
   * 删除网关
   */
  deleteGateway(id: string): Promise<ResponseData<Gateway>> {
    return AppAxios.delete<Gateway>(`/gateways/${id}`);
  }
}

// 导出所有API
export default {
  user: userApi,
  call: callApi,
  gateway: gatewayApi
  // 可继续添加其他模块
};
