<template>
  <div v-if="!isLoading" class="dashboard-editor-container minvh">

    <el-row :gutter="8">
      <el-col :xs="{span: 24}" :sm="{span: 24}" :md="{span: 24}" :lg="{span: 12}" :xl="{span: 12}" style="padding-right:8px;margin-bottom:30px;">
        <h1>List of AMI Pipeline</h1>
        <plugin-list-table @onRowSelect="onRowSelect"/>
      </el-col>
      <el-col :xs="{span: 24}" :sm="{span: 12}" :md="{span: 12}" :lg="{span: 12}" :xl="{span: 12}" style="margin-bottom:30px;">
        <h1>List of Clusters</h1>
        <project-list-table @onRowSelect="onProjectRowSelect"/>
      </el-col>
      <el-col :xs="{span: 24}" :sm="{span: 12}" :md="{span: 12}" :lg="{span: 12}" :xl="{span: 12}" style="margin-bottom:30px;">
        <br>
      </el-col>
    </el-row>

  </div>
</template>

<script>
import permission from '@/directive/permission/index.js' // 权限判断指令
import checkPermission from '@/utils/permission' // 权限判断函数
import PluginListTable from './components/PluginListTable'
import ProjectListTable from './components/ProjectListTable'
import { mapGetters } from 'vuex'

export default{
  name: 'DirectivePermission',
  components: {
    PluginListTable,
    ProjectListTable
  },
  directives: { permission },
  data() {
    return {
      key: 1 // 为了能每次切换权限的时候重新初始化指令
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
  methods: {
    checkPermission,
    handleRolesChange() {
      this.key++
    },
    onRowSelect(r) {
      console.error('received forw ofw')
      console.error(r)

      // this.$router.push({ name: 'ami', params: { amiid: r.id }})
      /*
      this.$router.push({
            name: 'ami',
            params: {
                amiid: r.id,
                otherProp: {
                    "a": "b"
                }
            }
        })
      */
      this.$router.push({
        name: 'ami',
        params: {
          id: r.id,
          data: r
        }
      })
      // this.$router.push({ path: '/ami/' + r.id })

      // this.$router.push("/ami/" + row.id, r)
    },
    onProjectRowSelect(r) {
      this.$router.push({
        name: 'cluster',
        params: {
          id: r.id,
          data: r
        }
      })
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
.app-container {
  /deep/ .permission-alert {
    width: 320px;
    margin-top: 30px;
    background-color: #f0f9eb;
    color: #67c23a;
    padding: 8px 16px;
    border-radius: 4px;
    display: block;
  }
  /deep/ .permission-tag{
    background-color: #ecf5ff;
  }
}
.dashboard-editor-container {
  padding: 32px;
  background-color: rgb(240, 242, 245);
  .chart-wrapper {
    background: #fff;
    padding: 16px 16px 0;
    margin-bottom: 32px;
  }
}
.minvh {
  min-height: 95vh;
}

</style>
