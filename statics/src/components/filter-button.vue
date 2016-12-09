<template> 
  <div class="btn-group" role="group">
    <button class="btn btn-default btn-sm" readonly>
      {{ label }} : {{ value }}
    </button>
    <button class="btn btn-danger btn-sm" v-on:click="remove">
      <span class="glyphicon glyphicon-remove"></span>
    </button>
  </div>
</template>

<script>
export default {

  props: ['label', 'value'],

  methods: {

    remove() {
      this.$emit('remove', this.label, this.value);
    }

  }

}
</script>
