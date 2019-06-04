<template>
  <div class="app-container">
    <div class="filter-container">
      <!--
         <el-input :placeholder="$t('table.title')" v-model="listQuery.title" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter"/>
      -->
    </div>

    <el-table
      v-loading="listLoading"
      :key="tableKey"
      :default-sort = "{prop: 'created', order: 'descending'}"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;"
      @sort-change="sortChange">
      <el-table-column :label="$t('table.date')" prop="created" sortable class="col-sm" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.created | parseTime('{y}-{m}-{d} {h}:{i}') }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('table.title')" class="col-sm">
        <template slot-scope="scope">
          <span class="text-type">{{ scope.row.title }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('table.status')" class-name="status-col">
        <template slot-scope="scope">
          <el-tag :type="scope.row.status | statusFilter">{{ scope.row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('table.role')" align="center" width="200" class-name="small-padding fixed-width">
        <template slot-scope="scope">
          <RoleDropdown v-model="scope.row.role" @input="onInput($event, scope)"/>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
// import { fetchList, fetchPv, createArticle, updateArticle } from '@/api/projects'

import Vue from 'vue'
import { RoleDropdown } from './index'

import waves from '@/directive/waves' // Waves directive
import { parseTime } from '@/utils'

// import { ALL_PROJECTS, deleteProject, updateProject, onCreateProject } from '@/gql/queries/projects.gql'
import { ALL_PROJECTS, deleteProject, updateProject } from '@/gql/queries/projects.gql'
// import { listProjects } from '@/gql/queries/projects.gql'
import { mapGetters } from 'vuex'

export default {
  name: 'ComplexTable',
  components: { RoleDropdown },
  directives: { waves },
  filters: {
    statusFilter(status) {
      const statusMap = {
        active: 'success',
        inactive: 'danger'
      }
      return statusMap[status]
    }
  },
  data() {
    return {
      tableKey: 0,
      sortKey: 'created',
      sortOrder: 'descending',
      list: [],
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        limit: 2000,
        importance: undefined,
        title: undefined,
        type: undefined,
        sort: '+id'
      },
      importanceOptions: [1, 2, 3],
      sortOptions: [{ label: 'ID Ascending', key: '+id' }, { label: 'ID Descending', key: '-id' }],
      statusOptions: ['active', 'inactive'],
      showReviewer: false,
      temp: {
        id: undefined,
        importance: 1,
        remark: '',
        timestamp: new Date(),
        title: '',
        type: '',
        status: 'published'
      },
      dialogFormVisible: false,
      dialogStatus: '',
      textMap: {
        update: 'Edit',
        create: 'Create'
      },
      dialogPvVisible: false,
      pvData: [],
      rules: {
        type: [{ required: true, message: 'type is required', trigger: 'change' }],
        timestamp: [{ type: 'date', required: true, message: 'timestamp is required', trigger: 'change' }],
        title: [{ required: true, message: 'title is required', trigger: 'blur' }]
      },
      downloadLoading: false
    }
  },
  computed: {
    ...mapGetters([
      'projects'
    ]),
    activeClient() { return this.$store.getters.apollo }
  },
  apollo: {
    $client() {
      return this.activeClient
    }
  },
  created() {
    this.getList()
  },
  methods: {
    onInput(evt, b, val) {
      // this.inputVal = val;
      /*
      console.error(this)
      console.log(a);
      console.log(b);
      console.log(val);

      // const data = this.$store.readQuery({ query: ALL_PROJECTS })
      // console.error(data)
      console.error(this.list)
      console.error(b.row)
      const index = this.list.findIndex(item => item.id === b.row.id)
      if (a == 'is_editor') {
        b.row.role = 'editor'
      }
      this.list[index].role='editor' // = b.row
      var clonedList = this.list
      this.list = []
      this.list = clonedList
      console.error(b.row)
      */
      this.$emit('onInput', this.list)
    },
    getList() {
      this.listLoading = true

      this.list.length = 0
      this.list = JSON.parse(JSON.stringify(this.projects.filter(project => {
        // project.role = 'none'
        // this.$set('contacts[' + newPsgId + ']', newObj)
        Vue.set(project, 'role', 'none')
        return project.status !== 'deleted'
      })))
      this.$emit('onInput', this.list)
      this.listLoading = false
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    handleModifyStatus(row, status) {
      var modifyResponse
      console.error('row')
      row.status = status // for Optimistic response
      if (status === 'deleted') {
        modifyResponse = this.deleteProject(row)
      }
      var responseState = 'error'
      if (modifyResponse === true) {
        responseState = 'success'
      }
      this.$message({
        message: 'Deletion successful',
        type: responseState
      })
      row.status = status
    },

    async updateProject(project) {
      const result = await new Promise((resolve, reject) => {
        this.$apollo.mutate({
          mutation: updateProject,
          variables: {
            id: project.id,
            title: project.title,
            status: project.status
          },
          update: (store, { data: { updateProject }}) => {
            // Update the ALL PROJECT
            const data = store.readQuery({ query: ALL_PROJECTS })
            const index = data.listAllProjects.items.findIndex(item => item.id === updateProject.id)
            data.listAllProjects.items[index] = updateProject
            store.writeQuery({ query: ALL_PROJECTS, data })
          },
          optimisticResponse: {
            __typename: 'Mutation',
            updateProject: {
              __typename: 'Project',
              ...project
            }
          }
        }).then(data => {
          console.log(data)
          resolve(true)
        }).catch(error => {
          console.error(error)
          resolve(false)
        })
      })
      return result
    },

    async deleteProject(project) {
      const result = await new Promise((resolve, reject) => {
        this.$apollo.mutate({
          mutation: deleteProject,
          variables: {
            id: project.id
          },
          update: (store, { data: { deleteProject }}) => {
            const data = store.readQuery({ query: ALL_PROJECTS })
            data.listAllProjects.items = data.listAllProjects.items.filter(project => project.id !== deleteProject.id)
            store.writeQuery({ query: ALL_PROJECTS, data })
          },
          optimisticResponse: {
            __typename: 'Mutation',
            deleteProject: {
              __typename: 'Project',
              ...project
            }
          }
        }).then(data => {
          console.log(data)
          resolve(true)
        }).catch(error => {
          console.error(error)
          resolve(false)
        })
      })
      return result
    },
    sortChange(data) {
      const { prop, order } = data
      if (prop === 'id') {
        this.sortByID(order)
      }
    },
    sortByID(order) {
      if (order === 'ascending') {
        this.listQuery.sort = '+id'
      } else {
        this.listQuery.sort = '-id'
      }
      this.handleFilter()
    },
    resetTemp() {
      this.temp = {
        id: undefined,
        importance: 1,
        remark: '',
        timestamp: new Date(),
        title: '',
        status: 'published',
        type: ''
      }
    },
    handleCreate() {
      this.resetTemp()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          this.temp.id = parseInt(Math.random() * 100) + 1024 // mock a id
          this.temp.author = 'vue-element-admin'
          /*
          createArticle(this.temp).then(() => {
            this.list.unshift(this.temp)
            this.dialogFormVisible = false
            this.$notify({
              title: '成功',
              message: '创建成功',
              type: 'success',
              duration: 2000
            })
          })
          */
        }
      })
    },
    handleUpdate(row) {
      this.temp = Object.assign({}, row) // copy obj
      this.temp.created = new Date(this.temp.created)
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.temp)
          tempData.updated = new Date().toISOString().split('.')[0] + 'Z'
          this.updateProject(tempData).then(() => {
            for (const v of this.list) {
              if (v.id === this.temp.id) {
                const index = this.list.indexOf(v)
                this.list.splice(index, 1, this.temp)
                break
              }
            }
            this.dialogFormVisible = false
            this.$notify({
              title: 'Update project',
              message: 'The prorject was updated',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    },
    handleDelete(row) {
      this.$notify({
        title: '成功',
        message: '删除成功',
        type: 'success',
        duration: 2000
      })
      const index = this.list.indexOf(row)
      this.list.splice(index, 1)
    },
    handleFetchPv(pv) {
      /*
      fetchPv(pv).then(response => {
        this.pvData = response.data.pvData
        this.dialogPvVisible = true
      })
      */
    },
    handleDownload() {
      this.downloadLoading = true
      /*
      impo
      rt('@/vendor/Export2Excel').then(excel => {
        const tHeader = ['timestamp', 'title', 'type', 'importance', 'status']
        const filterVal = ['timestamp', 'title', 'type', 'importance', 'status']
        const data = this.formatJson(filterVal, this.list)
        excel.export_json_to_excel({
          header: tHeader,
          data,
          filename: 'table-list'
        })
        this.downloadLoading = false
      })
      */
    },
    formatJson(filterVal, jsonData) {
      return jsonData.map(v => filterVal.map(j => {
        if (j === 'timestamp') {
          return parseTime(v[j])
        } else {
          return v[j]
        }
      }))
    }
  }
}
</script>
