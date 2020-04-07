<template>
  <el-dropdown :show-timeout="100" trigger="click">
    <el-button plain>{{ $t(is_editor) }}
      <i class="el-icon-caret-bottom el-icon--right" />
    </el-button>
    <el-dropdown-menu slot="dropdown" class="no-padding">
      <el-dropdown-item>
        <el-radio-group v-model="is_editor" style="padding: 10px;">
          <el-radio label="editor">{{ $t('editor') }}</el-radio>
          <el-radio label="none">{{ $t('none') }}</el-radio>
        </el-radio-group>
      </el-dropdown-item>
    </el-dropdown-menu>
  </el-dropdown>
</template>

<script>
export default {
  props: {
    value: {
      type: String,
      default: 'none'
    }
  },
  computed: {
    is_editor: {
      get() {
        return this.value
      },
      set(val) {
        console.error('emit value?')
        console.error(val)
        this.$emit('input', val)
      }
    }
  }
}
</script>
