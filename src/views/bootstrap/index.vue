<template>
  <div class="login-container">
    <div v-if="isMissingDB">
      <!-- missing entry -->
      <el-form ref="passwordForm" :model="passwordForm" class="login-form" auto-complete="on" label-position="left">
        <div class="title-container">
          <br><br><br>
          <h3 class="title">{{ $t('login.siteConfigErrorReturn') }}</h3>
        </div>
      </el-form>
    </div>
    <div v-if="!isMissingDB && isMissingEntry">
      <div v-if="isAccountIdSet">
        <el-form ref="passwordForm" :model="passwordForm" class="login-form" auto-complete="on" label-position="left">
          <div class="title-container">
            <br><br><br>
            <a href="/"><h3 class="title">{{ $t('setupSite.initialize') }}</h3></a>
          </div>
          <el-input
            v-model="textInfo"
            :autosize="{ minRows: 2, maxRows: 4}"
            placeholder="Please input"
            type="textarea"
            readonly
          />

          <div class="margin-clear" />
          <el-checkbox v-model="passwordForm.eula" class="title-small">{{ $t('setupSite.eula') }}</el-checkbox>

          <el-button :loading="loading" type="primary" style="width:100%;margin-top:30px;margin-bottom:30px;" @click.native.prevent="handleInitializeAccount">{{ $t('login.loginSubmit') }}</el-button>

        </el-form>
      </div>
      <div v-if="!isAccountIdSet">
        <el-form ref="passwordForm" :model="passwordForm" class="login-form" auto-complete="on" label-position="left">
          <div class="title-container">
            <br><br><br>
            <a href="/"><h3 class="title">{{ $t('login.siteConfigErrorReturn') }}</h3></a>
          </div>
        </el-form>
      </div>
    </div>
    <div v-if="!isMissingDB && !isMissingEntry">
      <el-form ref="passwordForm" :model="passwordForm" class="login-form" auto-complete="on" label-position="left">
        <div class="title-container">
          <br><br><br>
          <h3 class="title">{{ $t('login.siteConfigErrorReturn') }}</h3>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
// import LangSelect from '@/components/LangSelect'

