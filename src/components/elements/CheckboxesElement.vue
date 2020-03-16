<template>
  <div>
    <label class="editable" contenteditable="true">{{ field.label }}</label>
    <div class="form-element">
      <div class="outer">
        <div v-for="option in splitStringByLine(field.options)" :key="option" class="checkbox">
          <label><input type="checkbox" value="">{{ option }}</label>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    field: {
      type: Object,
      default: null
    }
  },
  methods: {
    splitStringByLine: function(str) {
      return str.split('\n')
    }
  }
}
</script>
