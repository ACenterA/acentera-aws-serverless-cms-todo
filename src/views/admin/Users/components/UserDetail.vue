<template>
  <div class="createUser-container">
    <el-form ref="userForm" :model="userForm" :rules="rules" class="form-container">

      <sticky :class-name="'sub-navbar '+userForm.status">
        <el-button v-loading="loading" style="margin-left: 10px;" type="success" @click="submitForm">{{ $t('submit') }}
        </el-button>
      </sticky>
      <div class="createUser-main-container">
        <el-row>

          <!--
            <Warning />
          -->

          <el-col :span="24">
            <el-form-item style="margin-bottom: 40px;" prop="title">
              <MDinput v-model="userForm.title" :maxlength="100" name="title" required>
                {{ $t('user.title') }}
              </MDinput>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item style="margin-bottom: 40px;" prop="email">
          <MDinput v-model="userForm.email" :maxlength="100" name="email" required>
            {{ $t('user.email') }}
          </MDinput>
        </el-form-item>

        <!--
          <RoleDropdown v-model="userForm.role" />
        -->
        <list-projects @onInput="onInput($event)" />
      </div>
    </el-form>

  </div>
</template>

<script>
import MDinput from '@/components/MDinput'
import Sticky from '@/components/Sticky'
import { validateEmail } from '@/utils/validate'
// import Warning from './Warning'
import { ListProjects } from './Dropdown'
import { createUserMutation } from '@/gql/queries/users.gql'

const defaultForm = {
  status: 'active',
  title: '',
  email: '',
  role: 'editor',
  id: undefined,
  projects: []
}

