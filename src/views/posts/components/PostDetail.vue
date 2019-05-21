<template>
  <div class="createPost-container">
    <el-form ref="postForm" :model="postForm" :rules="rules" class="form-container">

      <sticky :class-name="'sub-navbar '+postForm.status">
        <StatusDropdown v-model="postForm.status" />
        <!--
          <CommentDropdown v-model="postForm.comment_disabled" />
          <PlatformDropdown v-model="postForm.platforms" />
          <SourceUrlDropdown v-model="postForm.source_uri" />
        -->
        <el-button v-loading="loading" style="margin-left: 10px;" type="success" @click="submitForm">{{ $t('submit') }}
        </el-button>
        <!--
        <el-button v-loading="loading" type="warning" @click="draftForm">{{ $t('table.draft') }}</el-button>
        -->
      </sticky>

      <div class="createPost-main-container">
        <el-row>

          <!--
          <Warning />
          -->

          <el-col :span="24">
            <el-form-item style="margin-bottom: 40px;" prop="title">
              <MDinput v-model="postForm.title" :maxlength="100" name="name" required>
                {{ $t('post.title') }}
              </MDinput>
            </el-form-item>

            <div class="postInfo-container">
              <el-row>
                <el-col :span="8">
                  <el-form-item :label="form_userinfo_label" label-width="45px" class="postInfo-container-item">
                    <el-select :placeholder="form_userinfo_placeholder" v-model="postForm.author" :remote-method="getRemoteUserList" filterable remote>
                      <el-option v-for="(item,index) in userListOptions" :key="item+index" :label="item" :value="item"/>
                    </el-select>
                  </el-form-item>
                </el-col>

                <el-col :span="10">
                  <el-form-item :label="form_display_time_label" label-width="80px" class="postInfo-container-item">
                    <el-date-picker :placeholder="form_display_time_placeholder" v-model="postForm.display_time" type="datetime" format="yyyy-MM-dd HH:mm:ss"/>
                  </el-form-item>
                </el-col>

                <el-col :span="6">
                  <el-form-item :label="form_author_rating_label" label-width="60px" class="postInfo-container-item">
                    <el-rate
                      v-model="postForm.importance"
                      :max="3"
                      :colors="['#99A9BF', '#F7BA2A', '#FF9900']"
                      :low-threshold="1"
                      :high-threshold="3"
                      style="margin-top:8px;"/>
                  </el-form-item>
                </el-col>
              </el-row>
            </div>
          </el-col>
        </el-row>

        <el-form-item :label="form_desc_label" style="margin-bottom: 40px;" label-width="45px">
          <el-input :rows="1" :placeholder="form_desc_placeholder" v-model="postForm.content_short" type="textarea" class="article-textarea" autosize/>
          <span v-show="contentShortLength" class="word-counter">{{ contentShortLength }}字</span>
        </el-form-item>

        <el-form-item prop="content" style="margin-bottom: 30px;">
          <Tinymce ref="editor" :height="400" v-model="postForm.content" />
        </el-form-item>

        <el-form-item :label="form_author_rating_label" prop="image_uri" style="margin-bottom: 30px;" label-width="60px">
          <Upload v-model="postForm.image_uri" />
        </el-form-item>
      </div>
    </el-form>

  </div>
</template>

<script>
import Tinymce from '@/components/Tinymce'
import Upload from '@/components/Upload/singleImage3'
import MDinput from '@/components/MDinput'
import Sticky from '@/components/Sticky' // 粘性header组件
// import { validateURL } from '@/utils/validate'
// import { fetchArticle } from '@/api/article'
// import { userSearch } from '@/api/remoteSearch'
import Warning from './Warning'
import { CommentDropdown, PlatformDropdown, SourceUrlDropdown, StatusDropdown } from './Dropdown'
import { CreatePostMutation } from '@/gql/queries/posts.gql'

const defaultForm = {
  status: 'draft',
  project: '',
  title: '', // 文章题目
  content: '', // 文章内容
  content_short: '', // 文章摘要
  source_uri: '', // 文章外链
  image_uri: '', // 文章图片
  display_time: undefined, // 前台展示时间
  id: undefined,
  platforms: ['a-platform'],
  comment_disabled: false,
  importance: 0
}

