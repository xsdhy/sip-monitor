import axios, { AxiosRequestConfig } from 'axios';
import { message } from 'antd';
import { ResponseData } from '@/@types/entity';



// 封装请求方法
class AppRequest {
  private instance;

  constructor() {
    this.instance = axios.create({
      baseURL: "/api",
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    // 请求拦截器
    this.instance.interceptors.request.use(
      (config) => {
        // 获取 token 并添加到请求头
        const token = localStorage.getItem('token');
        if (token && config.headers) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );

    // 响应拦截器
    this.instance.interceptors.response.use(
      (response) => {
        const res = response.data;

        // 根据自定义状态码进行处理
        if (res.code && (res.code !== 2000 && res.code !== 200)) {
          message.error(res.message || '请求失败');
          
          // token 失效处理
          if (res.code === 401) {
            localStorage.removeItem('token');
            window.location.href = '/login';
          }
          
          return Promise.reject(new Error(res.message || '请求失败'));
        }
        
        return res;
      },
      (error) => {
        const { response } = error;
        if (response && response.data) {
          message.error(response.data.message || '请求失败');
        } else {
          message.error('网络错误，请检查您的网络连接');
        }
        return Promise.reject(error);
      }
    );
  }

  // GET 请求
  get<T = any, R = ResponseData<T>>(url: string, params?: any, config?: AxiosRequestConfig): Promise<R> {
    return this.instance.get(url, { params, ...config });
  }

  // POST 请求
  post<T = any, R = ResponseData<T>>(url: string, data?: any, config?: AxiosRequestConfig): Promise<R> {
    return this.instance.post(url, data, config);
  }

  // PUT 请求
  put<T = any, R = ResponseData<T>>(url: string, data?: any, config?: AxiosRequestConfig): Promise<R> {
    return this.instance.put(url, data, config);
  }

  // DELETE 请求
  delete<T = any, R = ResponseData<T>>(url: string, config?: AxiosRequestConfig): Promise<R> {
    return this.instance.delete(url, config);
  }
}

const AppAxios = new AppRequest();
export default AppAxios;