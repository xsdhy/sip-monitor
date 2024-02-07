import {Button, Form, Input} from 'antd';
import AppAxios from "../../utils/request";
import {customHistory} from "../../utils/history";



const Login = () => {
    const [form] = Form.useForm();

    const onFinish = (values:any) => {
        AppAxios.post("/login", values).then(res => {

            console.log(res, res.data.code)
            if (2000 === res.data.code) {
                window.localStorage.setItem("token", res.data.data)
                customHistory.push("/")
            }
        })
    };



    return (
        <div className="login-form">
            <h1 style={{textAlign:"center"}}>后台管理系统</h1>
            <Form
                labelCol={{span: 4,}}
                form={form}
                onFinish={onFinish}>
                <Form.Item
                    name="username"
                    label="用户名"
                    rules={[{required: true,message:"请输入用户名"}]}>
                    <Input/>
                </Form.Item>
                <Form.Item
                    name="password"
                    label="密码"
                    rules={[{required: true,message:"请输入密码"}]}>
                    <Input.Password/>
                </Form.Item>
                <Form.Item
                    name="code"
                    label="MFA"
                    rules={[{required: false,message:"请输入MFA"}]}>
                    <Input/>
                </Form.Item>
                <Form.Item
                    wrapperCol={{
                        offset: 4,
                        span: 16,
                    }}>
                    <Button type="primary" htmlType="submit">
                        登录
                    </Button>

                </Form.Item>
            </Form>
        </div>
    );
}
export default Login