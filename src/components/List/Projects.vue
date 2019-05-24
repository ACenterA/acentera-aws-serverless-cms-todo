<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input :placeholder="$t('table.title')" v-model="listQuery.title" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter"/>
      <el-button v-if="allowCreate" class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit" @click="handleCreate">{{ $t('table.add') }}</el-button>
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
          <span class="link-type" @click="handleUpdate(scope.row)">{{ scope.row.title }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="$t('table.status')" class-name="status-col">
        <template slot-scope="scope">
          <el-tag :type="scope.row.status | statusFilter">{{ scope.row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="$t('table.actions')" align="center" width="200" class-name="small-padding fixed-width">
        <template slot-scope="scope">
          <el-button type="primary" size="mini" @click="handleSelect(scope.row)">{{ $t('table.select') }}</el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="70px" style="width: 400px; margin-left:50px;">
        <el-form-item :label="$t('table.title')" prop="title">
          <el-input v-model="temp.title"/>
        </el-form-item>
        <el-form-item :label="$t('table.status')">
          <el-select v-model="temp.status" class="filter-item" placeholder="Please select">
            <el-option v-for="item in statusOptions" :key="item" :label="item" :value="item"/>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">{{ $t('table.cancel') }}</el-button>
        <el-button type="primary" @click="dialogStatus==='create'?createData():updateData()">{{ $t('table.confirm') }}</el-button>
      </div>
    </el-dialog>

    <el-dialog :visible.sync="dialogPvVisible" title="Reading statistics">
      <el-table :data="pvData" border fit highlight-current-row style="width: 100%">
        <el-table-column prop="key" label="Channel"/>
        <el-table-column prop="pv" label="Pv"/>
      </el-table>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="dialogPvVisible = false">{{ $t('table.confirm') }}</el-button>
      </span>
    </el-dialog>

  </div>
</template>

<script>
// import { fetchList, fetchPv, createArticle, updateArticle } from '@/api/projects'
import waves from '@/directive/waves' // Waves directive
import { parseTime } from '@/utils'
import Pagination from '@/components/Pagination' // Secondary package based on el-pagination

import { ALL_PROJECTS } from '@/gql/queries/projects.gql'

export default {
  name: 'ListProjects',
  components: { Pagination },
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
  props: {
    'allowCreate': {
      required: true,
      type: Boolean,
      default: true
    }
  },
  data() {
    return {
      tableKey: 0,
      sortKey: 'created',
      sortOrder: 'descending',
      list: null,
      total: 0,
      nextToken: null,
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
      statusOptions: ['open', 'pending', 'in-progress', 'closed'],
      showReviewer: false,
      temp: {
        id: undefined,
        importance: 1,
        remark: '',
        description: '',
        timestamp: new Date(),
        title: '',
        type: '',
        status: 'open'
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
    activeClient() { return this.$store.getters.apollo }
  },
  apollo: {
    $client() {
      return this.activeClient
    },
    projects: {
      query: () => ALL_PROJECTS,
      variables() {
        return {
        }
      },
      update(data) {
        console.error(data.listAllProjects)
        this.list = data.listAllProjects.items
        this.total = data.listAllProjects.items.length
        this.listLoading = false
        this.$store.dispatch('setProjects', this.list)
        return data.listAllProjects
      },
      fetchPolicy: 'cache-and-network' // 'network-only', // skip the cache
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      this.listLoading = true
      /*
      fetchList(this.listQuery).then(response => {
        this.list = response.data.items
        this.total = response.data.total

        // Just to simulate the time of the request
        setTimeout(() => {
          this.listLoading = false
        }, 1.5 * 1000)
      })
      */
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    /*
    handleModifyStatus(row, status) {
      var modifyResponse
      console.error('row')
      row.status = status // for Optimistic response
      if (status === 'deleted') {
        modifyResponse = this.deleteTask(row)
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
    */
    /*
    async updateTask(project) {
      const result = await new Promise((resolve, reject) => {
        this.$apollo.mutate({
          mutation: updateTask,
          variables: {
            id: project.id,
            title: project.title,
            status: project.status
          },
          update: (store, { data: { updateTask }}) => {
            // Update graph db store ?
            const data = store.readQuery({ query: listTasks })
            const index = data.listTasks.items.findIndex(item => item.id === updateTask.id)
            data.listTasks.items[index] = updateTask
            store.writeQuery({ query: listTasks, data })
          },
          optimisticResponse: {
            __typename: 'Mutation',
            updateTask: {
              __typename: 'Task',
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
    */
    /*
    async deleteTask(task) {
      const result = await new Promise((resolve, reject) => {
        this.$apollo.mutate({
          mutation: deleteTask,
          variables: {
            id: task.id
          },
          update: (store, { data: { deleteTask }}) => {
            const data = store.readQuery({ query: listTasks })
            data.listTasks.items = data.listTasks.items.filter(task => task.id !== deleteTask.id)
            store.writeQuery({ query: listTasks, data })
          },
          optimisticResponse: {
            __typename: 'Mutation',
            deleteProject: {
              __typename: 'Task',
              ...task
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
    */
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
        description: '',
        status: 'open',
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
    /*
    createData() {
      // const dat = this.$data.postForm
      const project = this.project
      const apollo = this.$apollo
      console.error('her this')
      console.error(this)
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          console.error('her this then')
          console.error(this.temp)
          this.temp.id = parseInt(Math.random() * 100) + 1024 // mock a id
          this.temp.author = 'vue-element-admin'
          const title = this.temp.title
          // const optimisticId = '' + Math.round(Math.random() * -1000000)
          var status = 'active'
          var completed = false
          this.loading = true
          apollo.mutate({
            mutation: CreateTaskMutation,
            variables: {
              title,
              project,
              status,
              completed
            },
            // see https://github.com/Akryum/vue-apollo-todos
            update: (store, { data: { createTask }}) => {

              // Add to ALL PROJECT list
              var queryKey = ALL_PROJECTS
              var itemKey = 'listTask'

              // Ensure no duplicas (offline graphql ...) if not using optimisticId logic
              // filtering doesnt seems to work with nested elements ???
              const queryWithFilteringOptimistic = {
                query: queryKey,
                // variables: { filter: { items: { id: optimisticId } } }
                variables: { filter: { items: { id: '' + createTask.id }}}
              }

              console.error('received data filtering of ' + optimisticId + ' and of ' + createTask.id)
              try {
                if (createTask.id !== optimisticId) {
                  const filteredData = store.readQuery(queryWithFilteringOptimistic)
                  if (filteredData) {
                    const tmpAddItem = filteredData[itemKey].items || filteredData[itemKey]
                    console.error(tmpAddItem)
                    var maxL = tmpAddItem.length
                    if (tmpAddItem.length >= 1) {
                      var found = false
                      try {
                        for (var i = maxL - 1; !found && i > maxL - 4 && i > 0; i--) {
                          console.error(tmpAddItem[i])
                          if (tmpAddItem[i].id === createTask.id) {
                            found = true
                          }
                        }
                      } catch (z) {
                        console.error(z.stack)
                      }

                      if (!found) {
                        tmpAddItem.push(createTask)
                      }
                      store.writeQuery({ ...queryWithFilteringOptimistic, data: filteredData })
                    }
                  }
                }
              } catch (err) {
                console.error(err)
              }

            }
            optimisticResponse: {
             createTask: {
                title: title,
                created: new Date(),
                status: 'active',
                id: optimisticId,
                __typename: 'Project',
             }
            }
          }).then(() => {
            this.$notify({
              title: this.$t('project.create.title'),
              message: this.$t('project.create.success.message'),
              type: 'success',
              duration: 2000
            })
            this.loading = false
            this.$router.push({ name: 'AdminProjectsList' })
          }).catch(error => {
            this.$notify({
              title: this.$t('project.create.title'),
              subtitle: '',
              message: this.$t('project.create.error.message', { error: error }),
              type: 'error',
              duration: 5000
            })
            this.loading = false
          })
        }
      })
    },
    */
    handleSelect(row) {
      console.error(row)
      this.$store.dispatch('setProject', row.id)
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
          this.updateTask(tempData).then(() => {
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
