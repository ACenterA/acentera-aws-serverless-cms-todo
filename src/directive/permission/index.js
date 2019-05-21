import permission from './permission'

if (!window['permission']) {
  const install = function(Vue) {
    Vue.directive('permission', permission)
  }

  if (window.Vue) {
    if (!window['permission']) {
      window['permission'] = permission
      Vue.use(install); // eslint-disable-line
    }
  }

  permission.install = install
}
export default permission
