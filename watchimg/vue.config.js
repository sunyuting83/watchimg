const cdn = {
  css: ['//lf6-cdn-tos.bytecdntp.com/cdn/expire-1-M/bulma/0.9.3/css/bulma.min.css'],
  js: [
    '//lf6-cdn-tos.bytecdntp.com/cdn/expire-1-M/vue/3.2.13/vue.global.min.js'
  ]
}
 
module.exports = {
  publicPath: '/',
  productionSourceMap: process.env.NODE_ENV !== 'production',
  devServer: {
    port: 8077
  },
  chainWebpack: config => {
    config.resolve.symlinks(true)
  },
  configureWebpack: {
    externals: {
      'vue': 'Vue'
    }
  },
  pages: {
    index: {
      entry: 'src/main.js',
      template: 'public/index.html',
      filename: 'index.html',
      title: 'admin',
      chunks: ['chunk-vendors', 'chunk-common', 'index'],
      cdn: cdn
    }
  }
}