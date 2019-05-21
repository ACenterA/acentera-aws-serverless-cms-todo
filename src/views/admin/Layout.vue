<!-- TOOD: We could totally override the app if needed? -->
<template>
  <div v-if="!isLoading">
    <!--
    If we want to show thits on all pages under this layout
    Welcome Administrator!
    <router-view v-if="isSettingRoute" :key="key" :loading="loading"/>
    -->
    <router-view :loading="loading"/>
  </div>
</template>

<script>
// import permission from '@/directive/permission/index.js' // 权限判断指令
// import PluginListTable from './components/PluginListTable'
import checkPermission from '@/utils/permission' // 权限判断函数
// import ClusterEcsFirstStep from './components/ClusterEcsFirstStep.vue'
import { mapGetters } from 'vuex'

export default {
  name: 'LayoutAdministration',
  components: {
    // ClusterEcsFirstStep
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
    ]),
    isSettingRoute() {
      if (this.$route.name.match(/^(.*ettings).*$/)) { // must match ettings ie: ClusterSettings
        return true
      }
      return false
    }
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
.example-metat-plugin-ecseks-container {
  background-color: rgb(240, 242, 245);
}
</style>

<style>
.wizard-tab-content {
  max-width: 650px;
  margin: auto;
}

.wizard-card-footer {
  max-width: 650px;
  margin: auto;
}

.wizard-btn.wizard-footer-right.invalid {
  background-color: white!important;
  color: black!important;
}

.wizard-tab-container {
  min-height:300px;
}
</style>
