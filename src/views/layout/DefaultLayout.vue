<!-- TOOD: We could totally override the app if needed? -->
<template>
  <div v-if="!isLoading">
    <router-view :loading="loading"/>
  </div>
</template>

<script>
// import permission from '@/directive/permission/index.js' // 权限判断指令
import checkPermission from '@/utils/permission' // 权限判断函数
import { mapGetters } from 'vuex'

export default {
  name: 'DefaultPluginInnerLayout',
  components: {
  },
  mixins: [],
  data() {
    return {
      loading: false
    }
  },
  computed: {
    ...mapGetters([
      'isLoading',
      'name',
      'avatar',
      'roles'
    ])
  },
  mounted() {
    // TODO: Optimization: If received a Data object...
  },
  methods: {
    checkPermission
  }
}
</script>

<style rel="stylesheet/scss" lang="scss">
</style>
