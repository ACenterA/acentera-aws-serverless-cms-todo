import Vue from 'vue'
import router from './router'
import i18n from './lang'
import App from './App'
import FormWizard from 'vue-form-wizard'
import 'vue-form-wizard/dist/vue-form-wizard.min.css'

import vueNcform from '@ncform/ncform'
import Element from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import ncformStdComps from '@ncform/ncform-theme-elementui'

// import Home from './pages/Home'
import store from '@/store'

import 'assets/css/app.styl'

window.VueObj = window.VueObj || Vue

VueObj.use(FormWizard)
VueObj.use(Element)
VueObj.use(vueNcform, { extComponents: ncformStdComps, lang: 'fr-ca' })
// sort on values

/* eslint-disable */
window.sortArr = function(desc) {
  return function(a, b) {
    return desc ? ~~(a < b) : ~~(a > b)
  }
}

// sort on key values
window.keySort = function(key, desc) {
  return function(a, b) {
    return desc ? ~~(a[key] < b[key]) : ~~(a[key] > b[key])
  }
}

window.srt = function(on, descending) {
  on = on && on.constructor === Object ? on : {}
  return function(a,b){
    if (on.string || on.key) {
      a = on.key ? a[on.key] : a
      a = on.string ? String(a).toLowerCase() : a
      b = on.key ? b[on.key] : b
      b = on.string ? String(b).toLowerCase() : b
      // if key is not present, move to the end
      if (on.key && (!b || !a)) {
        return !a && !b ? 1 : !a ? 1 : -1
      }
    }
    return descending ? ~~(on.string ? b.localeCompare(a) : a < b)
                     : ~~(on.string ? a.localeCompare(b) : a > b)
  }
}
window.srtInt = function (on, descending) {
 on = on && on.constructor === Object ? on : {}
 return function(a,b){
   if (on.string || on.key) {
     a = on.key ? a[on.key] : a
     a = on.string ? String(a).toLowerCase() : a
     b = on.key ? b[on.key] : b
     b = on.string ? String(b).toLowerCase() : b;
     // if key is not present, move to the end
     if (on.key && (!b || !a)) {
      return !a && !b ? 1 : !a ? 1 : -1
     }
   }
   return descending ? ~~(on.string ? b.localeCompare(a) : a < b)
                     : ~~(on.string ? a.localeCompare(b) : a > b)
  };
}
/* eslint-enable */

if (window.asyncTestRouterMapTemp) {
  // As Plugin
  // console.error('will perform async loading of the routes...')
} else {
  // As Non Plugin
  /*
    window.asyncTestRouterMapTemp.push({
      path: '/permission',
      layout: 'Layout',
      redirect: '/permission/index',
      alwaysShow: true, // will always show the root menu
      meta: {
        title: 'permission',
        icon: 'lock',
        roles: ['admin', 'editor'] // you can set roles in root nav
      },
      children: [
        {
          path: 'page',
          component: Home,
          name: 'PagePermission',
          meta: {
            title: 'pagePermission',
            roles: ['admin'] // or you can only set roles in sub nav
          }
        }
      ]
    })
  */
  /* eslint-disable-next-line no-new */
  new Vue({
    el: '#app',
    router,
    store,
    i18n,
    render: h => h(App)
  })
}
