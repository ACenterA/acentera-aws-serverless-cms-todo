// If we would want a setup wizard on plugin installation we could force to use a specific
// router such as 'routerSetup' with only few routes available...
// import { routerSetup, routerValid, routerInvalid } from '@/router'
import { routerValid, routerInvalid } from '@/router'

import store from '@/store'
import { getPluginSettings } from '@/api/plugin'

var plugName = `${PLUGIN_NAME}`
/* eslint-disable-next-line */
// var adminPluginNameMenu = plugName + ':administration'

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
    GetPluginRoutes({ commit }, data) {
      console.error('[PLUGIN:GetPluginRoutes] REceived event')
      console.error(data)
      if (data === plugName) {
        const routerRoutes = routerValid
        const Obj = {
          plugin: plugName,
          routes: routerRoutes
        }
        return Obj
      } else {
        return {}
      }
    },
    PluginLoaded() {
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
    Ready(resolvedRoute) {
      // Core Application is ready.
      // Only do something here if needed, ie: the main plugin need to initialize a new custom dashboard maybe ...
      // This is the main Site Entry Point (or Plugin entry point)
      // ie: we can do things like get the current information we need for the application if needed or not ...
      // console.log('Application is ready...')

      // retreive the current route informations, and call the next(xxx) if we are the main plugin...
      const routeInfo = store.getters.GetRouteInfo
      // var loadAmplifyAuth = store.getters.Auth
      // console.error(loadAmplifyAuth)
      store.dispatch('PluginLoaded').then(() => {
        store.dispatch('PluginAWSECLoad').then((res) => {
          // In here we could now load the routes we want based on if the app is configured or not..
          // ie: if (res.IsConfigured) { show all routes } else { show only route and menu X for wizard configuration }
          // add all routes to be loaded ...
          if (res < 0) {
            console.error('IN HERE WILL RETURN 404??')
            routeInfo.next({ path: '/login' })
            // store.commit('NPROGRESS_END')
            return
          }
          return res
        }).then((res) => {
          // We could limit the roles to only X if we want ...
          // Everything is ready to be displayed
          console.error('OK READY COMPLETED')
          routeInfo.next({ ...routeInfo.to, replace: true })
          // hack ... addRoutes, set the replace: true so the navigation will not leave a history record
          // return
        }).catch((err) => {
          console.error(err.stack)
          routeInfo.next({ path: `/404?plugin=${plugName}&error=loading` })
          return
        })
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
