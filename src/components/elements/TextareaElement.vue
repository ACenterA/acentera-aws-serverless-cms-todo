<template>
  <div>
    <label class="editable" contenteditable="true">{{ field.label }}</label>
    <textarea class="form-control"/>
  </div>
</template>

<script>
export default {
  props: {
    field: {
      type: Object,
      default: null
    }
  }
}
</script>
