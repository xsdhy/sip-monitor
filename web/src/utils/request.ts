import axios from 'axios';

import {notification} from 'antd';
import {customHistory} from './history'

const AppAxios = axios.create({
    baseURL: "/api",
})


const openNotificationWithIcon = (type:string, title:string, content:string) => {
    // @ts-ignore
    notification[type]({
        message: title,
        description: content,
    });
};

//添加拦截
AppAxios.interceptors.request.use(config => {
    let token = window.localStorage.getItem("token");
    if (token) {
        config.headers.token = token;
    }
    return config
}, error => {
    return Promise.reject(error);
})

AppAxios.interceptors.response.use(response => {
    if (4001 === response.data.code) {
        customHistory.push("/login");
        return Promise.reject("需要重新登录");
    } else if (response.data.code >= 5000) {
        openNotificationWithIcon("error", "操作失败", response.data.msg)
        return  Promise.reject("操作失败");
    }
    return Promise.resolve(response);
}, error => {
    return error;
})

export default AppAxios;