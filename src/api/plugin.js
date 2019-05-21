import request from '@/utils/request'

export function getPluginSettings() {
  return request({
    url: `/plugins/${PLUGIN_NAME}/settings`,  // eslint-disable-line
    method: 'get'
  }, {
    cache: true,
    maxAge: 15 * 60 * 1000
  })
}

export function validateAccountId(data) {
  return request({
    url: '/app/setup/validate',
    method: 'post',
    data
  })
}

export function performAppInitialization(data) {
  return request({
    url: `/plugins/${PLUGIN_NAME}/setup/bootstrap`,  // eslint-disable-line
    method: 'post',
    data
  })
}
