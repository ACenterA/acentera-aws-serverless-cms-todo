<template>
  <el-dropdown :show-timeout="100" trigger="click">
    <el-button plain>{{ !is_draft ? $t('published') : $t('draft') }}
      <i class="el-icon-caret-bottom el-icon--right"/>
    </el-button>
    <el-dropdown-menu slot="dropdown" class="no-padding">
      <el-dropdown-item>
        <el-radio-group v-model="is_draft" style="padding: 10px;">
          <el-radio :label="draft">{{ $t('draft') }}</el-radio>
          <el-radio :label="published">{{ $t('publish') }}</el-radio>
          <el-radio :label="scheduled">{{ $t('scheduled') }}</el-radio>
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
      default: 'draft'
    }
  },
  computed: {
    is_draft: {
      get() {
        return this.value
      },
      set(val) {
        this.$emit('input', val)
      }
    }
  }
}
</script>
