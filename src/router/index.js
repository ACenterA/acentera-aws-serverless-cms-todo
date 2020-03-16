import Vue from 'vue'
import Router from 'vue-router'

import DashboardIndex from '@/views/dashboard/Index'
import Bootstrap from '@/views/bootstrap/index'

import DefaultLayout from '@/views/layout/DefaultLayout'
import AdminLayout from '@/views/admin/Layout'
import AdminDashboard from '@/views/admin/Dashboard'
import AdminProjects from '@/views/admin/Projects'
import AdminProjectsCreate from '@/views/admin/Projects/Create'

import AdminUsers from '@/views/admin/Users'
import AdminUsersCreate from '@/views/admin/Users/Create'

import AdminSettingsLayout from '@/views/admin/settings/Layout'
import AdminSettingsIndex from '@/views/admin/settings/Index'
import AdminSettingsExample from '@/views/admin/settings/Example'

import ContextSelector from '@/components/LeftMenu/ContextSelector'

// import PostsLayout from '@/views/posts/Layout'
// import PostsIndex from '@/views/posts/Index'
// import PostsCreate from '@/views/posts/Create'

/*
import ClusterLayout from '@/views/cluster/Layout'
import ClusterDashboard from '@/views/cluster/Dashboard'
import ClusterSettingsDanger from '@/views/cluster/settings/Danger'
import ClusterSettingsAdmin from '@/views/cluster/settings/Admin'
import ClusterDiscovery from '@/views/cluster/Discovery'
import ClusterNodes from '@/views/cluster/Nodes'
*/

Vue.use(Router)

/* eslint-disable-next-line */
var plugName = `${PLUGIN_NAME}`

/* eslint-disable-next-line */
var adminPluginNameMenu = plugName + ':administration'
// TODO: We could use this router if we want to limit the application into a setup phase
/*
export const routerSetup = [
  {
    path: '/',
    component: HomeTest
  },
  {
    path: '/dashboard',
    // component: Layout,
    layout: 'Layout',
    replacePrecedence: 100,
    replacePath: '/dashboard',
    hidden: true,
    children: [
      {
        path: 'setup',
        component: HomeTest,
        name: 'HomeTestSetup',
        meta: { title: 'HomeTestSetup Setup', icon: 'tab', hidden: true }
      }
    ]
  }
]
*/

