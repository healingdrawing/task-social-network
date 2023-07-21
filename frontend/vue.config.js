const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 3000,
    // headers: {
    //   'Cache-Control': 'no-store',
    // },
  }
})
/* no-store was added to dev more clear */