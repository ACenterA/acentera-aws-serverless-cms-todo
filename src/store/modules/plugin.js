// If we would want a setup wizard on plugin installation we could force to use a specific
// router such as 'routerSetup' with only few routes available...
// import { routerSetup, routerValid, routerInvalid } from '@/router'
// import { routerValid, routerInvalid } from '@/router'
import { routerValid } from '@/router'

// import store from '@/store'
import { getPluginSettings } from '@/api/plugin'

const plugin = {
  state: {
  },
  mutations: {
  },
  getters: {
    pluginExample: () => {
      return {
        test: 'fromPlugin'
      }
    }
  },
  actions: {
    GetPluginRoutes({ commit, state }, obj) {
      console.error('GetPluginRoutes plugin from ehre CHILD')
      const info = {
        pluginName: 'serverless01',
        routes: routerValid
      }
      console.error('returning of')
      console.error(info)
      return info
      // return 'fff'
    },

    ActivatePlugins({ commit, state }, obj) {
      console.error('activate plugin from ehre CHILD')
      const info = {
        pluginName: 'serverless01',
        routes: routerValid
      }
      return info
      // return 'fff'
    },
    PluginAWSECLoad({ commit }, store) {
      return new Promise((resolve, reject) => {
        getPluginSettings().then(response => {
          var isReady = 0
          /*
          */
          isReady = 1
          resolve(isReady)
        }).catch((err) => {
          if (err) {
            console.error(err)
          }
          resolve(-1)
        })
      })
    },
    Ready(store) {
      // Core Application is ready.
      // Only do something here if needed, ie: the main plugin need to initialize a new custom dashboard maybe ...
      // This is the main Site Entry Point (or Plugin entry point)
      // ie: we can do things like get the current information we need for the application if needed or not ...
      // console.log('Application is ready...')

      // retreive the current route informations, and call the next(xxx) if we are the main plugin...
      const routeInfo = store.getters.GetRouteInfo
      // var loadAmplifyAuth = store.getters.Auth
      // console.error(loadAmplifyAuth)
      window.plugin_loaded++
      store.dispatch('PluginAWSECLoad').then((res) => {
        console.error('RECEIVED PluginAWSECLoad ' + res)
        // In here we could now load the routes we want based on if the app is configured or not..
        // ie: if (res.IsConfigured) { show all routes } else { show only route and menu X for wizard configuration }
        // add all routes to be loaded ...
        var routerRoutes = null
        /* if (res < 0) {
          routeInfo.next({ path: '/404' })
          // store.commit('NPROGRESS_END')
          return
        }*/

        if (res === 0) { // SETUP REQUIRED : ie: First time this plugin is being initiaized ...
          // routerRoutes = routerSetup
          console.error('SENDING ROUTER RES 0')
          routerRoutes = routerValid
        } else if (res === 1) { // Ok we have we were initialized...
          console.error('SENDING ROUTER RES 1')
          routerRoutes = routerValid
        } else if (res === 2) {
          console.error('SENDING ROUTER RES 2')
          // TODO ?
        } else if (res === -1) {
          // If its invalid: ?
          console.error('SENDING ROUTER INVALID?')
          // routerRoutes = routerInvalid //unless we modify api to send 200 but will have res -1?
          routerRoutes = routerValid
        }

        // Ok Add teh routes
        if (routerRoutes) {
          var l = routerRoutes.length
          for (var z = 0; z < l; z++) {
            window.asyncTestRouterMapTemp.push(routerRoutes[z])
          }
        }

        if (res < 0) {
          // routeInfo.next({ path: '/login' })
          // store.commit('NPROGRESS_END')
          return res
        }
        return res
      }).then((res) => {
        // We could limit the roles to only X if we want ...
        // store.dispatch('ActivatePlugins', store.getters.roles).then(function(r) {
        // Everything is ready to be displayed
        console.error('why 404? we just go to login actually')
        if (res < 0) {
          // error
          routeInfo.next({ path: '/login' })
          // routeInfo.next({ ...routeInfo.to, replace: true }) // hack ... addRoutes, set the replace: true so the navigation will not leave a history record
          return
        } else {
          console.error('RES IS.... ' + res)
          /*
          console.error('GOING TO')
          console.error(routeInfo.to)
          routeInfo.next({ ...routeInfo.to, replace: true }) // hack ... addRoutes, set the replace: true so the navigation will not leave a history record
          */
        }
        store.commit('NPROGRESS_END')
        return
        // routeInfo.push('/404')

        // window.app._router.push({ path: '/404' })
        // routeInfo.next({ path: '/404' }) // { ...routeInfo.to, replace: true }) // hack ... addRoutes, set the replace: true so the navigation will not leave a history record
      }).catch((err) => {
        // why??
        console.error(err.stack)
        routeInfo.next({ path: '/404' })
        return
      })
    },
    RouteChange(store) {
      // console.error('route change occured ')
      // console.error('NEW ROUTE:', window.app.$router)
      // console.error(store)
    }
  }
}

export default plugin
