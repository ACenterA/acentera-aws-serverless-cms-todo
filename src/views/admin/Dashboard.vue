<template>
  <div v-if="!isLoading" class="minvh">
    Dashboard ...
    <div v-if="project && project.id">
      Current Project: {{ project.title }}
    </div>
  </div>
</template>

<script>
import permission from '@/directive/permission/index.js' // 权限判断指令
// import checkPermission from '@/utils/permission' // 权限判断函数
// import PluginListTable from './components/PluginListTable'
import { mapGetters } from 'vuex'
import moment from 'moment'

export default {

  components: {
  },
  directives: { permission },
  filters: {
    dateToStr(str) {
      if (str) {
        var localTime = moment(str, 'YYYY-MM-DDTHH24:MI:ss.SSSZ').local()
        return localTime.format('YYYY/MM/DD hh:mm')
      }
    }
  },
  props: {
    id: {
      type: String,
      default: ''
    },
    data: {
      type: Object,
      default: function() {
        return null
      }
    },
    stack: {
      type: Object,
      default: function() {
        return null
      }
    }
  },
  data() {
    return {
      loading: false,
      cluster: this.data,
      buildedAmi: [],
      key: 1
    }
  },
  computed: {
    ...mapGetters([
      'isLoading',
      'name',
      'project',
      'avatar',
      'roles'
    ]),
    clusterObject() {
      if (this.cluster) {
        return [this.cluster]
      }
      return []
    }
  },
  mounted() {
    console.error('fa mounted fe')
    console.error(this.$parent)
  },
  created() {
    console.error('fa created fe')
    console.error(this.$parent.cluster)
    console.error(this.$parent.data)
  },
  methods: {
  }
}
</script>

<style>
 #bad_nprogress {
    opacity: 0.4;
    background-color: #a3a3a3;
    pointer-events: none;
    position: fixed;
    top: 0px;
    left: 0px;
    min-height: 100vh;
    width: 100%;
 }
</style>

<style rel="stylesheet/scss" lang="scss" scoped>

</style>
