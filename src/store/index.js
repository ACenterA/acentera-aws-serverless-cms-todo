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

window.pluginStore = window.pluginStore || null

const initStore = () => {
  let StoreTmp = null
  const storeObj = {
    // no namepsacing namespaced: true,
    modules: {
      'serverless-cmsplugin': plugin,
      'serverless-cmspluginapollo': apollo,
      'serverless-cmspluginapp': app
    }
  }
  if (typeof window !== 'undefined') {
    console.error('test a')
    if (window.pluginStore) {
      return window.pluginStore
    }
    if (window && window.app) {
      StoreTmp = window.coreStore
    } else {
      if (window && window.app) {
        StoreTmp = window.coreStore || window.app.$store
      } else {
        StoreTmp = window.coreStore || null
      }
    }
    if (window.deepClone) {
      console.error('wlll create store bocj from inner (SUBPROJECT)')
      // const vObj = new Vuex.Store(storeObj)
      // console.error(StoreTmp)
      for (const v of Object.keys(storeObj.modules)) {
        console.error(`subproject will register module of ${v}`)
        console.error(storeObj.modules[v])
        // StoreTmp.registerModule(`serverless-cms${v}`, storeObj.modules[v])
        StoreTmp.registerModule(`${v}`, storeObj.modules[v])
      }
      /*
      const currentState = window.deepClone(StoreTmp.state)
      console.error('test store aa')
      const newStateTmp = window.deepClone(vObj.state)
      console.error('test store aa1')
      const newState = Object.assign(currentState, newStateTmp)
      console.error('test store aa2')
      console.error(vObj)
      const merged = { ...StoreTmp, ...vObj }

      for (const v of Object.keys(storeObj.modules)) {
        console.error(`subproject will register module of ${v}`)
        console.error(storeObj.modules[v])
        // StoreTmp.registerModule(`serverless-cms${v}`, storeObj.modules[v])
        StoreTmp.registerModule(`${v}`, storeObj.modules[v])
      }

      console.error('test store aa3')
      console.error(merged)

      console.error('test store aa')
      console.error('test store aa4')
      console.error(newState)

      // StoreTmp.replaceState(vObj, { module: true })
      StoreTmp.replaceState(newState, { module: true })
      StoreTmp = vObj
      */
      window.pluginStore = StoreTmp
      // window.store = StoreTmp
      return StoreTmp
    }
  }
  StoreTmp = new Vuex.Store(storeObj)
  window.pluginStore = StoreTmp
  return StoreTmp
}

const createStore = () => (initStore())
export default createStore()
/*
const initStore = () => {

  let store = window.app.$store || null
  // console.error('GOT STORE TEST OF')
  // console(store)
  return store
  || (store = new Vuex.Store({
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
*/
