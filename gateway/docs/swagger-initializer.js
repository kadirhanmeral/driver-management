window.onload = function () {
    const ui = SwaggerUIBundle({
        url: "/swagger/doc.json",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
            SwaggerUIBundle.presets.apis,
            SwaggerUIStandalonePreset
        ],
        plugins: [
            SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout",
        requestInterceptor: (req) => {
            if (req.headers.Authorization && !req.headers.Authorization.startsWith('Bearer ')) {
                req.headers.Authorization = 'Bearer ' + req.headers.Authorization;
            }
            return req;
        }
    });
    window.ui = ui;
};
