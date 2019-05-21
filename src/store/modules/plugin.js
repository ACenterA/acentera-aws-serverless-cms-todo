// If we would want a setup wizard on plugin installation we could force to use a specific
// router such as 'routerSetup' with only few routes available...
// import { routerSetup, routerValid, routerInvalid } from '@/router'
import { routerValid, routerInvalid } from '@/router'

import store from '@/store'
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
    AddValidRoutes({ commit }, data) {
      return new Promise((resolve, reject) => {
        // Check if any projects existed in the past.. if not lets add the new routes
        // if (!store.getters.getAmiProjectsLst || store.getters.getAmiProjectsLst.length <= 0) {
        var routerRoutes = routerValid
        var l = routerRoutes.length
        for (var z = 0; z < l; z++) {
          window.asyncTestRouterMapTemp.push(routerRoutes[z])
        }
        store.dispatch('ActivatePlugins', store.getters.roles).then(function(r) {
          resolve(r)
        })
        // } else {
        //   resolve([])
        // }
      })
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
        // In here we could now load the routes we want based on if the app is configured or not..
        // ie: if (res.IsConfigured) { show all routes } else { show only route and menu X for wizard configuration }
        // add all routes to be loaded ...
        var routerRoutes = null
        if (res < 0) {
          routeInfo.next({ path: '/404' })
          // store.commit('NPROGRESS_END')
          return
        }

        if (res === 0) { // SETUP REQUIRED : ie: First time this plugin is being initiaized ...
          // routerRoutes = routerSetup
          routerRoutes = routerValid
        } else if (res === 1) { // Ok we have we were initialized...
          routerRoutes = routerValid
        } else if (res === 2) {
          // TODO ?
        } else if (res === -1) {
          // If its invalid: ?
          routerRoutes = routerInvalid
        }

        // Ok Add teh routes
        if (routerRoutes) {
          var l = routerRoutes.length
          for (var z = 0; z < l; z++) {
            window.asyncTestRouterMapTemp.push(routerRoutes[z])
          }
        }
        return res
      }).then((res) => {
        // We could limit the roles to only X if we want ...
        store.dispatch('ActivatePlugins', store.getters.roles).then(function(r) {
          // Everything is ready to be displayed
          if (res < 0) {
            // error
            routeInfo.next({ path: '/404' })
            return
          }
          routeInfo.next({ ...routeInfo.to, replace: true }) // hack ... addRoutes, set the replace: true so the navigation will not leave a history record
          return
        })
        // routeInfo.push('/404')

        // window.app._router.push({ path: '/404' })
        // routeInfo.next({ path: '/404' }) // { ...routeInfo.to, replace: true }) // hack ... addRoutes, set the replace: true so the navigation will not leave a history record
      }).catch((err) => {
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