export const routerValid = [
  {
    path: '/',
    component: DashboardIndex
  },
  {
    path: '/dashboard/index',
    // component: Layout,
    layout: 'Layout',
    replacePrecedence: 99, // Lower than the Setup Route
    replacePath: '/dashboard',
    hidden: true,
    children: [
      {
        path: '',
        component: DashboardIndex,
        name: 'DashboardIndex',
        meta: {
          title: 'DashboardIndex', icon: 'tab', hidden: true
        }
      }
    ]
  },
  {
    path: '/bootstrap/index',
    // component: Layout,
    // layout: 'Custom',
    layout: 'Custom',
    component: Bootstrap,
    replacePrecedence: 99, // Lower than the Setup Route
    replacePath: '/bootstrap', // this does not work we hardcoded to /bootstrap/index redirect
    hidden: true,
    children: [
      {
        path: '',
        component: Bootstrap,
        name: 'BootstrapPlugin',
        meta: {
          title: 'Bootstrap', icon: 'tab', hidden: true
        }
      }
    ]
  },
  {
    path: '/invalid/component/leftmenu/contextselector',
    // unused layout: 'Layout',
    component: ContextSelector,
    hidden: false,
    iscomponent: true,
    children: []
  },
  {
    path: '/posts',
    // component: Layout,
    layout: 'Layout',
    redirect: '/posts/index',
    alwaysShow: true, // will always show the root menu
    meta: {
      title: 'Posts',
      icon: 'lock',
      roles: ['admin', 'editor'] // you can set roles in root nav
    },
    children: [
      {
        path: 'index',
        component: () => import('@/views/posts/Index'),
        name: 'PostIndex',
        meta: {
          title: 'List',
          roles: ['admin', 'editor'] // or you can only set roles in sub nav
        }
      },
      {
        path: 'create',
        component: () => import('@/views/posts/Create'),
        name: 'PostCreate',
        meta: {
          title: 'Create'
          // if do not set roles, means: this page does not require permission
        }
      }
    ]
  },
  {
    path: '/administration',
    id: 3,
    layout: 'Layout',
    hidden: false,
    meta: {
      roles: ['admin']
    },
    children: [
      {
        id: 8,
        path: '',
        component: AdminLayout,
        name: 'AdministrationIndex',
        redirect: '/administration/index',
        meta: {
          title: 'Administration', icon: 'tab', hidden: false,
          roles: ['admin'] // you can set roles in root nav
        }
      }
    ]
  },
  {
    path: '/administration',
    id: 5,
    layout: 'Layout',
    plugin: plugName,
    innerMenu: adminPluginNameMenu,
    hidden: false,
    flatChildrens: true,
    meta: {
      roles: ['admin']
    },
    children: [
      {
        path: '',
        id: 9,
        name: 'adminRoot',
        layout: 'Layout',
        component: AdminLayout,
        hidden: false,
        // alwaysShow: false,
        flatChildrens: true,
        meta: {
          title: 'Dashboard',
          hideMainMenu: true,
          showMenu: adminPluginNameMenu,
          icon: 'lock',
          roles: ['admin'] // you can set roles in root nav
        },
        children: [
          {
            path: 'index',
            name: 'admin',
            layout: 'Layout',
            component: AdminDashboard,
            hidden: false,
            meta: {
              title: 'Dashboard',
              hideMainMenu: true,
              showMenu: adminPluginNameMenu,
              // example: class: 'serverless-cms-container',
              icon: 'lock',
              roles: ['admin'] // you can set roles in root nav
            }
          },
          {
            path: 'projects',
            name: 'AdminProjects',
            layout: 'Layout',
            component: DefaultLayout,
            hidden: false,
            meta: {
              title: 'Projects',
              hideMainMenu: true,
              showMenu: adminPluginNameMenu,
              // example: class: 'serverless-cms-container',
              icon: 'lock',
              roles: ['admin'] // you can set roles in root nav
            },
            children: [
              {
                path: 'list',
                name: 'AdminProjectsList',
                layout: 'Layout',
                component: AdminProjects,
                hidden: false,
                meta: {
                  title: 'List',
                  hideMainMenu: true,
                  showMenu: adminPluginNameMenu,
                  icon: 'lock',
                  roles: ['admin'] // you can set roles in root nav
                }
              },
              {
                path: 'create',
                name: 'AdminProjectsCreate',
                layout: 'Layout',
                component: AdminProjectsCreate,
                hidden: false,
                meta: {
                  title: 'Create',
                  hideMainMenu: true,
                  showMenu: adminPluginNameMenu,
                  icon: 'lock',
                  roles: ['admin'] // you can set roles in root nav
                }
              }
            ]
          },
          {
            path: 'users',
            name: 'AdminUsers',
            layout: 'Layout',
            component: DefaultLayout,
            hidden: false,
            meta: {
              title: 'Users',
              hideMainMenu: true,
              showMenu: adminPluginNameMenu,
              // example: class: 'serverless-cms-container',
              icon: 'lock',
              roles: ['admin'] // you can set roles in root nav
            },
            children: [
              {
                path: 'list',
                name: 'AdminUsersList',
                layout: 'Layout',
                component: AdminUsers,
                hidden: false,
                meta: {
                  title: 'List',
                  hideMainMenu: true,
                  showMenu: adminPluginNameMenu,
                  icon: 'lock',
                  roles: ['admin'] // you can set roles in root nav
                }
              },
              {
                path: 'create',
                name: 'AdminUsersCreate',
                layout: 'Layout',
                component: AdminUsersCreate,
                hidden: false,
                meta: {
                  title: 'Create',
                  hideMainMenu: true,
                  showMenu: adminPluginNameMenu,
                  icon: 'lock',
                  roles: ['admin'] // you can set roles in root nav
                }
              }
            ]
          },
          {
            path: 'settings',
            name: 'AdminSetitngs',
            layout: 'Layout',
            hidden: false,
            component: AdminSettingsLayout,
            meta: {
              title: 'Settings',
              hideMainMenu: true,
              showMenu: adminPluginNameMenu,
              icon: 'lock',
              roles: ['admin'] // you can set roles in root nav
            },
            children: [
              {
                path: 'info',
                name: 'AdminSettingsInfo',
                layout: 'Layout',
                component: AdminSettingsIndex,
                hidden: false,
                meta: {
                  title: 'Info',
                  hideMainMenu: true,
                  showMenu: adminPluginNameMenu,
                  icon: 'lock',
                  roles: ['admin'] // you can set roles in root nav
                }
              },
              {
                path: 'example',
                name: 'AdminSettingsExample',
                layout: 'Layout',
                component: AdminSettingsExample,
                hidden: false,
                meta: {
                  title: 'Example',
                  hideMainMenu: true,
                  showMenu: adminPluginNameMenu,
                  icon: 'lock',
                  roles: ['admin'] // you can set roles in root nav
                }
              }
            ]
          }
        ]
      }
    ]
  },
  {
    path: '/builder',
    // component: Layout,
    layout: 'Layout',
    redirect: '/builder/v1',
    alwaysShow: true, // will always show the root menu
    meta: {
      title: 'BUilder',
      icon: 'lock',
      roles: ['admin', 'editor'] // you can set roles in root nav
    },
    children: [
      {
        path: 'v1',
        component: () => import('@/views/builder/Index'),
        name: 'BuildIndex',
        meta: {
          title: 'Builder',
          roles: ['admin', 'editor'] // or you can only set roles in sub nav
        }
      },
      {
        path: 'v2',
        component: () => import('@/views/builderv2/Index'),
        name: 'BuildIndexV2',
        meta: {
          title: 'BuilderV2',
          roles: ['admin', 'editor'] // or you can only set roles in sub nav
        }
      }
    ]
  }
  /*
  {
    path: '/create',
    layout: 'Layout',
    hidden: true,
    meta: {
      title: 'create',
      icon: 'lock',
      roles: ['admin', 'editor'] // you can set roles in root nav
    },
    children: [
      {
        path: '/create/ami',
        name: 'createami',
        props: true,
        layout: 'Layout',
        component: Setup,
        hidden: true,
        meta: {
          title: 'ami',
          icon: 'lock',
          roles: ['admin', 'editor'] // you can set roles in root nav
        }
      },
      {
        path: '/create/cluster',
        layout: 'Layout',
        plugin: plugName,
        // innerMenu: plugName,
        name: 'createcluster',
        component: SetupCluster,
        hidden: true,
        meta: {
          title: 'Create',
          hideMainMenu: false,
          // showMenu: false,
          icon: 'lock',
          roles: ['admin', 'editor'] // you can set roles in root nav
        },
        children: [
        ]
      }
    ]
  },
  {
    path: '/cluster/:id',
    layout: 'Layout',
    plugin: plugName,
    innerMenu: clusterPlugName,
    hidden: false,
    flatChildrens: true,
    children: [
      {
        path: '',
        name: 'clusterRoot',
        layout: 'Layout',
        component: ClusterLayout,
        hidden: false,
        // alwaysShow: false,
        flatChildrens: true,
        meta: {
          title: 'Cluster Dashboard',
          hideMainMenu: true,
          showMenu: clusterPlugName,
          icon: 'lock',
          roles: ['admin', 'editor'] // you can set roles in root nav
        },
        children: [
          {
            path: '',
            name: 'cluster',
            layout: 'Layout',
            component: ClusterDashboard,
            hidden: false,
            meta: {
              title: 'Dashboard',
              hideMainMenu: true,
              showMenu: clusterPlugName,
              // example: class: 'serverless-cms-container',
              icon: 'lock',
              roles: ['admin', 'editor'] // you can set roles in root nav
            }
          },
          {
            path: 'discovery',
            name: 'ClusterDiscovery',
            layout: 'Layout',
            component: ClusterDiscovery,
            hidden: false,
            meta: {
              title: 'Discovery',
              hideMainMenu: true,
              showMenu: clusterPlugName,
              icon: 'lock',
              roles: ['admin', 'editor'] // you can set roles in root nav
            }
          },
          {
            path: 'nodes',
            name: 'ClusterNodes',
            layout: 'Layout',
            component: ClusterNodes,
            hidden: false,
            meta: {
              title: 'Nodes',
              hideMainMenu: true,
              showMenu: clusterPlugName,
              icon: 'lock',
              roles: ['admin', 'editor'] // you can set roles in root nav
            }
          },
          {
            path: 'settings',
            name: 'ClusterSettings',
            layout: 'Layout',
            hidden: false,
            component: ClusterLayout,
            meta: {
              title: 'Settings',
              hideMainMenu: true,
              showMenu: clusterPlugName,
              icon: 'lock',
              roles: ['admin', 'editor'] // you can set roles in root nav
            },
            children: [
              {
                path: 'admin',
                name: 'ClusterSettingsIndex',
                layout: 'Layout',
                component: ClusterSettingsAdmin,
                hidden: false,
                meta: {
                  title: 'Admin',
                  hideMainMenu: true,
                  showMenu: clusterPlugName,
                  icon: 'lock',
                  roles: ['admin', 'editor'] // you can set roles in root nav
                }
              },
              {
                path: 'danger',
                name: 'ClusterSettingsDanger',
                layout: 'Layout',
                component: ClusterSettingsDanger,
                hidden: false,
                meta: {
                  title: 'Danger Zone',
                  hideMainMenu: true,
                  showMenu: clusterPlugName,
                  icon: 'lock',
                  roles: ['admin', 'editor'] // you can set roles in root nav
                }
              }
            ]
          }
        ]
      }
    ]
  },
  {
    path: '/ami/:id',
    layout: 'Layout',
    plugin: plugName,
    innerMenu: amiPlugName,
    hidden: false,
    flatChildrens: true,
    children: [
      {
        path: '',
        name: 'amiRoot',
        layout: 'Layout',
        component: AmiLayout,
        hidden: false,
        // alwaysShow: false,
        flatChildrens: true,
        meta: {
          title: 'AMI Dashboard',
          hideMainMenu: true,
          showMenu: amiPlugName,
          icon: 'lock',
          roles: ['admin', 'editor'] // you can set roles in root nav
        },
        children: [
          {
            path: '',
            name: 'ami',
            layout: 'Layout',
            component: AmiDashboard,
            hidden: false,
            meta: {
              title: 'Dashboard',
              hideMainMenu: true,
              showMenu: amiPlugName,
              // example: class: 'serverless-cms-container',
              icon: 'lock',
              roles: ['admin', 'editor'] // you can set roles in root nav
            }
          },
          {
            path: 'settings',
            name: 'AmiSettings',
            layout: 'Layout',
            hidden: false,
            component: AmiLayout,
            meta: {
              title: 'Settings',
              hideMainMenu: true,
              showMenu: amiPlugName,
              icon: 'lock',
              roles: ['admin', 'editor'] // you can set roles in root nav
            },
            children: [
              {
                path: 'admin',
                name: 'AmiSettingsIndex',
                layout: 'Layout',
                component: AmiSettingsAdmin,
                hidden: false,
                meta: {
                  title: 'Admin',
                  hideMainMenu: true,
                  showMenu: amiPlugName,
                  icon: 'lock',
                  roles: ['admin', 'editor'] // you can set roles in root nav
                }
              },
              {
                path: 'danger',
                name: 'AmiSettingsDanger',
                layout: 'Layout',
                component: AmiSettingsDanger,
                hidden: false,
                meta: {
                  title: 'Danger Zone',
                  hideMainMenu: true,
                  showMenu: amiPlugName,
                  icon: 'lock',
                  roles: ['admin', 'editor'] // you can set roles in root nav
                }
              }
            ]
          }
        ]
      }
    ]
  }
  */
]

export const routerInvalid = [
  {
    path: '/',
    component: DashboardIndex
  },
  {
    path: '/dashboard/index',
    // component: Layout,
    layout: 'Layout',
    replacePrecedence: 100,
    replacePath: '/dashboard',
    hidden: true,
    children: [
      {
        path: 'index',
        component: DashboardIndex,
        name: 'Tab',
        meta: { title: 'tab', icon: 'tab', hidden: true }
      }
    ]
  }
]

const initRouter = function() {
  const r = window.app.$router || new Router({
    // routerSetup,
    actions: {
      RouteChange() {
        console.error('route change aaa FA ZZ')
      }
    }
  })
  return r
}
const router = initRouter()

export default router
