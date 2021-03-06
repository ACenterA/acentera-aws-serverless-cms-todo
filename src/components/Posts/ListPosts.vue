<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input v-model="listQuery.title" :placeholder="$t('table.title')" style="width: 200px;" class="filter-item" @keyup.enter.native="handleFilter" />
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-edit" @click="handleCreate">{{ $t('table.add') }}</el-button>
      listing posts
    </div>

    <el-table
      :key="tableKey"
      v-loading="listLoading"
      :default-sort="{prop: 'created', order: 'descending'}"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;"
      @sort-change="sortChange"
    >
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
          <el-button type="primary" size="mini" @click="handleUpdate(scope.row)">{{ $t('table.edit') }}</el-button>
          <!--
          <el-button v-if="scope.row.status!='published'" size="mini" type="success" @click="handleModifyStatus(scope.row,'published')">{{ $t('table.publish') }}
          </el-button>
          <el-button v-if="scope.row.status!='draft'" size="mini" @click="handleModifyStatus(scope.row,'draft')">{{ $t('table.draft') }}
          </el-button>
          -->
          <el-button v-if="scope.row.status!='deleted'" size="mini" type="danger" @click="handleModifyStatus(scope.row,'deleted')">{{ $t('table.delete') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination v-show="total>0" :total="total" :page.sync="listQuery.page" :limit.sync="listQuery.limit" @pagination="getList" />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form ref="dataForm" :rules="rules" :model="temp" label-position="left" label-width="70px" style="width: 400px; margin-left:50px;">
        <el-form-item :label="$t('table.title')" prop="title">
          <el-input v-model="temp.title" />
        </el-form-item>
        <el-form-item :label="$t('table.status')">
          <el-select v-model="temp.status" class="filter-item" placeholder="Please select">
            <el-option v-for="item in statusOptions" :key="item" :label="item" :value="item" />
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
        <el-table-column prop="key" label="Channel" />
        <el-table-column prop="pv" label="Pv" />
      </el-table>
      <span slot="footer" class="dialog-footer">
        <el-button type="primary" @click="dialogPvVisible = false">{{ $t('table.confirm') }}</el-button>
      </span>
    </el-dialog>

  </div>
</template>

<script>
import waves from '@/directive/waves' // Waves directive
import { parseTime } from '@/utils'
import Pagination from '@/components/Pagination' // Secondary package based on el-pagination

import { updatePost, listPosts, CreatePostMutation } from '@/gql/queries/posts.gql'

export default {
  name: 'ComplexTable',
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
    project: {
      type: String,
      default: null
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
      query: () => listPosts,
      variables() {
        return {
          project: this.project,
          nextToken: this.nextToken
        }
      },
      update(data) {
        if (!(data.listPosts && data.listPosts.items)) {
          return []
        }
        this.list = data.listPosts.items
        this.total = data.listPosts.items.length
        this.listLoading = false
        return data.listPosts
      },
      /*
      // example if we want to have live updates
      subscribeToMore: {
        document: onCreateProject,
        // Variables passed to the subscription. Since we're using a function,
        // they are reactive
        //variables () {
        //  return {
        //    param: this.param,
        //  }
        // },
        // Mutate the previous result
        updateQuery: (previousResult, { subscriptionData }) => {
          console.error('recive de vent ehre')
          const previousItems = previousResult.listPosts.items

          // Here, return the new result from the previous with the new data
          var newData = subscriptionData.data.onCreateProject
          previousItems.push(newData)
        },
      },
      */
      fetchPolicy: 'cache-and-network' // 'network-only', // skip the cache
      /*
      query: () => listProjects,
      fetchPolicy: 'cache-and-network' // 'network-only', // skip the cache
      variables: function() {
        return {
          limit: this.$data.listQuery.limit
        }
      },
      update(data) {
        this.list = data.listProjects.items
        this.total = data.listProjects.items.length
        this.listLoading = false
        return data.listProjects
      }
      */
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
    async updatePost(data) {
      const result = await new Promise((resolve, reject) => {
        this.$apollo.mutate({
          mutation: updatePost,
          variables: {
            id: data.id,
            title: data.title
          },
          update: (store, { data: { updatePost }}) => {
            // Update graph db store ?
            const data = store.readQuery({ query: listPosts })
            const index = data.listPosts.items.findIndex(item => item.id === updatePost.id)
            data.listPosts.items[index] = updatePost
            store.writeQuery({ query: listPosts, data })
          },
          optimisticResponse: {
            __typename: 'Mutation',
            updatePost: {
              __typename: 'Post',
              ...data
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
        description: '',
        status: 'open',
        type: ''
      }
    },
    handleCreate() {
      this.resetTemp()
      this.$router.push({ 'name': 'PostCreate' })
      /*
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
      */
    },
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
            mutation: CreatePostMutation,
            variables: {
              title,
              project,
              status,
              completed
            },
            // see https://github.com/Akryum/vue-apollo-todos
            update: (store, { data: { createTask }}) => {
              /*
              // Add to ALL PROJECT list
              var queryKey = ALL_PROJECTS
              var itemKey = 'ListPosts'

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
              */
            }
            /*
            optimisticResponse: {
             createTask: {
                title: title,
                created: new Date(),
                status: 'active',
                id: optimisticId,
                __typename: 'Project',
             }
            }
            */
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
          this.updatePost(tempData).then(() => {
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
