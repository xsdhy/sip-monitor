const {createProxyMiddleware} = require('http-proxy-middleware');

module.exports = function (app) {
    app.use(
        '/api',                              //访问端口已是api开头的请求都代理到8000
        createProxyMiddleware({
            target: 'http://127.0.0.1:9059',   //代理地址
            changeOrigin: true,
        })
    )
};