export default {
  name: 'Bootstrap',
  components: {
    // LangSelect
  },
  data() {
    /*
    const validatePassword = (rule, value, callback) => {
      if (value && value.length < 6) {
        callback(new Error(this.$t('login.PasswordDigitRequirements'))) // 'The password can not be less than 6 digits'))
      } else {
        callback()
      }
    }
    const validateJWT = (rule, value, callback) => {
      if (value && value.length < 12) {
        callback(new Error(this.$t('setupSite.JWTRequirements')))
      } else {
        callback()
      }
    }
    */
    return {
      passwordForm: {
        password: null,
        serverjwt: null,
        clientjwt: null,
        passwordConfirm: null,
        eula: false,
        email: false
      },
      passwordType: 'password',
      textInfo: '',
      loading: false
    }
  },
  computed: {
    ...mapGetters([
      'needMFARegistrationStr',
      'needMFARegistration',
      'getCognitoUser',
      'isLoginCodeReset',
      'isAccountIdSet'
    ]),
    textinfo: function() {
      // $t('setupSite.initialize')
      return 'aaaaa'
    },
    settings: function() {
      if (this.$store.getters.settings) {
        return this.$store.getters.settings
      }
      return {}
    },
    isMissingEntry() {
      if (this.$store && this.$store.getters && this.$store.getters.settings) {
        return this.$store.getters.settings.missingSiteEntry === true
      }
      return false
    },
    isMissingDB() {
      if (this.$store && this.$store.getters && this.$store.getters.settings) {
        return this.$store.getters.settings.missingDB === true
      }
      return true
    }
  },
  created() {
    this.textInfo = window.app.$t('setupSite.textinfo')
  },
  destroyed() {
  },
  methods: {
    showPwd() {
      if (this.passwordType === 'password') {
        this.passwordType = ''
      } else {
        this.passwordType = 'password'
      }
    },
    resetForm() {
      this.passwordForm = {
        password: null,
        serverjwt: null,
        clientjwt: null,
        passwordConfirm: null,
        eula: false,
        email: false
      }
    },
    handleInitializeAccount() {
      this.loading = true
      var self = this
      this.$store.dispatch('PerformPluginAppInitialization', self.passwordForm).then((resp) => {
        if (resp === true) {
          self.resetForm()
          // window.location.href = '/'

          window.location.href = '/'
          window.location.reload(true)

          // self.$router.push({ path: '/' })
          setTimeout(function() {
            self.loading = false
          }, 2000)
        } else {
          // todo: show error ??
          self.loading = false
        }
      }).catch(() => {
        // todo: show error message
        self.loading = false
      })
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss">
  /* 修复input 背景不协调 和光标变色 */
  /* Detail see https://github.com/PanJiaChen/vue-element-admin/pull/927 */

  $bg:#283443;
  $light_gray:#eee;
  $cursor: #fff;

  @supports (-webkit-mask: none) and (not (cater-color: $cursor)) {
    .login-container .el-input input{
      color: $cursor;
      &::first-line {
        color: $light_gray;
      }
    }
  }

  /* reset element-ui css */
  .login-container {
    .el-input {
      display: inline-block;
      height: 47px;
      width: 85%;
      input {
        background: transparent;
        border: 0px;
        -webkit-appearance: none;
        border-radius: 0px;
        padding: 12px 5px 12px 15px;
        color: $light_gray;
        height: 47px;
        caret-color: $cursor;
        &:-webkit-autofill {
          -webkit-box-shadow: 0 0 0px 1000px $bg inset !important;
          -webkit-text-fill-color: $cursor !important;
        }
      }
    }
    .el-form-item {
      border: 1px solid rgba(255, 255, 255, 0.1);
      background: rgba(0, 0, 0, 0.1);
      border-radius: 5px;
      color: #454545;
    }
  }
</style>

<style rel="stylesheet/scss" lang="scss" scoped>
$bg:#2d3a4b;
$dark_gray:#889aa4;
$light_gray:#eee;

$break-small: 320px;
$break-large: 700px;

.title-container {
  padding-top:20px;
}

.el-button+.el-button {
    margin-left: 0px;
}

.login-container {
  overflow: auto;
  position: fixed;
  height: 100%;
  width: 100%;
  background-color: $bg;
  .login-form {
    position: absolute;
    left: 0;
    right: 0;
    width: 520px;
    max-width: 100%;
    padding: 35px 35px 15px 35px;
    margin: 120px auto;
  }
  .tips {
    font-size: 14px;
    color: #fff;
    margin-bottom: 10px;
    span {
      &:first-of-type {
        margin-right: 16px;
      }
    }
  }
  .svg-container {
    padding: 6px 5px 6px 15px;
    color: $dark_gray;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
    &_login {
      font-size: 20px;
    }
  }
  .title-container {
    position: relative;
    .title {
      font-size: 26px;
      color: $light_gray;
      margin: 0px auto 40px auto;
      text-align: center;
      font-weight: bold;
    }
    .title-small {
      font-size: 18px;
      color: $light_gray;
      margin: 0px auto 30px auto;
      text-align: center;
      font-weight: bold;
    }
    .titleFirstTime {
      color: $light_gray;
      margin: 0px auto 10px auto;
      text-align: center;
      font-weight: bold;
    }
    .set-language {
      color: #fff;
      position: absolute;
      top: 5px;
      right: 0px;
    }
  }
  .show-pwd {
    position: absolute;
    right: 10px;
    top: 7px;
    font-size: 16px;
    color: $dark_gray;
    cursor: pointer;
    user-select: none;
  }
  .thirdparty-button {
    position: absolute;
    right: 35px;
    bottom: 28px;
  }
}

@media screen and (max-width: $break-large) {
  .login-container {
    .login-form {
      margin: 10px auto;
    }
  }
}
.forceblack input::placeholder {
    color: black;
    opacity: 1; /* Firefox */
}
.forcewhite {
  background-color: white;
  color: black!important;
}

.black {
  color: black!important;
}
.el-form-item.forceblack {
  background-color:black!important;
}

.el-input-number--medium {
  width: 80%
}

.margin-clear {
  margin-bottom: 20px;
}

.title-small-nopad {
  font-size: 12px;
  color: $light_gray;
  text-align: center;
  font-weight: bold;
}
.title-small-nopad-nocenter {
  font-size: 12px;
  color: $light_gray;
  font-weight: bold;
}

.title-small {
  font-size: 12px;
  color: $light_gray;
  margin: 0px auto 30px auto;
  text-align: center;
  font-weight: bold;
}

</style>