export default {
  name: 'ArticleDetail',
  components: {
    Tinymce, MDinput, Upload, Sticky, Warning, CommentDropdown, PlatformDropdown, SourceUrlDropdown,
    StatusDropdown
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
        this.$message({
          message: rule.field + '为必传项',
          type: 'error'
        })
        callback(new Error(rule.field + '为必传项'))
      } else {
        callback()
      }
    }
    /*
    const validateSourceUri = (rule, value, callback) => {
      if (value) {
        if (validateURL(value)) {
          callback()
        } else {
          this.$message({
            message: '外链url填写不正确',
            type: 'error'
          })
          callback(new Error('外链url填写不正确'))
        }
      } else {
        callback()
      }
    }
    */
    return {
      postForm: Object.assign({}, defaultForm),
      loading: false,
      userListOptions: [],
      rules: {
        status: [{ validator: validateRequire }],
        // image_uri: [{ validator: validateRequire }],
        title: [{ validator: validateRequire }],
        content: [{ validator: validateRequire }]
        // source_uri: [{ validator: validateSourceUri, trigger: 'blur' }]
      },
      tempRoute: {}
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
    contentShortLength() {
      return this.postForm.content_short.length
    },
    lang() {
      return this.$store.getters.language
    },
    projectid() {
      return this.$store.getters.projectid
    },
    form_userinfo_label() {
      return this.$t('post.userinfo') + ':'
    },
    form_userinfo_placeholder() {
      return this.$t('post.author_select')
    },
    form_display_time_label() {
      return this.$t('post.display_time') + ':'
    },
    form_display_time_placeholder() {
      return this.$t('post.display_time_placeholder')
    },
    form_author_rating_label() {
      return this.$t('post.author_rating') + ':'
    },
    form_desc_label() {
      return this.$t('post.content') + ':'
    },
    form_desc_placeholder() {
      return this.$t('post.content_placeholder')
    }
  },
  created() {
    if (this.isEdit) {
      const id = this.$route.params && this.$route.params.id
      this.fetchData(id)
    } else {
      this.postForm = Object.assign({}, defaultForm)
    }

    // Why need to make a copy of this.$route here?
    // Because if you enter this page and quickly switch tag, may be in the execution of the setTagsViewTitle function, this.$route is no longer pointing to the current page
    // https://github.com/PanJiaChen/vue-element-admin/issues/1221
    this.tempRoute = Object.assign({}, this.$route)
  },
  methods: {
    fetchData(id) {
      /*
      fetchArticle(id).then(response => {
        this.postForm = response.data
        // Just for test
        this.postForm.title += `   Article Id:${this.postForm.id}`
        this.postForm.content_short += `   Article Id:${this.postForm.id}`

        // Set tagsview title
        this.setTagsViewTitle()
      }).catch(err => {
        console.log(err)
      })
      */
    },
    setTagsViewTitle() {
      const title = this.lang === 'zh' ? '编辑文章' : 'Edit Article'
      const route = Object.assign({}, this.tempRoute, { title: `${title}-${this.postForm.id}` })
      this.$store.dispatch('updateVisitedView', route)
    },

    submitForm() {
      this.postForm.display_time = parseInt(this.display_time / 1000)
      // console.log(this.postForm)
      var _self = this
      this.$refs.postForm.validate(valid => {
        if (valid) {
          // const { title } = this.$data.postForm
          // const optimisticId = '' + Math.round(Math.random() * -1000000)
          var postVariables = this.postForm
          postVariables.project = _self.projectid
          // var status = 'active'
          this.loading = true
          console.error(postVariables)
          try {
            this.$apollo.mutate({
              mutation: CreatePostMutation,
              variables: postVariables,
              // see https://github.com/Akryum/vue-apollo-todos
              update: (store, { data: { createPost }}) => {
                /*
                // Add to ALL PROJECT list
                var queryKey = ALL_PROJECTS
                var itemKey = 'listAllProjects'

                // Ensure no duplicas (offline graphql ...) if not using optimisticId logic
                // filtering doesnt seems to work with nested elements ???
                const queryWithFilteringOptimistic = {
                  query: queryKey,
                  // variables: { filter: { items: { id: optimisticId } } }
                  variables: { filter: { items: { id: '' + createPost.id }}}
                }

                console.error('received data filtering of ' + optimisticId + ' and of ' + createPost.id)
                try {
                  if (createPost.id !== optimisticId) {
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
                            if (tmpAddItem[i].id === createPost.id) {
                              found = true
                            }
                          }
                        } catch (z) {
                          console.error(z.stack)
                        }

                        if (!found) {
                          tmpAddItem.push(createPost)
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
               createPost: {
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
      this.postForm.display_time = parseInt(this.display_time / 1000)
      console.log(this.postForm)
      this.$refs.postForm.validate(valid => {
        console.error(this.postForm)
        if (valid) {
          this.loading = true
          this.$notify({
            title: 'Creation',
            message: 'Successfully created post',
            type: 'success',
            duration: 2000
          })
          // this.postForm.status = 'published'
          this.loading = false
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    draftForm() {
      if (this.postForm.content.length === 0 || this.postForm.title.length === 0) {
        this.$message({
          message: 'No content or post title',
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
      // this.postForm.status = 'draft'
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
.createPost-container {
  position: relative;
  .createPost-main-container {
    padding: 40px 45px 20px 50px;
    .postInfo-container {
      position: relative;
      @include clearfix;
      margin-bottom: 10px;
      .postInfo-container-item {
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
