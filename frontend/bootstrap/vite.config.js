import pluginPurgeCss from "@mojojoejo/vite-plugin-purgecss";

const path = require('path')
const debug = true

function pluginList() {
    if (!debug) {
        return [
            pluginPurgeCss({
                content: ["../../templates/**/*.html"],
                variables: true,
            }),
        ]
    }

    return []
}

export default {
    root: path.resolve(__dirname, 'src'),

    plugins: pluginList(),

    build: {
        outDir: '../dist',

        rollupOptions: {
            output: {
                //название главного файла javascript
                entryFileNames: 'bootstrap.js',
                //название остальных файлов. Если не прописать, то к имени будет дописывать случайный id
                assetFileNames: 'bootstrap.[ext]',
            },
        },
    },
    server: {
        port: 8080,
        hot: true
    }
}