// import store from '@/store'
import { performAppInitialization } from '@/api/plugin'
import store from '@/store'
import { Auth } from 'aws-amplify'
import { ALL_PROJECTS } from '@/gql/queries/projects.gql'

const app = {
  state: {
    projectid: null,
    project: {},
    projects: null
  },
  mutations: {
    SET_PROJECTS: (state, val) => {
      state.projects = val
    },
    SET_ACTIVE_PROJECT: (state, val) => {
      const oldProjectId = state.projectid
      const ttmpProjectId = val
      if (val) {
        if (ttmpProjectId !== oldProjectId) {
          state.projectid = ttmpProjectId
        }
      } else {
        state.project = { title: '' }
        state.projectid = null
      }

      /*
      store.getters.Auth.currentSession()
         .then(data => console.log(data))
         .catch(err => console.log(err))
     */

      /*
      //refreshes credentials using AWS.CognitoIdentity.getCredentialsForIdentity()
      AWS.config.credentials.refresh((error) => {
          if (error) {
               console.error(error);
          } else {
               // Instantiate aws sdk service objects now that the credentials have been updated.
               // example: var s3 = new AWS.S3();
               console.error(this)
           }
       })
       */
      /* Auth.currentSession()
         .then(data => console.log(data))
         .catch(err => console.log(err));
         */

      // console.error(store.getters.credentials)
      /* OPTION 1
       var user = store.getters.getCognitoUser
       Auth.updateUserAttributes(user, {
         'custom:project': 'patate'
       }).then((f) => {
         console.error('upated to')
         console.error(f)
       })
       */

      // OPTION 2
      if (state.projectid !== null) {
        if (oldProjectId !== state.projectid) {
          store.getters.Auth.currentAuthenticatedUser().then((user) => {
            console.error('current user is a')
            console.error(user)
            Auth.updateUserAttributes(user, {
              'custom:project': state.projectid
            }).then((uu) => {
              // After rupddatte make sure client wize we have proper attributes ...
              store.dispatch('REFRESH_COGNITO_USER')
            })
          })
        }
        try {
          const data = store.getters.apollo.readQuery({ query: ALL_PROJECTS })
          store.dispatch('setProjects', data.listAllProjects.items)
          var projectItem = data.listAllProjects.items.filter(project => project.id === state.projectid)
          state.project = projectItem[0]
        } catch (ezz) {
          state.project = { title: '' }
          if (ezz) {
            console.error(ezz)
            console.error('ignore')
          }
        }
      }
    }
  },
  getters: {
    appExample: () => {
      return {
        test: 'appExample fromPlugin'
      }
    },
    hasProjects: (state) => {
      if (!state.projects) {
        return true
      }

      return state.projects.length > 0
    },
    project: (state) => {
      return state.project
    },
    projectid: (state) => {
      if (!state.projectid) {
        // without cache ... ?
        /*
        store.getters.Auth.currentAuthenticatedUser({
          bypassCache: true
        }).then((curUser) => {
        */
        /*
        store.getters.Auth.currentAuthenticatedUser().then((curUser) => {
          console.error(curUser.attributes)
          resolve(curUser.attributes['custom:project'])
        })
        */
        if (store.getters.getCognitoUser && store.getters.getCognitoUser.attributes) {
          console.error(store.getters.getCognitoUser.attributes)
          state.projectid = store.getters.getCognitoUser.attributes['custom:project']
          store.dispatch('setProject', store.getters.getCognitoUser.attributes['custom:project'])
        }
      }
      // always run it ... in case we refresh using f5 and no cookies / data?
      // THIS CAUSE LOOPING ... store.dispatch('setProject', state.projectid)
      return state.projectid
    }
  },
  actions: {
    setProjects({ commit, state }, data) {
      console.error('set projects of')
      console.error(data)
      commit('SET_PROJECTS', data)
    },
    setProject({ commit, state }, data) {
      commit('SET_ACTIVE_PROJECT', data)
    },
    PerformPluginAppInitialization({ commit }, data) {
      // Lets validate if thet accountid match and if website has not been initialized yet ...
      console.error(this.state)
      data['accountid'] = this.state.app.accountid
      var self = this

      return new Promise((resolve, reject) => {
        // TODO: Add some security such as signed url / random token in the web UI ?
        performAppInitialization(data).then(response => {
          // This returns the list of plugins
          if (response && response.data) {
            // commit('SET_ACCOUNT_ID', null)
            commit('SET_LOADING', false)
            window.preventLoop = false // ?? Hack see in settings ....
            self.dispatch('GetSiteSettings').then(res => { // Required for firstTime loging notification in login page and other plugins infos?
              return resolve(true)
            })
          } else {
            return resolve(false)
          }
        }).catch((ex) => {
          console.error(ex)
          resolve(ex)
        })
      })
    }
  }
}

export default app
