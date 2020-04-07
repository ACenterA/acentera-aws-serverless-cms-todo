<template>
  <div v-if="!isLoading" class="minvh createPost-container">
    <el-form ref="postForm" :model="postForm" class="form-container">
      <h3>Create Project</h3>
      <div class="createPost-main-container">
        <el-row>
          <el-col :span="24">
            <el-form-item style="margin-bottom: 40px;" prop="title">
              <MDinput v-model="postForm.title" :maxlength="100" name="name" required>
                Name
              </MDinput>
            </el-form-item>
          </el-col>
        </el-row>
      </div>
    </el-form>
    <el-button v-loading="loading" type="primary" @click="createProject()">Submit</el-button>
  </div>
</template>

<script>
import permission from '@/directive/permission/index.js' // 权限判断指令
// import checkPermission from '@/utils/permission' // 权限判断函数
// import PluginListTable from './components/PluginListTable'
import { mapGetters } from 'vuex'
import moment from 'moment'
import { ALL_PROJECTS, CreateProjectMutation } from '@/gql/queries/projects.gql'
import MDinput from '@/components/MDinput'

export default {

  components: {
    MDinput
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
  data() {
    return {
      loading: false,
      postForm: {
        title: ''
      }
    }
  },
  computed: {
    ...mapGetters([
      'isLoading',
      'name',
      'avatar',
      'roles'
    ]),
    activeClient() { return this.$store.getters.apollo }
  },
  apollo: {
    $client() {
      console.error('APOLLO CLIENT')
      console.error(window.Apollo)
      console.error('APOLLO CLIENT EST ')
      console.error(this.activeClient)
      // console.error(this.$store.getters.apollo)
      // return this.$store.getters.apollo
      // return window.Apollo
      return this.activeClient // window.Apollo // this.activeClient
    }
  },
  mounted() {
  },
  created() {
  },
  methods: {
    createProject() {
      const { title } = this.$data.postForm
      const optimisticId = '' + Math.round(Math.random() * -1000000)
      var status = 'active'
      this.loading = true
      this.$apollo.mutate({
        mutation: CreateProjectMutation,
        variables: {
          title,
          status
        },
        // see https://github.com/Akryum/vue-apollo-todos
        update: (store, { data: { createProject }}) => {
          // Add to ALL PROJECT list
          var queryKey = ALL_PROJECTS
          var itemKey = 'listAllProjects'

          // Ensure no duplicas (offline graphql ...) if not using optimisticId logic
          // filtering doesnt seems to work with nested elements ???
          const queryWithFilteringOptimistic = {
            query: queryKey,
            // variables: { filter: { items: { id: optimisticId } } }
            variables: { filter: { items: { id: '' + createProject.id }}}
          }

          /*
          if (store.data.data) {
            // Ensure to delete offline cached bad data...
            if (createProject.id !== optimisticId) {
              delete store.data.data[createProject.__typename +':' + optimisticId]
            }
          }
          */

          console.error('received data filtering of ' + optimisticId + ' and of ' + createProject.id)
          try {
            if (createProject.id !== optimisticId) {
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
                      if (tmpAddItem[i].id === createProject.id) {
                        found = true
                      }
                    }
                  } catch (z) {
                    console.error(z.stack)
                  }

                  if (!found) {
                    tmpAddItem.push(createProject)
                  }
                  store.writeQuery({ ...queryWithFilteringOptimistic, data: filteredData })
                }
              }
            }
          } catch (err) {
            console.error(err)
          }
        }
        /*
        optimisticResponse: {
         createProject: {
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
