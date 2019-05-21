import Vue from 'vue'
import Vuex from 'vuex'
import plugin from './modules/plugin'

import apollo from './modules/apollo'
import app from './modules/app'

// import s3 from './modules/s3'
// import ami from './modules/ami'
// import iam from './modules/iam'
// import cloudformation from './modules/cloudformation'
Vue.use(Vuex)

let store = window.app.$store

const initStore = () => {
  return store || (store = new Vuex.Store({
    modules: {
    },
    actions: {
    }
  }))
}

if (!initStore().state['serverless-cmsplugin']) {
  initStore().registerModule('serverless-cmsplugin', plugin)
  initStore().registerModule('serverless-cmspluginapollo', apollo)
  initStore().registerModule('serverless-cmspluginapp', app)
}

export default initStore()