export default {
  name: 'UserDetail',
  components: {
    MDinput, Sticky, ListProjects
    // , RoleDropdown
  },
  props: {
    isEdit: {
      type: Boolean,
      default: false
    }
  },
  data() {
    const validateRequire = (rule, value, callback) => {
      if (value === '') {
        console.error(rule)
        if (rule.error) {
          this.$message({
            message: this.$t('error.input.' + rule.field),
            type: rule.error_type || 'error'
          })
          callback(new Error(this.$t('error.input.' + rule.field)))
        } else {
          this.$message({
            message: this.$t('error.input.' + rule.field),
            type: rule.error_type || 'error'
          })
          callback(new Error(this.$t('error.input.' + rule.field)))
        }
      } else {
        callback()
      }
    }
    const validateRequireEmail = (rule, value, callback) => {
      if (!validateEmail(value)) {
        callback(new Error(this.$t('error.input.' + rule.field)))
      } else {
        callback()
      }
    }
    return {
      userForm: Object.assign({}, defaultForm),
      loading: false,
      userListOptions: [],
      rules: {
        email: [{
          validator: validateRequireEmail
        }],
        // image_uri: [{ validator: validateRequire }],
        title: [{
          validator: validateRequire
        }],
        content: [{ validator: validateRequire }]
        // source_uri: [{ validator: validateSourceUri, trigger: 'blur' }]
      }
    }
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
  computed: {
    activeClient() { return this.$store.getters.apollo },
    lang() {
      return this.$store.getters.language
    },
    projectid() {
      return this.$store.getters.projectid
    },
    isValid() {
      return false
    },
    form_email_placeholder() {
      return this.$t('user.email_placeholder')
    }
  },
  created() {
    this.userForm = Object.assign({}, defaultForm)
  },
  methods: {
    onInput(event) {
      console.error('wtf')
      this.userForm.projects = event.map(function(product) {
        return {
          id: product.id,
          role: product.role
        }
      })
      console.error('wtf 1 ')
      console.error(this.userForm)
    },
    fetchData(id) {
      /*
      fetchArticle(id).then(response => {
        this.userForm = response.data
        // Just for test
        this.userForm.title += `   Article Id:${this.userForm.id}`
        this.userForm.email += `   Article Id:${this.userForm.id}`
        // Set tagsview title
      }).catch(err => {
        console.log(err)
      })
      */
    },
    submitForm() {
      this.userForm.display_time = parseInt(this.display_time / 1000)
      // console.log(this.userForm)
      var _self = this
      this.$refs.userForm.validate(valid => {
        if (valid) {
          // const { title } = this.$data.userForm
          // const optimisticId = '' + Math.round(Math.random() * -1000000)
          var userVariables = this.userForm
          userVariables.project = _self.projectid
          // var status = 'active'
          this.loading = true
          console.error(userVariables)
          try {
            this.$apollo.mutate({
              mutation: createUserMutation,
              variables: userVariables,
              // see https://github.com/Akryum/vue-apollo-todos
              update: (store, { data: { createUser }}) => {
                /*
                // Add to ALL PROJECT list
                var queryKey = ALL_PROJECTS
                var itemKey = 'listAllProjects'

                // Ensure no duplicas (offline graphql ...) if not using optimisticId logic
                // filtering doesnt seems to work with nested elements ???
                const queryWithFilteringOptimistic = {
                  query: queryKey,
                  // variables: { filter: { items: { id: optimisticId } } }
                  variables: { filter: { items: { id: '' + createUser.id }}}
                }

                console.error('received data filtering of ' + optimisticId + ' and of ' + createUser.id)
                try {
                  if (createUser.id !== optimisticId) {
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
                            if (tmpAddItem[i].id === createUser.id) {
                              found = true
                            }
                          }
                        } catch (z) {
                          console.error(z.stack)
                        }

                        if (!found) {
                          tmpAddItem.push(createUser)
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
               createUser: {
                  title: title,
                  created: new Date(),
                  status: 'active',
                  id: optimisticId,
                  __typename: 'Project',
               }
              }
              */
            }).then(() => {
              this.loading = false
              this.$notify({
                title: this.$t('create.title'),
                message: this.$t('create.success.message'),
                type: 'success',
                duration: 2000
              })
              this.loading = false
              // this.$router.push({ name: 'AdminProjectsList' })
            }).catch(error => {
              this.loading = false
              this.$notify({
                title: this.$t('create.title'),
                subtitle: '',
                message: this.$t('create.error.message', { error: error }),
                type: 'error',
                duration: 5000
              })
              this.loading = false
            })
          } catch (ex) {
            console.error(ex.stack)
            this.loading = false
            this.$notify({
              title: 'GraprhQL Error',
              message: 'Error with graphql schema',
              type: 'danger',
              duration: 3000
            })
            console.log('error submit!!')
            this.loading = false
            return false
          }
        } else {
          this.loading = false
          this.$notify({
            title: 'Validation failed',
            message: 'Missing required fields',
            type: 'warning',
            duration: 2000
          })
          console.log('error submit!!')
          this.loading = false
          return false
        }
      })
    },
    submitFormTmp() {
      this.userForm.display_time = parseInt(this.display_time / 1000)
      console.log(this.userForm)
      this.$refs.userForm.validate(valid => {
        console.error(this.userForm)
        if (valid) {
          this.loading = true
          this.$notify({
            title: 'Creation',
            message: 'Successfully created user',
            type: 'success',
            duration: 2000
          })
          // this.userForm.status = 'published'
          this.loading = false
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    draftForm() {
      if (this.userForm.content.length === 0 || this.userForm.title.length === 0) {
        this.$message({
          message: 'No content or user title',
          type: 'warning'
        })
        return
      }
      this.$message({
        message: 'Draft created...',
        type: 'success',
        showClose: true,
        duration: 1000
      })
      // this.userForm.status = 'draft'
    },
    getRemoteUserList(query) {
      /*
      userSearch(query).then(response => {
        if (!response.data.items) return
        this.userListOptions = response.data.items.map(v => v.name)
      })
      */
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
@import "~@/styles/mixin.scss";
.createUser-container {
  position: relative;
  .createUser-main-container {
    padding: 40px 45px 20px 50px;
    .userInfo-container {
      position: relative;
      @include clearfix;
      margin-bottom: 10px;
      .userInfo-container-item {
        float: left;
      }
    }
  }
  .word-counter {
    width: 40px;
    position: absolute;
    right: -10px;
    top: 0px;
  }
}
</style>
