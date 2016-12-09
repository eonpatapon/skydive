<template>
  <div style="position:relative">
    <input class="form-control input-sm has-left-icon"
           :placeholder="placeholder"
           readonly
           @focus="select" 
           @mouseover="highlight(true)"
           @mouseout="highlight(false)"
           :value="value" />
    <span class="fa fa-crosshairs form-control-feedback"></span>
  </div>
</template>

<script>
export default {

  props: {
    value: {
      type: String,
      required: true,
    },
    placeholder: {
      type: String,
    },
  },

  methods: {

    highlight: function(bool) {
      topologyLayout.SetNodeClass(this.value, "highlighted", bool);
    },

    select: function() {
      var self = this;
      $(".topology-d3").on('click', function(e) {
        if (e.target.__data__) {
          self.$emit('input', e.target.__data__.ID);
          e.preventDefault();
          $(".topology-d3").off('click');
        }
      });
    }

  }

}
</script>
